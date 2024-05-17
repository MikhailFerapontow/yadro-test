package main

import (
	"fmt"
	"strconv"
)

type Time struct {
	hour   int
	minute int
}

func NewTime(hours, minutes string) Time {
	hour, _ := strconv.Atoi(hours)
	minute, _ := strconv.Atoi(minutes)
	return Time{hour: hour, minute: minute}
}

func (t1 *Time) Before(t2 Time) bool {
	return t1.hour < t2.hour || (t1.hour == t2.hour && t1.minute < t2.minute)
}

func (t1 *Time) After(t2 Time) bool {
	return t1.hour > t2.hour || (t1.hour == t2.hour && t1.minute > t2.minute)
}

// Subtract time t2 from t1
func (t1 *Time) Subtract(t2 Time) Time {
	time := (t1.hour-t2.hour)*60 + (t1.minute - t2.minute)

	return Time{hour: time / 60, minute: time % 60}
}

func (t *Time) String() string {
	return fmt.Sprintf("%02d:%02d", t.hour, t.minute)
}

type Table struct {
	occupied bool
	startUse Time
	inUse    Time
}

// stop usage of the table at time end
// update the inUse time
func (t *Table) StopUsage(end Time) {
	usageTime := end.Subtract(t.startUse)

	// fmt.Printf("%d, %d\n", usageTime.hour, t.inUse.hour)
	// fmt.Printf("%d, %d\n", usageTime.minute, t.startUse.minute)
	minutesInUse := (usageTime.hour+t.inUse.hour)*60 + (usageTime.minute + t.inUse.minute)

	// fmt.Printf("min = %d\n", minutesInUse)

	t.inUse = Time{hour: minutesInUse / 60, minute: minutesInUse % 60}
	t.occupied = false
}

func (t *Table) CalculatePrice(tariff int) int {
	sum := t.inUse.hour * tariff
	if t.inUse.minute > 0 {
		sum += tariff
	}

	return sum
}

type Client struct {
	inClub   bool
	tableNum int
}
