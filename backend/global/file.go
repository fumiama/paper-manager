package global

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"

	"github.com/corona10/goimagehash"
	base14 "github.com/fumiama/go-base16384"
	"github.com/fumiama/go-docx"

	"github.com/fumiama/paper-manager/backend/utils"
)

const (
	FileTableFile     = "file"
	FileTableQuestion = "question"
)

var (
	ErrMajorSplitsTooShort = errors.New("major splits too short")
)

// PaperType [4 开 一页纸 闭] [4 上下] [4 中末] [4 AB]
type PaperType uint16

// AB default A
func (pt PaperType) AB() byte {
	switch pt & 0x0f {
	case 1:
		return 'A'
	case 2:
		return 'B'
	default:
		return 'A'
	}
}

func (pt PaperType) SetAB(x byte) PaperType {
	n := PaperType(0)
	switch x {
	case 'A':
		n = 1
	case 'B':
		n = 2
	}
	return pt | n
}

// MiddleFinal default 平时
func (pt PaperType) MiddleFinal() string {
	switch (pt & 0xf0) >> 4 {
	case 1:
		return "期中"
	case 2:
		return "期末"
	default:
		return "平时"
	}
}

func (pt PaperType) SetMiddleFinal(x string) PaperType {
	n := PaperType(0)
	switch x {
	case "中":
		n = 1 << 4
	case "末":
		n = 2 << 4
	}
	return pt | n
}

// FirstSecond default is 年度
func (pt PaperType) FirstSecond() string {
	switch (pt & 0x0f00) >> 8 {
	case 1:
		return "第1学期"
	case 2:
		return "第2学期"
	default:
		return "年度"
	}
}

func (pt PaperType) SetFirstSecond(x byte) PaperType {
	n := PaperType(0)
	switch x {
	case '1':
		n = 1 << 8
	case '2':
		n = 2 << 8
	}
	return pt | n
}

// OpenClose default 闭卷
func (pt PaperType) OpenClose() string {
	switch (pt & 0xf000) >> 12 {
	case 1:
		return "开卷"
	case 2:
		return "一页纸开卷"
	case 3:
		return "闭卷"
	default:
		return "闭卷"
	}
}

func (pt PaperType) SetOpenClose(x string) PaperType {
	n := PaperType(0)
	switch x {
	case "开卷":
		n = 1 << 12
	case "一页纸开卷":
		n = 2 << 12
	case "闭卷":
		n = 3 << 12
	}
	return pt | n
}

// StudyYear 学年
type StudyYear uint16

// String ex. 2022-2023学年
func (sy StudyYear) String() string {
	next := sy + 1
	return strconv.Itoa(int(sy)) + "-" + strconv.Itoa(int(next)) + "学年"
}

func init() {
	err := FileDB.db.Open(time.Hour)
	if err != nil {
		panic(err)
	}
	err = FileDB.db.Create(FileTableFile, &File{})
	if err != nil {
		panic(err)
	}
	err = FileDB.db.Create(FileTableQuestion, &Question{})
	if err != nil {
		panic(err)
	}
	err = FileDB.db.Close()
	if err != nil {
		panic(err)
	}
	err = os.Chmod(FileDB.db.DBPath, 0600)
	if err != nil {
		panic(err)
	}
	err = FileDB.db.Open(time.Hour)
	if err != nil {
		panic(err)
	}
}

type File struct {
	ID        uint64 // ID is the first 8 bytes of the original file's md5
	Year      StudyYear
	Type      PaperType
	Date      uint32        // Date is the yyyymmdd of 考试日期
	UID       int           // UID is the uploader's ID
	UpTime    int64         // UpTime is time.Now().Unix() when uploading
	Size      int64         // Size of the original file
	Time      time.Duration // Time is 考试时长
	Class     string        // Class is 考试科目
	Rate      string        // Rate is 成绩构成比例
	Path      string        // Path is like paper/Class/2023/第一学期/期末/A/xxx.docx
	Questions []byte        // Questions is for json struct QuestionJSON
}

