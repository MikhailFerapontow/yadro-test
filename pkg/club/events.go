package club

import "fmt"

// Prints line to std::out with newline
func PrintLine(line string) {
	fmt.Printf("%s\n", line)
}

// Prints error message NotOpenYer to std::out
// with newline
func NotOpenYet(time string) {
	fmt.Printf("%s 13 NotOpenYet\n", time)
}

// Prints error message YouShallNotPass to std::out
// with newline
func YouShallNotPass(time string) {
	fmt.Printf("%s 13 YouShallNotPass\n", time)
}

// Prints error message ClientUnknown to std::out
// with newline
func ClientUnknown(time string) {
	fmt.Printf("%s 13 ClientUnknown\n", time)
}

// Prints error message ClientLeft to std::out
// with newline
func PlaceIsBusy(time string) {
	fmt.Printf("%s 13 PlaceIsBusy\n", time)
}

// Prints error message ICanWaitNoLonger to std::out
// with newline
func ICanWaitNoLonger(time string) {
	fmt.Printf("%s 13 ICanWaitNoLonger!\n", time)
}

// Prints output event with ID 11 to std::out
// with newline
func ClientLeft(time, clientName string) {
	fmt.Printf("%s 11 %s\n", time, clientName)
}

// Prints output event with ID 12 to std::out
// with newline
func ClientTakeTable(time, clientName string, tableIdx int) {
	fmt.Printf("%s 12 %s %d\n", time, clientName, tableIdx)
}
