package timePlus

import (
	"time"
)

func StartDay(t time.Time) (nt time.Time) {
	nt, _ = time.ParseInLocation("2006-01-02", t.Format("2006-01-02"), t.Location())
	return
}

func EndDay(t time.Time) (nt time.Time) {
	nt, _ = time.ParseInLocation("2006-01-02", t.Format("2006-01-02"), t.Location())
	nt = nt.Add(24*time.Hour - time.Nanosecond)
	return
}