// AddFile from FileFolder+tempath and copy it to File.Path.
// The para res must belong to a valid user
func (f *FileDatabase) AddFile(tempath string, reg *Regex, istemp bool, progress func(uint)) (*File, error) {
	user, err := UserDB.GetUserByID(reg.ID)
	if err != nil {
		return nil, err
	}
	if !user.IsFileManager() && !istemp {
		return nil, ErrInvalidRole
	}
	if strings.Contains(tempath, "..") {
		return nil, os.ErrNotExist
	}
	tempath = FileFolder + tempath
	docf, err := os.Open(tempath)
	if err != nil {
		return nil, err
	}
	defer docf.Close()
	h := md5.New()
	_, err = io.Copy(h, docf)
	if err != nil {
		return nil, err
	}
	var buf [md5.Size]byte
	id := binary.LittleEndian.Uint64(h.Sum(buf[:0]))
	_, err = docf.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	stat, err := docf.Stat()
	if err != nil {
		return nil, err
	}
	sz := stat.Size()
	doc, err := docx.Parse(docf, sz)
	if err != nil {
		return nil, err
	}
	doc.Document.Body.DropDrawingOf("NilPicture")
	majorre, err := regexp.Compile(reg.Major)
	if err != nil {
		return nil, err
	}
	docs := doc.SplitByParagraph(docx.SplitDocxByPlainTextRegex(majorre))
	if len(docs) < 2 {
		return nil, ErrMajorSplitsTooShort
	}
	// filling File struct
	file := &File{
		ID:     id,
		UID:    *user.ID,
		UpTime: time.Now().Unix(),
		Size:   sz,
	}
	titlere, err := regexp.Compile(reg.Title)
	if err != nil {
		return nil, err
	}
	classre, err := regexp.Compile(reg.Class)
	if err != nil {
		return nil, err
	}
	opclre, err := regexp.Compile(reg.OpenCl)
	if err != nil {
		return nil, err
	}
	datere, err := regexp.Compile(reg.Date)
	if err != nil {
		return nil, err
	}
	timere, err := regexp.Compile(reg.Time)
	if err != nil {
		return nil, err
	}
	ratere, err := regexp.Compile(reg.Rate)
	if err != nil {
		return nil, err
	}
	for _, it := range docs[0].Document.Body.Items {
		if p, ok := it.(*docx.Paragraph); ok {
			text := p.String()
			title := titlere.FindStringSubmatch(text)
			if len(title) >= 5 {
				years, semesters, mfs, abs := title[1], title[2], title[3], title[4]
				y, err := strconv.Atoi(years)
				if err != nil {
					return nil, err
				}
				file.Year = StudyYear(y)
				if len(semesters) > 0 {
					file.Type = file.Type.SetFirstSecond(semesters[0])
				}
				file.Type = file.Type.SetMiddleFinal(mfs)
				if len(abs) > 0 {
					file.Type = file.Type.SetAB(abs[0])
				}
			}
			class := classre.FindStringSubmatch(text)
			if len(class) >= 2 {
				file.Class = class[1]
			}
			opcl := opclre.FindStringSubmatch(text)
			if len(opcl) >= 2 {
				file.Type = file.Type.SetOpenClose(opcl[1])
			}
			date := datere.FindStringSubmatch(text)
			if len(date) >= 4 {
				y, m, d := date[1], date[2], date[3]
				if y != "" && m != "" && d != "" {
					yyyy, err := strconv.ParseUint(y, 10, 64)
					if err == nil && yyyy > 1600 {
						mm, err := strconv.ParseUint(m, 10, 64)
						if err == nil && mm >= 1 && mm <= 12 {
							dd, err := strconv.ParseUint(d, 10, 64)
							if err == nil && dd >= 1 && dd <= 31 {
								file.Date = uint32(yyyy*10000 + mm*100 + dd)
							}
						}
					}
				}
			}
			times := timere.FindStringSubmatch(text)
			if len(times) >= 2 {
				min, err := strconv.Atoi(times[1])
				if err == nil && min > 0 {
					file.Time = time.Minute * time.Duration(min)
				}
			}
			rate := ratere.FindStringSubmatch(text)
			if len(rate) >= 2 {
				file.Rate = rate[1]
			}
		}
	}
	docs = docs[1:]
	// parse questions
	subre, err := regexp.Compile(reg.Sub)
	if err != nil {
		return nil, err
	}
	for _, majordoc := range docs {
		majorq := QuestionJSON{}
		for _, it := range majordoc.Document.Body.Items {
			if p, ok := it.(*docx.Paragraph); ok {
				text := p.String()
				majorinfo := majorre.FindStringSubmatch(text)
				if len(majorinfo) >= 6 {
					name, points := majorinfo[2], majorinfo[5]
					majorq.Name = name
					majorq.Points, _ = strconv.Atoi(points)
				}
			}
		}
		subdocs := majordoc.SplitByParagraph(docx.SplitDocxByPlainTextRegex(subre))
		majorq.Sub = make([]QuestionJSON, 0, len(subdocs))
		for _, subdoc := range subdocs {
			sb := bytes.NewBuffer(make([]byte, 0, 4096))
			for _, it := range subdoc.Document.Body.Items {
				sb.WriteString(fmt.Sprint(it))
			}
			m := md5.Sum(sb.Bytes())
			que := &Question{
				ID:    binary.LittleEndian.Uint64(m[:8]),
				Plain: base14.BytesToString(sb.Bytes()),
				Images: func() []byte {
					m := make(map[string]string)
					_ = subdoc.RangeRelationships(func(r *docx.Relationship) error {
						if r.Type != docx.REL_IMAGE {
							return nil
						}
						if r.Target == "" {
							return nil
						}
						i := strings.LastIndex(r.Target, "/")
						if i < 0 {
							return nil
						}
						name := r.Target[i+1:]
						if name == "" {
							return nil
						}
						md := subdoc.Media(name)
						if md == nil {
							return nil
						}
						img, _, err := image.Decode(bytes.NewReader(md.Data))
						if err != nil {
							return nil
						}
						dh, err := goimagehash.DifferenceHash(img)
						if err != nil {
							return nil
						}
						var buf [8]byte
						binary.LittleEndian.PutUint64(buf[:], dh.GetHash())
						m[name] = hex.EncodeToString(buf[:])
						return nil
					})
					if len(m) == 0 {
						return nil
					}
					data, err := json.Marshal(m)
					if err != nil {
						return nil
					}
					return data
				}(),
				Vector: func() []byte {
					words := utils.Segmenter.Cut(base14.BytesToString(sb.Bytes()), true)
					if len(words) == 0 {
						return nil
					}
					v := make(map[string]uint8, len(words)*2)
					for _, word := range words {
						v[word]++
					}
					data, err := json.Marshal(v)
					if err != nil {
						return nil
					}
					return data
				}(),
			}
			var q Question
			dupmap := make(map[string]float64, 64)
			FileDB.mu.RLock()
			err = FileDB.db.FindFor(FileTableQuestion, &q, "", func() error {
				r, err := q.GetDuplicateRate(que)
				if err != nil {
					return err
				}
				var buf [8]byte
				binary.LittleEndian.PutUint64(buf[:], q.ID)
				dupmap[hex.EncodeToString(buf[:])] = r
				return nil
			})
			FileDB.mu.RUnlock()
			if err == nil {

			}
		}
	}
	return file, nil
}

