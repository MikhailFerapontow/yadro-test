package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Club struct {
	scanner *bufio.Scanner

	queue  chan string
	client map[string]Client
	table  map[int]Table

	tariff    int
	startHour Time
	endHour   Time
}

func NewClub(file *os.File) (*Club, error) {
	_, err := file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	numOfTable, _ := strconv.Atoi(scanner.Text())

	scanner.Scan()
	workHours := strings.FieldsFunc(scanner.Text(), func(c rune) bool {
		return c == ':' || c == ' '
	})

	scanner.Scan()
	tariff, _ := strconv.Atoi(scanner.Text())

	table := make(map[int]Table, numOfTable)
	for i := 0; i < numOfTable; i++ {
		table[i] = Table{}
	}

	return &Club{
		queue:     make(chan string, numOfTable),
		client:    make(map[string]Client),
		table:     table,
		tariff:    tariff,
		startHour: NewTime(workHours[0], workHours[1]),
		endHour:   NewTime(workHours[2], workHours[3]),
		scanner:   scanner,
	}, nil
}

func (c *Club) StartSimulation() {
	fmt.Printf("%s\n", c.startHour.String())

	for {
		if !c.scanner.Scan() {
			break
		}

		tokens := strings.FieldsFunc(c.scanner.Text(), func(c rune) bool {
			return c == ':' || c == ' '
		})

		time := NewTime(tokens[0], tokens[1])
		client := tokens[3]

		switch command := tokens[2]; command {
		case "1":
			fmt.Printf("%s\n", c.scanner.Text())

			c.clientArrive(time, client)
		case "2":
			fmt.Printf("%s\n", c.scanner.Text())

			tableId, _ := strconv.Atoi(tokens[4])
			c.clientTakeTable(time, client, tableId)
		case "3":
			fmt.Printf("%s\n", c.scanner.Text())

			c.clientWait(time, client)
		case "4":
			fmt.Printf("%s\n", c.scanner.Text())

			c.clientLeave(time, client)
		}
	}

	c.CloseClub()
	fmt.Printf("%s\n", c.endHour.String())
	c.TableInfo()
}

func (c *Club) clientArrive(time Time, client string) {
	//check working hours
	if time.Before(c.startHour) || time.After(c.endHour) {
		fmt.Printf("%s 13 NotOpenYet\n", time.String())
		return
	}

	//check if client is in club
	val, ok := c.client[client]
	if !ok {
		c.client[client] = Client{inClub: true, tableNum: 0}
		return
	}

	if val.inClub {
		fmt.Printf("%s 13 YouShallNotPass\n", time.String())
	} else {
		c.client[client] = Client{inClub: true, tableNum: 0}
	}
}

func (c *Club) clientTakeTable(time Time, client string, tableId int) {
	curClient, ok := c.client[client]
	if !ok || !curClient.inClub {
		fmt.Printf("%s 13 ClientUnknown\n", time.String())
		return
	}

	curTable := c.table[tableId-1]
	if curTable.occupied {
		fmt.Printf("%s 13 PlaceIsBusy\n", time.String())
		return
	}

	c.client[client] = Client{inClub: true, tableNum: tableId}
	c.table[tableId-1] = Table{occupied: true, startUse: time}
}

func (c *Club) clientWait(time Time, client string) {
	//check if free tables exist
	for _, t := range c.table {
		if !t.occupied {
			fmt.Printf("%s 13 ICanWaitNoLonger!\n", time.String())
			return
		}
	}

	select {
	case c.queue <- client:
	default:
		fmt.Printf("%s 11 %s\n", time.String(), client)
	}
}

func (c *Club) clientLeave(time Time, client string) {
	curClient, ok := c.client[client]
	if !ok || !curClient.inClub {
		fmt.Printf("%s 13 ClientUnknown\n", time.String())
		return
	}

	tableIdx := curClient.tableNum - 1
	t := c.table[tableIdx]
	t.StopUsage(time)

	c.client[client] = Client{inClub: false, tableNum: 0}

	select {
	case queueClient := <-c.queue:
		c.client[queueClient] = Client{inClub: true, tableNum: tableIdx + 1}
		c.table[tableIdx] = Table{occupied: true, startUse: time, inUse: t.inUse}
		fmt.Printf("%s 12 %s %d\n", time.String(), queueClient, tableIdx+1)
	default:
		c.table[tableIdx] = t
	}
}

func (c *Club) CloseClub() {
	clients := make([]string, 0)

	for clientName, client := range c.client {
		if client.inClub {
			tableIdx := client.tableNum - 1
			val, ok := c.table[tableIdx]
			if ok {
				val.StopUsage(c.endHour)
				c.table[tableIdx] = val
			}

			clients = append(clients, clientName)
		}
	}

	sort.Slice(clients, func(i, j int) bool {
		return clients[i] < clients[j]
	})

	for _, client := range clients {
		fmt.Printf("%s 11 %s\n", c.endHour.String(), client)
	}
}

func (c *Club) TableInfo() {
	for i, t := range c.table {
		fmt.Printf("%d %d %s\n", i+1, t.CalculatePrice(c.tariff), t.inUse.String())
	}
}
