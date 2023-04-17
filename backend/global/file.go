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

	sql "github.com/FloatTech/sqlite"
	"github.com/corona10/goimagehash"
	base14 "github.com/fumiama/go-base16384"
	"github.com/fumiama/go-docx"
	"github.com/sirupsen/logrus"

	"github.com/fumiama/paper-manager/backend/utils"
)

const (
	FileTableList         = "lst"
	FileTableFile         = "file"
	FileTableTempFile     = "tmpfile"
	FileTableQuestion     = "question"
	FileTableTempQuestion = "tmpqstn"
)

var (
	ErrMajorSplitsTooShort       = errors.New("major splits too short")
	ErrEmptyClass                = errors.New("empty class")
	ErrHasntAnalyzed             = errors.New("hasn't analyzed")
	ErrNoGetFileStatusPermission = errors.New("no get file status permission")
)

func init() {
	err := FileDB.db.Open(time.Hour)
	if err != nil {
		panic(err)
	}
	err = FileDB.db.Create(FileTableList, &List{})
	if err != nil {
		panic(err)
	}
	err = FileDB.db.Create(FileTableFile, &File{},
		"FOREIGN KEY(ListID) REFERENCES "+FileTableList+"(ID)",
	)
	if err != nil {
		panic(err)
	}
	err = FileDB.db.Create(FileTableTempFile, &File{})
	if err != nil {
		panic(err)
	}
	err = FileDB.db.Create(FileTableQuestion, &Question{})
	if err != nil {
		panic(err)
	}
	err = FileDB.db.Create(FileTableTempQuestion, &Question{})
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

// File stores to paper/Class/2022-2023学年/第一学期/期末/A/xxx.docx
type File struct {
	ID        int64 // ID is the first 8 bytes of the original file's md5
	ListID    int   // ListID is the foreign key to List(ID)
	Year      StudyYear
	Type      PaperType
	Date      uint32        // Date is the yyyymmdd of 考试日期
	Time      time.Duration // Time is 考试时长
	Class     string        // Class is 考试科目
	Rate      string        // Rate is 成绩构成比例
	Questions []byte        // Questions is for []QuestionJSON
}

// StudyYear 学年
type StudyYear uint16

// String ex. 2022-2023学年
func (sy StudyYear) String() string {
	next := sy + 1
	return strconv.Itoa(int(sy)) + "-" + strconv.Itoa(int(next)) + "学年"
}

// AddFile from lst and copy it to analyzed path.
// The para reg must belong to a valid user
func (f *FileDatabase) AddFile(lstid int, reg *Regex, istemp bool, progress func(uint)) error {
	user, err := UserDB.GetUserByID(reg.ID)
	if err != nil {
		return err
	}
	if !user.IsFileManager() && !istemp {
		return ErrInvalidRole
	}
	progress(1)
	f.mu.RLock()
	lst, err := sql.Find[List](&f.db, FileTableList, "WHERE ID="+strconv.Itoa(lstid))
	f.mu.RUnlock()
	if err != nil {
		return err
	}
	if lst.Path == "" || strings.Contains(lst.Path, "..") {
		return os.ErrNotExist
	}
	tempath := lst.Path
	docf, err := os.Open(tempath)
	if err != nil {
		return err
	}
	defer docf.Close()
	progress(2)
	h := md5.New()
	_, err = io.Copy(h, docf)
	if err != nil {
		return err
	}
	var buf [md5.Size]byte
	id := int64(binary.LittleEndian.Uint64(h.Sum(buf[:0])))
	_, err = docf.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	stat, err := docf.Stat()
	if err != nil {
		return err
	}
	sz := stat.Size()
	progress(3)
	doc, err := docx.Parse(docf, sz)
	if err != nil {
		return err
	}
	progress(5)
	doc.Document.Body.DropDrawingOf("NilPicture")
	majorre, err := regexp.Compile(reg.Major)
	if err != nil {
		return err
	}
	docs := doc.SplitByParagraph(docx.SplitDocxByPlainTextRegex(majorre))
	if len(docs) < 2 {
		return ErrMajorSplitsTooShort
	}
	progress(9)
	// filling File struct
	file := &File{
		ID:     id,
		ListID: *lst.ID,
	}
	titlere, err := regexp.Compile(reg.Title)
	if err != nil {
		return err
	}
	classre, err := regexp.Compile(reg.Class)
	if err != nil {
		return err
	}
	opclre, err := regexp.Compile(reg.OpenCl)
	if err != nil {
		return err
	}
	datere, err := regexp.Compile(reg.Date)
	if err != nil {
		return err
	}
	timere, err := regexp.Compile(reg.Time)
	if err != nil {
		return err
	}
	ratere, err := regexp.Compile(reg.Rate)
	if err != nil {
		return err
	}
	progress(10)
	for _, it := range docs[0].Document.Body.Items {
		if p, ok := it.(*docx.Paragraph); ok {
			text := p.String()
			title := titlere.FindStringSubmatch(text)
			if len(title) >= 5 {
				years, semesters, mfs, abs := title[1], title[2], title[3], title[4]
				y, err := strconv.Atoi(years)
				if err != nil {
					return err
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
			if len(class) >= 3 {
				file.Class = class[2]
			}
			opcl := opclre.FindStringSubmatch(text)
			if len(opcl) >= 2 {
				file.Type = file.Type.SetOpenClose(opcl[1])
			}
			date := datere.FindStringSubmatch(text)
			if len(date) >= 4 {
				y, m, d := date[1], date[2], date[3]
				if y != "" && m != "" {
					if d == "" {
						d = "1"
					}
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
			if len(rate) >= 3 {
				file.Rate = rate[2]
			}
		}
	}
	progress(19)
	if file.Class == "" || strings.Contains(file.Class, "..") || strings.ContainsAny(file.Class, `/\`) {
		return ErrEmptyClass
	}
	filebasepath := ""
	if istemp {
		filebasepath = PaperFolder + "temp/" + strconv.Itoa(*user.ID) + "/" + file.Class + "/"
	} else {
		filebasepath = fmt.Sprintf(
			PaperFolder+file.Class+"/%v/%v/%v/%c/",
			file.Year, file.Type.FirstSecond(), file.Type.MiddleFinal(), file.Type.AB(),
		)
	}
	questionpath := filebasepath + "questions/"
	err = os.MkdirAll(questionpath, 0755)
	if err != nil {
		return err
	}
	docs = docs[1:]
	// parse questions
	subre, err := regexp.Compile(reg.Sub)
	if err != nil {
		return err
	}
	filequestions := make([]QuestionJSON, 0, len(docs))
	lst.QuesC = 0
	progress(20)
	p := uint(20)
	delta := uint(70 / len(docs))
	if delta == 0 {
		delta = 1
	}
	for _, majordoc := range docs {
		p += delta
		if p > 90 {
			p = 90
		}
		progress(p)
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
		if len(subdocs) < 2 {
			continue
		}
		subdocs = subdocs[1:]
		majorq.Sub = make([]QuestionJSON, 0, len(subdocs))
		for _, subdoc := range subdocs {
			sb := bytes.NewBuffer(make([]byte, 0, 4096))
			for _, it := range subdoc.Document.Body.Items {
				sb.WriteString(fmt.Sprint(it))
			}
			m := md5.Sum(sb.Bytes())
			que := &Question{
				ID:    int64(binary.LittleEndian.Uint64(m[:8])),
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
						if word != "" && !strings.Contains("\n 。，、的是使,.（）()1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ", word) {
							if v[word] == 0 { // 二值化
								v[word] = 1
							}
						}
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
			f.mu.RLock()
			err = f.db.FindFor(FileTableQuestion, &q, "", func() error {
				r, err := q.GetDuplicateRate(que)
				if err != nil {
					logrus.Warnln("[global.AddFile] GetDuplicateRate err:", err)
					return err
				}
				if r < 0.5 {
					return nil
				}
				var buf [8]byte
				binary.LittleEndian.PutUint64(buf[:], uint64(q.ID))
				dupmap[hex.EncodeToString(buf[:])] = r
				return nil
			})
			f.mu.RUnlock()
			if err == nil && len(dupmap) > 0 {
				que.Dup, _ = json.Marshal(dupmap)
			}
			w := bytes.NewBuffer(make([]byte, 0, 65536))
			_, err = subdoc.WriteTo(w)
			var buf [8]byte
			binary.LittleEndian.PutUint64(buf[:], uint64(que.ID))
			queidstr := hex.EncodeToString(buf[:])
			if err == nil {
				m5 := md5.Sum(w.Bytes())
				quepath := questionpath + hex.EncodeToString(m5[:]) + ".docx"
				qf, err := os.Create(quepath)
				if err == nil {
					_, _ = io.Copy(qf, w)
					_ = qf.Close()
				}
				que.Path = quepath
				if istemp {
					f.mu.Lock()
					_ = f.db.Insert(FileTableTempQuestion, que)
					f.mu.Unlock()
				} else {
					f.mu.Lock()
					for k, v := range dupmap {
						err = f.db.Find(FileTableQuestion, &q, "WHERE ID=0x"+k)
						if err == nil {
							thismap := make(map[string]float64, 64)
							err := json.Unmarshal(q.Dup, &thismap)
							if err == nil {
								thismap[queidstr] = v
								q.Dup, err = json.Marshal(thismap)
								if err == nil {
									_ = f.db.Insert(FileTableQuestion, &q)
								}
							}
						}
					}
					_ = f.db.Insert(FileTableQuestion, que)
					f.mu.Unlock()
				}
			}
			r := 0.0
			for _, v := range dupmap {
				if v > r {
					r = v
				}
			}
			majorq.Sub = append(majorq.Sub, QuestionJSON{
				Name: queidstr, //TODO: fill sub points
			})
		}
		filequestions = append(filequestions, majorq)
		lst.QuesC += len(majorq.Sub)
	}
	progress(90)
	file.Questions, _ = json.Marshal(filequestions)
	_, err = docf.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	lst.Path = filebasepath + file.Class + ".docx"
	lst.HasntAnalyzed = false
	lst.Desc = fmt.Sprintf("%s%v%v%v%c卷",
		file.Class, file.Year, file.Type.FirstSecond(), file.Type.MiddleFinal(), file.Type.AB(),
	)
	dstf, err := os.Create(lst.Path)
	if err != nil {
		return err
	}
	defer dstf.Close()
	_, err = io.Copy(dstf, docf)
	if err != nil {
		return err
	}
	progress(95)
	f.mu.Lock()
	if istemp {
		err = f.db.Insert(FileTableTempFile, file)
		lst.IsTemp = true
	} else {
		err = f.db.Insert(FileTableFile, file)
		lst.IsTemp = false
	}
	_ = f.db.Insert(FileTableList, &lst)
	f.mu.Unlock()
	progress(100)
	return err
}

// DelFile by listid
func (f *FileDatabase) DelFile(lstid, uid int, istemp bool) error {
	user, err := UserDB.GetUserByID(uid)
	if err != nil {
		return err
	}
	if !user.IsSuper() {
		return ErrInvalidRole
	}
	ftable := ""
	if istemp {
		ftable = FileTableTempFile
	} else {
		ftable = FileTableFile
	}
	f.mu.RLock()
	lst, err := sql.Find[List](&f.db, FileTableList, "WHERE ID="+strconv.Itoa(lstid))
	f.mu.RUnlock()
	if err != nil {
		return err
	}
	if lst.Path == "" || strings.Contains(lst.Path, "..") {
		return os.ErrNotExist
	}
	err = f.db.Del(FileTableList, "WHERE ID="+strconv.Itoa(lstid))
	if err != nil {
		return err
	}
	if lst.HasntAnalyzed {
		return os.RemoveAll(lst.Path)
	}
	i := strings.LastIndex(lst.Path, "/")
	if i <= 0 {
		return os.ErrNotExist
	}
	parentfolder := lst.Path[:i]
	if utils.IsNotExist(parentfolder) {
		return os.ErrNotExist
	}
	f.mu.RLock()
	file, err := sql.Find[File](&f.db, ftable, "WHERE ListID="+strconv.Itoa(lstid))
	f.mu.RUnlock()
	if err != nil {
		return err
	}
	err = f.db.Del(ftable, "WHERE ListID="+strconv.Itoa(lstid))
	if err != nil {
		return err
	}
	ques := make([]QuestionJSON, 0, 64)
	err = json.Unmarshal(file.Questions, &ques)
	if err != nil {
		return err
	}
	for _, q := range ques {
		q.Delete(f, istemp)
	}
	return os.RemoveAll(parentfolder)
}

// GetFile get analyzed file's structure from List(ID)
func (f *FileDatabase) GetFile(lstid, uid int) (file *File, sz int64, istemp bool, err error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	lst, err := sql.Find[List](&f.db, FileTableList, "WHERE ID="+strconv.Itoa(lstid))
	if err != nil {
		return
	}
	if lst.HasntAnalyzed {
		err = ErrHasntAnalyzed
		return
	}
	if lst.IsTemp && lst.Uploader != uid {
		err = ErrNoGetFileStatusPermission
		return
	}
	ftable := ""
	if lst.IsTemp {
		ftable = FileTableTempFile
	} else {
		ftable = FileTableFile
	}
	filestruct, err := sql.Find[File](&f.db, ftable, "WHERE ListID="+strconv.Itoa(lstid))
	if err != nil {
		return
	}
	return &filestruct, lst.Size, lst.IsTemp, nil
}
