package global

import (
	"errors"
	"os"
	"strconv"

	sql "github.com/FloatTech/sqlite"
	"github.com/fumiama/go-docx"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvalidGenerateConfig          = errors.New("invalid generate config")
	ErrMajorTooLarge                  = errors.New("major too large")
	ErrNoSuchMajor                    = errors.New("no such major")
	ErrNoEnoughQuestionToMatchRequire = errors.New("no enough question to match require")
	ErrRateLimitExceeded              = errors.New("rate limit exceeded")
)

// GenerateConfig 试卷生成配置
type GenerateConfig struct {
	Distribution map[string]uint // Distribution is map[majorname]subcount
	RateLimit    float64         // RateLimit 重复率上限
	YearStart    StudyYear       // YearStart 起始年份（空则直到最旧）
	YearEnd      StudyYear       // YearEnd 截止年份（空则直到最新）
	TypeMask     PaperType       // TypeMask & File.Type != 0 则匹配
}

// GenerateFile 用一些限定条件生成新试卷, 云端不保存
func (f *FileDatabase) GenerateFile(config *GenerateConfig) (docf *docx.Docx, err error) {
	if config == nil || config.Distribution == nil || len(config.Distribution) == 0 {
		return nil, ErrInvalidGenerateConfig
	}
	if len(config.Distribution) > 10 {
		return nil, ErrMajorTooLarge
	}
	mm := map[string]struct{}{}
	for _, m := range f.GetMajors() {
		mm[m] = struct{}{}
	}
	for n := range config.Distribution {
		if _, ok := mm[n]; !ok {
			return nil, ErrNoSuchMajor
		}
	}
	docf = docx.NewA4()
	f.mu.RLock()
	defer f.mu.RUnlock()
	i := 0
	for n, c := range config.Distribution {
		if c <= 0 {
			continue
		}
		docf.AddParagraph().AddText(string([]rune("一二三四五六七八九十")[i]) + "、" + n).Size("30").Bold()
		cond := " WHERE"
		hasfront := false
		if config.YearStart > 0 {
			cond += " Year>=" + strconv.Itoa(int(config.YearStart))
			hasfront = true
		}
		if config.YearEnd > 0 {
			if hasfront {
				cond += " AND"
			}
			cond += " Year<=" + strconv.Itoa(int(config.YearEnd))
			hasfront = true
		}
		if config.TypeMask > 0 {
			if hasfront {
				cond += " AND"
			}
			typmsk := strconv.Itoa(int(config.TypeMask))
			cond += " (Type&" + typmsk + ")==" + typmsk
			hasfront = true
		}
		var ques []*Question
		q := ""
		if hasfront {
			q = "SELECT * FROM " + FileTableQuestion +
				" WHERE Major='" + n + "' AND ListID IN (SELECT DISTINCT ListID FROM " +
				FileTableFile + cond +
				") ORDER BY RANDOM() limit " + strconv.Itoa(int(c)) + ";"
			ques, err = sql.QueryAll[Question](&f.db, q)
		} else {
			q = "SELECT * FROM " + FileTableQuestion +
				" WHERE Major='" + n + "' ORDER BY RANDOM() limit " + strconv.Itoa(int(c)) + ";"
			ques, err = sql.QueryAll[Question](&f.db, q)
		}
		if err != nil {
			logrus.Warnln(err, q)
			return nil, err
		}
		if len(ques) != int(c) {
			return nil, ErrNoEnoughQuestionToMatchRequire
		}
		rate := 0.0
		for _, q := range ques {
			rate += q.MaxDuplicateRate()
		}
		rate /= float64(len(ques))
		if rate > config.RateLimit {
			return nil, ErrRateLimitExceeded
		}
		for i, q := range ques {
			lst, err := sql.Find[List](&f.db, FileTableList, "WHERE ID="+strconv.Itoa(q.ListID))
			if err != nil {
				return nil, err
			}
			quesfile, err := os.Open(q.Path)
			if err != nil {
				return nil, err
			}
			stat, err := quesfile.Stat()
			if err != nil {
				return nil, err
			}
			docq, err := docx.Parse(quesfile, stat.Size())
			if err != nil {
				return nil, err
			}
			docf.AddParagraph().AddText(strconv.Itoa(i+1) + ". (" + lst.Desc + ")")
			docf.AppendFile(docq)
		}
		i++
	}
	return docf, nil
}
