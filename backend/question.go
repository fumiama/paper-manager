package backend

import (
	"encoding/json"

	"github.com/fumiama/paper-manager/backend/global"
)

type question struct {
	Count int    `json:"count"`
	Point int    `json:"point"`
	Name  string `json:"name"`
}

type duplication struct {
	Percent int    `json:"percent"`
	Name    string `json:"name"`
}

func parseFileQuestions(qb []byte) ([]question, []duplication, error) {
	ques := make([]global.QuestionJSON, 0, 16)
	qs := make([]question, 0, 16)
	ds := make([]duplication, 0, 16)
	err := json.Unmarshal(qb, &qs)
	if err != nil {
		return nil, nil, err
	}
	for _, q := range ques {
		qs = append(qs, question{
			Count: len(q.Sub),
			Point: q.Points,
			Name:  q.Name,
		})
		// TODO: use heap to get top 10 ds
	}
	return nil, nil, nil
}
