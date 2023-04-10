package global

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"strconv"

	"github.com/corona10/goimagehash"
	"github.com/fumiama/paper-manager/backend/utils"
	"github.com/sirupsen/logrus"
)

// QuestionJSON is the struct representation of File.Questions
type QuestionJSON struct {
	Name   string         `json:"name"`             // Name is name or Question ID
	Points int            `json:"points,omitempty"` // Points is sum of subs' points or self
	Rate   float64        `json:"rate,omitempty"`   // Rate is the avg(non-leaf) or max(leaf) similarity
	Sub    []QuestionJSON `json:"sub,omitempty"`
}

// Delete me and all subs
func (q *QuestionJSON) Delete(f *FileDatabase, istemp bool) {
	if b, err := hex.DecodeString(q.Name); err == nil {
		err = f.DelQuestion(int64(binary.LittleEndian.Uint64(b)), istemp)
		if err != nil {
			logrus.Warnln("[global.QuestionJSON] Delete", q.Name, "err:", err)
		}
	}
	for _, sq := range q.Sub {
		sq.Delete(f, istemp)
	}
}

// DelQuestion 删除问题, 其它问题的 dup 可能会残留有 id, 使用时需要排除
func (f *FileDatabase) DelQuestion(id int64, istemp bool) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if istemp {
		return f.db.Del(FileTableTempQuestion, "WHERE ID="+strconv.FormatInt(id, 10))
	}
	return f.db.Del(FileTableQuestion, "WHERE ID="+strconv.FormatInt(id, 10))
}

type Question struct {
	ID     int64  // ID is the first 8 bytes of the Plain's md5
	Path   string // Path is the question's docx position
	Plain  string // Plain is the plain text of the question (like markdown format)
	Images []byte // Images is json of the image dhash in XML, ex. ['rId1': '1234567890abcdef', ...]
	Vector []byte // Vector is json of {word: freq, ...}
	Dup    []byte // Dup is json of {queid: rate, ...}
}

// GetDuplicateRate calc q & que's dup rate
func (q *Question) GetDuplicateRate(que *Question) (float64, error) {
	v1, v2 := make(map[string]uint8, 64), make(map[string]uint8, 64)
	m1, m2 := make(map[string]string, 64), make(map[string]string, 64)
	if len(q.Images) > 2 {
		err := json.Unmarshal(q.Images, &m1)
		if err != nil {
			return 0, err
		}
	}
	if len(que.Images) > 2 {
		err := json.Unmarshal(que.Images, &m2)
		if err != nil {
			return 0, err
		}
	}
	if len(q.Vector) > 2 {
		err := json.Unmarshal(q.Vector, &v1)
		if err != nil {
			return 0, err
		}
	}
	if len(que.Vector) > 2 {
		err := json.Unmarshal(que.Vector, &v2)
		if err != nil {
			return 0, err
		}
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
	imgdupr := 0.0
	if len(m2) > 0 {
		imgdupr = float64(imgdsts) / float64(len(m2)) / 64.0
	}
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
	if imgdupr > 0 {
		return (8*utils.Similarity(v1space, v2space) + 2*imgdupr) / 10.0, nil
	}
	return utils.Similarity(v1space, v2space), nil
}
