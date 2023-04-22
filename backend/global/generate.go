package global

import "errors"

var (
	ErrInvalidGenerateConfig = errors.New("invalid generate config")
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
func (f *FileDatabase) GenerateFile(config *GenerateConfig) ([]byte, error) {
	if config == nil || config.Distribution == nil {
		return nil, ErrInvalidGenerateConfig
	}
	return nil, nil
}
