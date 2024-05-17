package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalf("invalid number of arguments. Need 1, got %d", len(args))
	}

	fileName := args[1]
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err = ValidateInput(file); err != nil {
		log.Fatal(err)
	}
}

func ValidateInput(inputFile *os.File) error {
	validNumOfTable := regexp.MustCompile(`^[0-9]+$`)
	validWorkHours := regexp.MustCompile(`^([0-1]\d|2[0-3]):([0-5]\d) ([0-1]\d|2[0-3]):([0-5]\d)$`)
	validProfit := regexp.MustCompile(`^[0-9]+$`)
	validCommand := regexp.MustCompile(
		`^([0-1]\d|2[0-3]):([0-5]\d) (([1,3,4]) ([a-z0-9_]+)|(2 ([a-z0-9_]+) [0-9]+))$`,
	)

	validTime := regexp.MustCompile(`^([0-1]\d|2[0-3]):([0-5]\d) `)

	line := 1
	scanner := bufio.NewScanner(inputFile)

	scanner.Scan()
	if !validNumOfTable.MatchString(scanner.Text()) {
		return fmt.Errorf("line %d invalid number of tables:\n%s", line, scanner.Text())
	}

	line++
	scanner.Scan()
	if !validWorkHours.MatchString(scanner.Text()) {
		return fmt.Errorf("line %d invalid work hours:\n%s", line, scanner.Text())
	}

	line++
	scanner.Scan()
	if !validProfit.MatchString(scanner.Text()) {
		return fmt.Errorf("line %d invalid hour payment:\n%s", line, scanner.Text())
	}

	prevTime := Time{hour: 0, minute: 0}
	for {
		line++

		if !scanner.Scan() {
			break
		}

		lineText := scanner.Text()

		if !validCommand.MatchString(lineText) {
			return fmt.Errorf("line %d invalid command:\n%s", line, scanner.Text())
		}

		//check if time flows forward
		matches := validTime.FindStringSubmatch(lineText)
		hour, _ := strconv.Atoi(matches[1])
		minutes, _ := strconv.Atoi(matches[2])

		if hour < prevTime.hour || (hour == prevTime.hour && minutes < prevTime.minute) {
			return fmt.Errorf("line %d Time can only flow forward...\n%s", line, scanner.Text())
		}

		prevTime.hour = hour
		prevTime.minute = minutes
	}

	return nil
}
