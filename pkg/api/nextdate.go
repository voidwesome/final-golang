package api

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func containsInt(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// NextDate implements all rules: d, y, w, m
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if strings.TrimSpace(repeat) == "" {
		return "", errors.New("empty repeat rule")
	}
	start, err := time.Parse(DateLayout, dstart)
	if err != nil {
		return "", fmt.Errorf("bad start date: %w", err)
	}

	parts := strings.Fields(repeat)
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid d format")
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil || n <= 0 || n > 400 {
			return "", errors.New("invalid d interval")
		}
		t := start
		for {
			t = t.AddDate(0, 0, n)
			if afterNow(t, now) {
				return t.Format(DateLayout), nil
			}
		}
	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid y format")
		}
		t := start
		for {
			t = t.AddDate(1, 0, 0)
			if afterNow(t, now) {
				return t.Format(DateLayout), nil
			}
		}
	case "w":
		// w <1..7,[comma]>
		if len(parts) != 2 {
			return "", errors.New("invalid w format")
		}
		days, err := parseWeekdays(parts[1])
		if err != nil {
			return "", err
		}
		// roll from start forward until > now
		t := start
		// fast-forward to first > now
		if !afterNow(t, now) {
			t = now
		}
		for i := 0; i < 800; i++ {
			t = t.AddDate(0, 0, 1)
			if afterNow(t, now) && containsWeekday(days, weekdayRU(t.Weekday())) {
				return t.Format(DateLayout), nil
			}
		}
		return "", errors.New("cannot find next date for w-rule")
	case "m":
		// m <days list> [months list]
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("invalid m format")
		}
		mlist := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		var err error
		days, err := parseMonthDays(parts[1])
		if err != nil {
			return "", err
		}
		if len(parts) == 3 {
			mlist, err = parseMonths(parts[2])
			if err != nil {
				return "", err
			}
		}
		sort.Ints(mlist)

		// algorithm:
		t := start
		if !afterNow(t, now) {
			t = now
		}
		// try from now+1 day onward
		for i := 0; i < 5000; i++ {
			t = t.AddDate(0, 0, 1)
			if !afterNow(t, now) {
				continue
			}
			if !containsInt(mlist, int(t.Month())) {
				continue
			}
			if containsInt(days, monthDayOrNegative(t)) {
				return t.Format(DateLayout), nil
			}
		}
		return "", errors.New("cannot find next date for m-rule")
	default:
		return "", errors.New("unsupported repeat format")
	}
}

func weekdayRU(w time.Weekday) int {
	// 1 - Monday ... 7 - Sunday
	if w == time.Sunday {
		return 7
	}
	return int(w)
}

func parseWeekdays(s string) ([]int, error) {
	ss := strings.Split(s, ",")
	days := make([]int, 0, len(ss))
	for _, x := range ss {
		v, err := strconv.Atoi(strings.TrimSpace(x))
		if err != nil || v < 1 || v > 7 {
			return nil, errors.New("invalid weekday in w-rule")
		}
		days = append(days, v)
	}
	return days, nil
}

func containsWeekday(arr []int, v int) bool {
	for _, x := range arr {
		if x == v {
			return true
		}
	}
	return false
}

func parseMonthDays(s string) ([]int, error) {
	ss := strings.Split(s, ",")
	days := make([]int, 0, len(ss))
	for _, x := range ss {
		v, err := strconv.Atoi(strings.TrimSpace(x))
		if err != nil {
			return nil, errors.New("invalid month day")
		}
		if !(v >= 1 && v <= 31) && v != -1 && v != -2 {
			return nil, errors.New("invalid month day value")
		}
		days = append(days, v)
	}
	return days, nil
}

func parseMonths(s string) ([]int, error) {
	ss := strings.Split(s, ",")
	months := make([]int, 0, len(ss))
	for _, x := range ss {
		v, err := strconv.Atoi(strings.TrimSpace(x))
		if err != nil || v < 1 || v > 12 {
			return nil, errors.New("invalid month in m-rule")
		}
		months = append(months, v)
	}
	return months, nil
}

func monthDayOrNegative(t time.Time) int {
	// returns day number; if it's last day -> -1, if pre-last -> -2
	year, month, day := t.Date()
	lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, t.Location()).Day()
	if day == lastDay {
		return -1
	}
	if day == lastDay-1 {
		return -2
	}
	return day
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	nowStr := q.Get("now")
	date := q.Get("date")
	rep := q.Get("repeat")

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateLayout, nowStr)
		if err != nil {
			writeError(w, fmt.Errorf("bad now: %w", err), http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, date, rep)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]string{"date": next})
}