// QuestionJSON is the struct representation of File.Questions
type QuestionJSON struct {
	Name   string         `json:"name"`   // Name is name or Question ID
	Points int            `json:"points"` // Points is sum of subs' points or self
	Rate   float64        `json:"rate"`   // Rate is the avg(non-leaf) or max(leaf) similarity
	Sub    []QuestionJSON `json:"sub,omitempty"`
}

type Question struct {
	ID     uint64 // ID is the first 8 bytes of the Plain's md5
	Plain  string // Plain is the plain text of the question (like markdown format)
	Images []byte // Images is json of the image dhash in XML, ex. ['rId1': '1234567890abcdef', ...]
	Vector []byte // Vector is json of {word: freq, ...}
	Dup    []byte // Dup is json of {queid: rate, ...}
}

// GetDuplicateRate calc q & que's dup rate
func (q *Question) GetDuplicateRate(que *Question) (float64, error) {
	v1, v2 := make(map[string]uint8, 64), make(map[string]uint8, 64)
	m1, m2 := make(map[string]string, 64), make(map[string]string, 64)
	err := json.Unmarshal(q.Images, &m1)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(que.Images, &m2)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(q.Vector, &v1)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(que.Vector, &v2)
	if err != nil {
		return 0, err
	}
	imgdsts := uint64(0)
	for _, dhstr2 := range m2 {
		d, err := hex.DecodeString(dhstr2)
		if err != nil {
			return 0, err
		}
		dh2 := goimagehash.NewImageHash(binary.LittleEndian.Uint64(d), goimagehash.DHash)
		r := 0
		for _, dhstr1 := range m1 {
			d, err := hex.DecodeString(dhstr1)
			if err != nil {
				return 0, err
			}
			dh1 := goimagehash.NewImageHash(binary.LittleEndian.Uint64(d), goimagehash.DHash)
			dst, err := dh2.Distance(dh1)
			if err != nil {
				return 0, err
			}
			if dst > r {
				r = dst
			}
		}
		imgdsts += uint64(r)
	}
	imgdupr := float64(imgdsts) / float64(len(m2)) / 64.0
	v1space := make([]uint8, 0, len(v1)+len(v2))
	v2space := make([]uint8, 0, len(v1)+len(v2))
	for k, v := range v1 {
		v1space = append(v1space, v)
		if tv, ok := v2[k]; ok {
			v2space = append(v2space, tv)
			delete(v2, k)
		} else {
			v2space = append(v2space, 0)
		}
	}
	for _, v := range v2 {
		v1space = append(v1space, 0)
		v2space = append(v2space, v)
	}
	return utils.Similarity(v1space, v2space) + imgdupr/2.0, nil
}
