package backend

import (
	"container/heap"
	"encoding/json"
	"math"
	"net/http"
	"strconv"

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

type duplications []duplication

func (d *duplications) Len() int {
	return len(*d)
}

// Less is actually more for a big-top heap
func (d *duplications) Less(i, j int) bool {
	return (*d)[i].Percent > (*d)[j].Percent
}

func (d *duplications) Swap(i, j int) {
	(*d)[i], (*d)[j] = (*d)[j], (*d)[i]
}

func (d *duplications) Push(x any) {
	*d = append(*d, x.(duplication))
}

func (d *duplications) Pop() any {
	if d.Len() == 0 {
		return nil
	}
	i := d.Len() - 1
	x := (*d)[i]
	*d = (*d)[:i]
	return x
}

func parseFileQuestions(qb []byte, istemp, getdup bool) ([]question, []duplication, float64, error) {
	ques := make([]global.QuestionJSON, 0, 16)
	qs := make([]question, 0, 16)
	err := json.Unmarshal(qb, &ques)
	if err != nil {
		return nil, nil, 0, err
	}
	dhp := (*duplications)(nil)
	if getdup {
		dh := make(duplications, 0, 16)
		heap.Init(&dh)
		dhp = &dh
	}
	sum := 0.0
	cnt := 0
	for _, q := range ques {
		qs = append(qs, question{
			Count: len(q.Sub),
			Point: q.Points,
			Name:  q.Name,
		})
		for i, subq := range q.Sub {
			qstruct, err := global.FileDB.GetQuestionByHex(subq.Name, istemp)
			if err != nil {
				continue
			}
			p := qstruct.MaxDuplicateRate()
			if getdup {
				heap.Push(dhp, duplication{
					Percent: int(math.Round(p * 100)),
					Name:    q.Name + "." + strconv.Itoa(i+1),
				})
			}
			sum += p
			cnt++
		}
	}
	if getdup {
		i := dhp.Len()
		ds := make([]duplication, 10)
		if i > 10 {
			i = 10
		} else {
			for j := i; j < 10; j++ {
				ds[j] = duplication{Name: "N/A"}
			}
		}
		for i--; i >= 0; i-- {
			ds[i] = heap.Pop(dhp).(duplication)
		}
		return qs, ds, sum / float64(cnt), nil
	}

	return qs, nil, sum / float64(cnt), nil
}

// getQuestionDupFromPaper returns rate, error
func getQuestionDupFromPaper(que global.QuestionJSON, ques []global.QuestionJSON, istemp bool) (float64, error) {
	myqstruct, err := global.FileDB.GetQuestionByHex(que.Name, istemp)
	if err != nil {
		return -1, err
	}
	maxp := 0.0
	for _, q := range ques {
		for _, subq := range q.Sub {
			qstruct, err := global.FileDB.GetQuestionByHex(subq.Name, istemp)
			if err != nil {
				continue
			}
			p, err := myqstruct.GetDuplicateRate(&qstruct)
			if err != nil {
				continue
			}
			if maxp < p {
				maxp = p
			}
		}
	}
	return maxp, nil
}

func init() {
	apimap["/api/getMajors"] = &apihandler{"GET", func(w http.ResponseWriter, r *http.Request) {
		user := usertokens.Get(r.Header.Get("Authorization"))
		if user == nil {
			writeresult(w, codeError, nil, errInvalidToken.Error(), typeError)
			return
		}
		majs := global.FileDB.GetMajors()
		type majret struct {
			Name string
		}
		majrets := make([]majret, len(majs))
		for i, s := range majs {
			majrets[i].Name = s
		}
		writeresult(w, codeSuccess, &majrets, messageOk, typeSuccess)
	}}
}
