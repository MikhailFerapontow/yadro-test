package club

import (
	"strings"
	"testing"

	"github.com/MikhailFerapontow/yadro-test/pkg/models"
)

func ExampleClub() {
	r := strings.NewReader(
		"1\n" +
			"09:00 19:00\n" +
			"10\n" +
			"08:48 1 client1\n" +
			"09:41 1 client1\n" +
			"10:25 2 client1 1\n" +
			"10:30 2 client3 1\n" +
			"12:33 4 client1\n",
	)

	c, _ := NewClub(r)
	c.StartSimulation()
	// Output: 09:00
	// 08:48 1 client1
	// 08:48 13 NotOpenYet
	// 09:41 1 client1
	// 10:25 2 client1 1
	// 10:30 2 client3 1
	// 10:30 13 ClientUnknown
	// 12:33 4 client1
	// 19:00
	// 1 30 02:08
}

func TestClub(t *testing.T) {
	testTable := []struct {
		testName    string
		input       string
		expectedAns []models.Table
	}{
		{
			testName: "test club",
			input: "3\n" +
				"09:00 19:00\n" +
				"10\n" +
				"08:48 1 client1\n" +
				"09:41 1 client1\n" +
				"09:48 1 client2\n" +
				"09:52 3 client1\n" +
				"09:54 2 client1 1\n" +
				"10:25 2 client2 2\n" +
				"10:58 1 client3\n" +
				"10:59 2 client3 3\n" +
				"11:30 1 client4\n" +
				"11:35 2 client4 2\n" +
				"11:45 3 client4\n" +
				"12:33 4 client1\n" +
				"12:43 4 client2\n" +
				"15:52 4 client4\n",
			expectedAns: []models.Table{
				{
					Occupied: false,
					StartUse: models.Time{Hour: 12, Minute: 33},
					InUse:    models.Time{Hour: 5, Minute: 58},
				},
				{
					Occupied: false,
					StartUse: models.Time{Hour: 10, Minute: 25},
					InUse:    models.Time{Hour: 2, Minute: 18},
				},
				{
					Occupied: false,
					StartUse: models.Time{Hour: 10, Minute: 59},
					InUse:    models.Time{Hour: 8, Minute: 01},
				},
			},
		},
		{
			testName: "work all day",
			input: "3\n" +
				"09:00 19:00\n" +
				"10\n" +
				"08:48 1 client1\n" +
				"09:00 1 client1\n" +
				"09:00 1 client2\n" +
				"09:00 3 client1\n" +
				"09:00 2 client1 1\n" +
				"09:00 2 client2 2\n" +
				"09:00 1 client3\n" +
				"09:00 2 client3 3\n" + //заняты все столы
				"10:30 1 client4\n" +
				"10:30 3 client4\n" +
				"10:30 1 client5\n" +
				"10:30 3 client5\n" +
				"10:30 1 client6\n" +
				"10:30 3 client6\n" + //заполнилась очередь
				"10:30 1 client7\n" +
				"10:30 3 client7\n" + // этот клиент ушёл
				"12:33 4 client1\n" +
				"12:43 4 client2\n" +
				"15:52 4 client3\n" + // их место заняли
				"18:00 4 client5\n", // его место занять некому
			expectedAns: []models.Table{
				{
					Occupied: false,
					StartUse: models.Time{Hour: 12, Minute: 33},
					InUse:    models.Time{Hour: 10, Minute: 00},
				},
				{
					Occupied: false,
					StartUse: models.Time{Hour: 12, Minute: 43},
					InUse:    models.Time{Hour: 9, Minute: 00},
				},
				{
					Occupied: false,
					StartUse: models.Time{Hour: 15, Minute: 52},
					InUse:    models.Time{Hour: 10, Minute: 00},
				},
			},
		},
	}

	for _, test := range testTable {
		t.Run(test.testName, func(t *testing.T) {
			r := strings.NewReader(test.input)
			c, _ := NewClub(r)
			c.StartSimulation()

			if len(c.table) != len(test.expectedAns) {
				t.Errorf("%s: expected len %v, got len %v", test.testName, len(test.expectedAns), len(c.table))
			}

			gotTable := make([]models.Table, len(c.table))
			for i, table := range c.table {
				gotTable[i] = table
			}

			for i, table := range test.expectedAns {
				if gotTable[i] != table {
					t.Errorf("%s: expected %v, got %v", test.testName, table, gotTable[i])
				}
			}
		})
	}
}
