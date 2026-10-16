package dashboardservice

import (
	"time"

	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

const (
	dayLayout = "2006-01-02"
	// defaultPeriodDays is the window used when from/to are omitted.
	defaultPeriodDays = 30
	maxPeriodDays     = 366
	// mongoTimezone is tehran as a Mongo $dateToString offset. Iran has had
	// no daylight saving since 2022, so a fixed offset is exact.
	mongoTimezone = "+03:30"
)

var tehran = time.FixedZone("IRST", 3*3600+30*60)

// window is one reporting range of whole Tehran-local days.
type window struct {
	from time.Time // first instant of the first day
	to   time.Time // last instant of the last day
	days int
}

func (w window) fromPtr() *time.Time { return &w.from }
func (w window) toPtr() *time.Time   { return &w.to }

func (w window) fromDay() string { return w.from.Format(dayLayout) }
func (w window) toDay() string   { return w.to.Format(dayLayout) }

// dayKeys lists every day of the window, in order.
func (w window) dayKeys() []string {
	keys := make([]string, 0, w.days)
	for d := w.from; d.Before(w.to); d = d.AddDate(0, 0, 1) {
		keys = append(keys, d.Format(dayLayout))
	}
	return keys
}

func newWindow(fromDay, toDay time.Time) window {
	days := int(toDay.Sub(fromDay).Hours()/24) + 1
	return window{
		from: fromDay,
		to:   toDay.AddDate(0, 0, 1).Add(-time.Nanosecond),
		days: days,
	}
}

// resolvePeriod turns the raw query into the current window and the
// same-length window right before it.
func resolvePeriod(fromRaw, toRaw string, now time.Time) (window, window, error) {
	const op = "dashboardservice.resolvePeriod"
	invalid := func() error {
		return richerror.New(op).WithKind(richerror.KindInvalid).WithMessageKey("dashboard.invalid_period")
	}

	today := now.In(tehran)
	toDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, tehran)
	if toRaw != "" {
		t, err := time.ParseInLocation(dayLayout, toRaw, tehran)
		if err != nil {
			return window{}, window{}, invalid()
		}
		toDay = t
	}

	fromDay := toDay.AddDate(0, 0, -(defaultPeriodDays - 1))
	if fromRaw != "" {
		t, err := time.ParseInLocation(dayLayout, fromRaw, tehran)
		if err != nil {
			return window{}, window{}, invalid()
		}
		fromDay = t
	}

	if fromDay.After(toDay) {
		return window{}, window{}, invalid()
	}
	cur := newWindow(fromDay, toDay)
	if cur.days > maxPeriodDays {
		return window{}, window{}, invalid()
	}

	prevTo := fromDay.AddDate(0, 0, -1)
	prevFrom := fromDay.AddDate(0, 0, -cur.days)
	return cur, newWindow(prevFrom, prevTo), nil
}
