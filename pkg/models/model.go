package models

import (
	"fmt"
	"strconv"
)

type Time struct {
	Hour   int
	Minute int
}

func NewTime(hours, minutes string) Time {
	hour, _ := strconv.Atoi(hours)
	minute, _ := strconv.Atoi(minutes)
	return Time{Hour: hour, Minute: minute}
}

func (t1 *Time) Before(t2 Time) bool {
	return t1.Hour < t2.Hour || (t1.Hour == t2.Hour && t1.Minute < t2.Minute)
}

func (t1 *Time) After(t2 Time) bool {
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
	Occupied bool
	StartUse Time
	InUse    Time
}

// stop usage of the table at time end
// update the inUse time
func (t *Table) StopUsage(end Time) {
	usageTime := end.Subtract(t.StartUse)

	// fmt.Printf("%d, %d\n", usageTime.hour, t.inUse.hour)
	// fmt.Printf("%d, %d\n", usageTime.minute, t.startUse.minute)
	minutesInUse := (usageTime.Hour+t.InUse.Hour)*60 + (usageTime.Minute + t.InUse.Minute)

	// fmt.Printf("min = %d\n", minutesInUse)

	t.InUse = Time{Hour: minutesInUse / 60, Minute: minutesInUse % 60}
	t.Occupied = false
}

func (t *Table) CalculatePrice(tariff int) int {
	sum := t.InUse.Hour * tariff
	if t.InUse.Minute > 0 {
		sum += tariff
	}

	return sum
}

type Client struct {
	InClub   bool
	TableNum int
}