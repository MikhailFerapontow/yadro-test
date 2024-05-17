package models

import "testing"

func TestTimeCmp(t *testing.T) {
	testTable := []struct {
		testName    string
		t1          Time
		t2          Time
		expectedAns int
	}{
		{"time1 < time2", Time{Hour: 1, Minute: 0}, Time{Hour: 2, Minute: 0}, -1},
		{"time1 > time2", Time{Hour: 2, Minute: 0}, Time{Hour: 1, Minute: 0}, 1},
		{"time1 == time2", Time{Hour: 1, Minute: 0}, Time{Hour: 1, Minute: 0}, 0},
	}

	for _, test := range testTable {
		t.Run(test.testName, func(t *testing.T) {
			ans := test.t1.Cmp(test.t2)
			if ans != test.expectedAns {
				t.Errorf("%s: expected %v, got %v", test.testName, test.expectedAns, ans)
			}
		})
	}
}

func TestTimeSubtract(t *testing.T) {
	testTable := []struct {
		testName    string
		t1          Time
		t2          Time
		expectedAns Time
	}{
		{"time1 > time2", Time{Hour: 2, Minute: 25}, Time{Hour: 1, Minute: 20}, Time{Hour: 1, Minute: 5}},
		{"time1 == time2", Time{Hour: 1, Minute: 0}, Time{Hour: 1, Minute: 0}, Time{Hour: 0, Minute: 0}},
		{"more minutes", Time{Hour: 1, Minute: 10}, Time{Hour: 0, Minute: 20}, Time{Hour: 0, Minute: 50}},
	}

	for _, test := range testTable {
		t.Run(test.testName, func(t *testing.T) {
			ans := test.t1.Subtract(test.t2)
			if ans != test.expectedAns {
				t.Errorf("%s: expected %v, got %v", test.testName, test.expectedAns, ans)
			}
		})
	}
}
