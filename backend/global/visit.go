package global

import (
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// MonthlyAPIVisit counts the api visit history
type MonthlyAPIVisit struct {
	YM    uint32 // YM is yyyymm
	Count uint32 // visit count this mounth
}

// VisitAPI increases count of this mounth by 1
func (u *UserDatabase) VisitAPI() {
	now := time.Now()
	ym := uint32(now.Year())*100 + uint32(now.Month())
	var v MonthlyAPIVisit
	u.mu.Lock()
	defer u.mu.Unlock()
	_ = u.db.Find(UserTableMonthlyAPIVisit, &v, "WHERE YM="+strconv.FormatUint(uint64(ym), 10))
	v.YM = ym
	v.Count++
	err := u.db.Insert(UserTableMonthlyAPIVisit, &v)
	if err != nil {
		logrus.Warnln("[global.user] insert visit error:", err)
	}
}

// GetAnnualAPIVisitCount get the latest 12 mounths' count
func (u *UserDatabase) GetAnnualAPIVisitCount() (cnts [12]uint32) {
	var v MonthlyAPIVisit
	var yms [12]uint32
	now := time.Now()
	y100 := uint32(now.Year()) * 100
	py100 := uint32(now.Year()-1) * 100
	nm := int(now.Month())
	for i := 0; i < nm; i++ {
		yms[i] = y100 + uint32(i+1)
	}
	for i := nm; i < 12; i++ {
		yms[i] = py100 + uint32(i+1)
	}
	u.mu.RLock()
	defer u.mu.RUnlock()
	i := 0
	for _, ym := range yms {
		_ = u.db.Find(UserTableMonthlyAPIVisit, &v, "WHERE YM="+strconv.FormatUint(uint64(ym), 10))
		cnts[i] = v.Count
		i++
		v.Count = 0
	}
	return
}
