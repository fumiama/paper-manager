package global

import (
	"errors"
	"strconv"

	sql "github.com/FloatTech/sqlite"
	"github.com/fumiama/go-docx"
)

var (
	ErrInvalidGenerateConfig          = errors.New("invalid generate config")
	ErrMajorTooLarge                  = errors.New("major too large")
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
func (f *FileDatabase) GenerateFile(config *GenerateConfig) (*docx.Docx, error) {
	if config == nil || config.Distribution == nil || len(config.Distribution) == 0 {
		return nil, ErrInvalidGenerateConfig
	}
	if len(config.Distribution) > 10 {
		return nil, ErrMajorTooLarge
	}
	docf := docx.NewA4()
	f.mu.RLock()
	defer f.mu.RUnlock()
	for n, c := range config.Distribution {
		if c == 0 {
			continue
		}
		docf.AddParagraph().AddText(string([]rune("一二三四五六七八九十")[c]) + "、" + n).Size("44").Bold()
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
			cond += " Year<=" + strconv.Itoa(int(config.YearStart))
			hasfront = true
		}
		if hasfront {
			cond += " AND"
		}
		cond += " (Type&" + strconv.Itoa(int(config.TypeMask)) + ")!=0"
		ques, err := sql.QueryAll[Question](&f.db,
			"SELECT * FROM "+FileTableQuestion+
				" WHERE FileID IN (SELECT FileID FROM "+
				FileTableFile+cond+
				") ORDER BY RANDOM() limit "+strconv.Itoa(int(c))+";",
		)
		if err != nil {
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
		// TODO: 写入question到docf
	}
	return nil, nil
}
