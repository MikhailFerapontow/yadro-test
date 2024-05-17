package club

import (
	"strings"
)

func ExampleClub() {
	r := strings.NewReader(
		"1\n" +
			"09:00 19:00\n" +
			"10\n" +
			"08:48 1 client1\n" +
			"09:41 1 client1\n" +
			"10:25 2 client1 1\n" +
			"12:33 4 client1\n",
	)

	c, _ := NewClub(r)
	c.StartSimulation()
	// Output: 09:00
	// 08:48 1 client1
	// 08:48 13 NotOpenYet
	// 09:41 1 client1
	// 10:25 2 client1 1
	// 12:33 4 client1
	// 19:00
	// 1 30 02:08
}
