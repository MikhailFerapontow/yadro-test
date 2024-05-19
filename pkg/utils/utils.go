package utils

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/MikhailFerapontow/yadro-test/pkg/models"
)

// Validates input file
// Returns error if invalid
func ValidateInput(inputFile io.Reader) error {
	validNum := regexp.MustCompile(`^[0-9]+$`)
	validWorkHours := regexp.MustCompile(`^([0-1]\d|2[0-3]):([0-5]\d) ([0-1]\d|2[0-3]):([0-5]\d)$`)
	validCommand := regexp.MustCompile(
		`^([0-1]\d|2[0-3]):([0-5]\d) (([1,3,4]) ([a-z0-9_]+)|(2 ([a-z0-9_]+) [0-9]+))$`,
	)
	validCommandType2 := regexp.MustCompile(`^([0-1]\d|2[0-3]):([0-5]\d) 2 ([a-z0-9_]+) ([0-9]+)$`)

	validTime := regexp.MustCompile(`^([0-1]\d|2[0-3]):([0-5]\d)`)

	line := 1
	scanner := bufio.NewScanner(inputFile)

	scanner.Scan()
	if !validNum.MatchString(scanner.Text()) {
		return fmt.Errorf("line %d invalid number of tables:\n%s", line, scanner.Text())
	}
	matches := validNum.FindStringSubmatch(scanner.Text())
	numOfTables, _ := strconv.Atoi(matches[0])

	line++
	scanner.Scan()
	if !validWorkHours.MatchString(scanner.Text()) {
		return fmt.Errorf("line %d invalid work hours:\n%s", line, scanner.Text())
	}

	matches = validWorkHours.FindStringSubmatch(scanner.Text())
	startHour, _ := models.NewTime(matches[1], matches[2]) // regex proves it's an int
	endHour, _ := models.NewTime(matches[3], matches[4])

	if endHour.Cmp(startHour) != 1 {
		return fmt.Errorf("line %d invalid work hours:\n%s", line, scanner.Text())
	}

	line++
	scanner.Scan()
	if !validNum.MatchString(scanner.Text()) {
		return fmt.Errorf("line %d invalid hour payment:\n%s", line, scanner.Text())
	}

	prevTime := models.Time{Hour: 0, Minute: 0}
	for {
		line++

		if !scanner.Scan() {
			break
		}

		lineText := scanner.Text()

		if !validCommand.MatchString(lineText) {
			return fmt.Errorf("line %d invalid command:\n%s", line, scanner.Text())
		}

		//check if client takes availiable table
		commandType2 := validCommandType2.MatchString(lineText)
		if commandType2 {
			matches = validCommandType2.FindStringSubmatch(lineText)
			tableNum, _ := strconv.Atoi(matches[4])
			if tableNum > numOfTables || tableNum == 0 {
				return fmt.Errorf("line %d Table number out of range:\n%s", line, scanner.Text())
			}
		}

		//check if time flows forward
		matches = validTime.FindStringSubmatch(lineText)
		curTime, _ := models.NewTime(matches[1], matches[2])

		if prevTime.Cmp(curTime) == 1 {
			return fmt.Errorf("line %d Time can only flow forward...\n%s", line, scanner.Text())
		}

		prevTime = curTime
	}

	return nil
}
