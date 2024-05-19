package main

import (
	"log"
	"os"

	"github.com/MikhailFerapontow/yadro-test/pkg/club"
	"github.com/MikhailFerapontow/yadro-test/pkg/utils"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalf("invalid number of arguments. Need 1, got %d", len(args)-1)
	}

	fileName := args[1]
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err = utils.ValidateInput(file); err != nil {
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
