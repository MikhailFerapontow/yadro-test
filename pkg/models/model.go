package models

import (
	"fmt"
	"strconv"
)

type Time struct {
	Hour   int
	Minute int
}

func NewTime(hours, minutes string) (Time, error) {
	hour, err := strconv.Atoi(hours)
	if err != nil {
		return Time{}, fmt.Errorf("error parsing hour: %s", err)
	}

	minute, err := strconv.Atoi(minutes)
	if err != nil {
		return Time{}, fmt.Errorf("error parsing minute: %s", err)
	}

	return Time{Hour: hour, Minute: minute}, nil
}

// Cmp return -1 if time1 is before time2;
// return 0 if time1 is equal time2;
// return 1 if time1 is after time2;
func (t1 *Time) Cmp(t2 Time) int {
	if t1.before(t2) {
		return -1
	}
	if t1.after(t2) {
		return 1
	}
	return 0
}

func (t1 *Time) before(t2 Time) bool {
	return t1.Hour < t2.Hour || (t1.Hour == t2.Hour && t1.Minute < t2.Minute)
}

func (t1 *Time) after(t2 Time) bool {
	return t1.Hour > t2.Hour || (t1.Hour == t2.Hour && t1.Minute > t2.Minute)
}

// Subtract time t2 from t1
func (t1 *Time) Subtract(t2 Time) Time {
	time := (t1.Hour-t2.Hour)*60 + (t1.Minute - t2.Minute)

	return Time{Hour: time / 60, Minute: time % 60}
}

func (t *Time) String() string {
	return fmt.Sprintf("%02d:%02d", t.Hour, t.Minute)
}

type Table struct {
	Occupied  bool
	StartUse  Time
	InUse     Time
	FullHours int // how many full hours the table was in use
}

// stop usage of the table at time end
// update the inUse time
func (t *Table) StopUsage(end Time) {
	usageTime := end.Subtract(t.StartUse)

	t.FullHours += usageTime.Hour
	if usageTime.Minute > 0 {
		t.FullHours++
	}

	minutesInUse := (usageTime.Hour+t.InUse.Hour)*60 + (usageTime.Minute + t.InUse.Minute)

	t.InUse = Time{Hour: minutesInUse / 60, Minute: minutesInUse % 60}
	t.Occupied = false
}

func (t *Table) CalculateProfit(tariff int) int {
	return t.FullHours * tariff
}

type Client struct {
	InClub   bool
	TableNum int
}
