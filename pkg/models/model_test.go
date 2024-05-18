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

func TestTableStopUsage(t *testing.T) {
	testTable := []struct {
		testName    string
		table       Table
		endTime     Time
		expectedAns Table
	}{
		{
			testName:    "stop usage",
			table:       Table{Occupied: true, StartUse: Time{Hour: 2, Minute: 28}, InUse: Time{Hour: 0, Minute: 00}},
			endTime:     Time{Hour: 4, Minute: 10},
			expectedAns: Table{Occupied: false, InUse: Time{Hour: 1, Minute: 42}},
		},
	}

	for _, test := range testTable {
		t.Run(test.testName, func(t *testing.T) {
			test.table.StopUsage(test.endTime)

			if test.table.InUse != test.expectedAns.InUse {
				t.Errorf("%s: expected %v, got %v", test.testName, test.expectedAns, test.table)
			}
		})
	}
}

func TestTableCalculateProfit(t *testing.T) {
	testTable := []struct {
		testName    string
		table       Table
		tariff      int
		expectedAns int
	}{
		{
			testName:    "calculate profit",
			table:       Table{InUse: Time{Hour: 2, Minute: 28}},
			tariff:      10,
			expectedAns: 30,
		},
		{
			testName:    "full hour",
			table:       Table{InUse: Time{Hour: 2, Minute: 0}},
			tariff:      10,
			expectedAns: 20,
		},
	}

	for _, test := range testTable {
		t.Run(test.testName, func(t *testing.T) {
			ans := test.table.CalculateProfit(test.tariff)
			if ans != test.expectedAns {
				t.Errorf("%s: expected %v, got %v", test.testName, test.expectedAns, ans)
			}
		})
	}
}
