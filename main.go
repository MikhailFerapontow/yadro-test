package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/MikhailFerapontow/yadro-test/pkg/club"
	"github.com/MikhailFerapontow/yadro-test/pkg/models"
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

	_, err = file.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	club, err := club.NewClub(file)
	if err != nil {
		log.Fatalf("error creating club:%s", err)
	}

	club.StartSimulation()
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

	matches := validWorkHours.FindStringSubmatch(scanner.Text())
	startHour := models.NewTime(matches[1], matches[2])
	endHour := models.NewTime(matches[3], matches[4])

	if !endHour.After(startHour) {
		return fmt.Errorf("line %d invalid work hours:\n%s", line, scanner.Text())
	}

	line++
	scanner.Scan()
	if !validProfit.MatchString(scanner.Text()) {
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

		//check if time flows forward
		matches := validTime.FindStringSubmatch(lineText)
		hour, _ := strconv.Atoi(matches[1])
		minutes, _ := strconv.Atoi(matches[2])

		if hour < prevTime.Hour || (hour == prevTime.Hour && minutes < prevTime.Minute) {
			return fmt.Errorf("line %d Time can only flow forward...\n%s", line, scanner.Text())
		}

		prevTime.Hour = hour
		prevTime.Minute = minutes
	}

	return nil
}
