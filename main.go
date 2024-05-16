package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
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
	validTableNum := regexp.MustCompile(`^[0-9]+$`)
	validWorkHours := regexp.MustCompile(`^([0-1]\d|2[0-3]):([0-5]\d) ([0-1]\d|2[0-3]):([0-5]\d)$`)
	validProfit := regexp.MustCompile(`^[0-9]+$`)
	validCommand := regexp.MustCompile(
		`^([0-1]\d|2[0-3]):([0-5]\d) (([1,3,4]) ([a-z0-9_]+)|(2 ([a-z0-9_]+) [0-9]+))$`,
	)

	line := 1
	scanner := bufio.NewScanner(inputFile)

	scanner.Scan()
	if !validTableNum.MatchString(scanner.Text()) {
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

	for {
		line++

		if !scanner.Scan() {
			break
		}

		if !validCommand.MatchString(scanner.Text()) {
			return fmt.Errorf("line %d invalid command:\n%s", line, scanner.Text())
		}
	}

	return nil
}
