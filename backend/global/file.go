package global

import (
	"os"
	"strconv"
	"time"
)

const (
	FileTableFile     = "file"
	FileTableQuestion = "question"
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
	XML    []byte // XML is the OpenXML bytes of the question
	Images []byte // Images is json of the image paths in XML, ex. ['md5.jpg', 'md5.png', ...]
	Vector []byte // Vector is json of {word: rate, ...} freq
	Dup    []byte // Dup is json of Duplication struct
}

// Duplication is the struct representation of Question.Dup
type Duplication struct {
	ID   string        `json:"id"`   // ID is hex string for json's 53 bits number
	Rate float64       `json:"rate"` // Rate is the avg(non-leaf) or max(leaf) similarity
	To   []Duplication `json:"to,omitempty"`
}
