package main

type Time struct {
	hour   int
	minute int
}

// Subtract time t2 from t1
func (t1 *Time) Sub(t2 Time) Time {
	return Time{hour: t1.hour - t2.hour, minute: t1.minute - t2.minute}
}
