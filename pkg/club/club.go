package club

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/MikhailFerapontow/yadro-test/pkg/models"
)

const (
	ClientArrive     = "1"
	ClientTakesTable = "2"
	ClientWait       = "3"
	ClientLeave      = "4"
)

type Club struct {
	scanner *bufio.Scanner

	queue  chan string
	client map[string]models.Client
	table  map[int]models.Table

	tariff    int
	startHour models.Time
	endHour   models.Time
}

func NewClub(file io.Reader) (*Club, error) {
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	numOfTable, _ := strconv.Atoi(scanner.Text())

	scanner.Scan()
	workHours := strings.FieldsFunc(scanner.Text(), func(c rune) bool {
		return c == ':' || c == ' '
	})

	startHour, _ := models.NewTime(workHours[0], workHours[1])
	endHour, _ := models.NewTime(workHours[2], workHours[3])

	scanner.Scan()
	tariff, _ := strconv.Atoi(scanner.Text())

	table := make(map[int]models.Table, numOfTable)
	for i := 0; i < numOfTable; i++ {
		table[i] = models.Table{}
	}

	return &Club{
		queue:     make(chan string, numOfTable),
		client:    make(map[string]models.Client),
		table:     table,
		tariff:    tariff,
		startHour: startHour,
		endHour:   endHour,
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

		time, _ := models.NewTime(tokens[0], tokens[1])
		client := tokens[3]

		PrintLine(c.scanner.Text())

		switch command := tokens[2]; command {
		case ClientArrive:

			c.clientArrive(time, client)
		case ClientTakesTable:

			tableId, _ := strconv.Atoi(tokens[4])
			c.clientTakeTable(time, client, tableId)
		case ClientWait:

			c.clientWait(time, client)
		case ClientLeave:

			c.clientLeave(time, client)
		}
	}

	c.CloseClub()
	fmt.Printf("%s\n", c.endHour.String())
	c.TableInfo()
}

// client arrives to club
func (c *Club) clientArrive(time models.Time, clientName string) {
	//check working hours
	if time.Cmp(c.startHour) == -1 || time.Cmp(c.endHour) == 1 {
		NotOpenYet(time.String())
		return
	}

	//check if client is in club
	client, ok := c.client[clientName]
	if !ok {
		c.client[clientName] = models.Client{InClub: true, TableNum: 0}
		return
	}

	if client.InClub {
		YouShallNotPass(time.String())
	} else {
		c.client[clientName] = models.Client{InClub: true, TableNum: 0}
	}
}

// client takes empty table
func (c *Club) clientTakeTable(time models.Time, clientName string, tableId int) {
	client, ok := c.client[clientName]
	if !ok || !client.InClub {
		ClientUnknown(time.String())
		return
	}

	curTable := c.table[tableId-1]
	if curTable.Occupied {
		PlaceIsBusy(time.String())
		return
	}

	c.client[clientName] = models.Client{InClub: true, TableNum: tableId}
	c.table[tableId-1] = models.Table{Occupied: true, StartUse: time}
}

// client waits to take table
func (c *Club) clientWait(time models.Time, clientName string) {
	client, ok := c.client[clientName]
	if !ok || !client.InClub {
		ClientUnknown(time.String())
		return
	}

	//check if free tables exist
	for _, t := range c.table {
		if !t.Occupied {
			ICanWaitNoLonger(time.String())
			return
		}
	}

	// fill client queue
	select {
	case c.queue <- clientName:
	default:
		//if queue is full
		ClientLeft(time.String(), clientName)
	}
}

// client leaves the club
func (c *Club) clientLeave(time models.Time, clientName string) {
	client, ok := c.client[clientName]
	if !ok || !client.InClub {
		ClientUnknown(time.String())
		return
	}

	tableIdx := client.TableNum - 1
	table, ok := c.table[tableIdx]
	if !ok {
		c.client[clientName] = models.Client{InClub: false, TableNum: 0}
		return
	}
	table.StopUsage(time)

	c.client[clientName] = models.Client{InClub: false, TableNum: 0}

	// take first client from queue
	select {
	case queueClient := <-c.queue:
		c.client[queueClient] = models.Client{InClub: true, TableNum: tableIdx + 1}
		c.table[tableIdx] = models.Table{Occupied: true, StartUse: time, InUse: table.InUse, FullHours: table.FullHours}
		ClientTakeTable(time.String(), queueClient, tableIdx+1)
	default:
		c.table[tableIdx] = table
	}
}

func (c *Club) CloseClub() {
	clients := make([]string, 0)

	for clientName, client := range c.client {
		if client.InClub {
			tableIdx := client.TableNum - 1
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
		ClientLeft(c.endHour.String(), client)
	}
}

func (c *Club) TableInfo() {
	output := make([]string, len(c.table))

	for i, table := range c.table {
		output[i] = fmt.Sprintf("%d %d %s", i+1, table.CalculateProfit(c.tariff), table.InUse.String())
	}

	for _, table := range output {
		fmt.Println(table)
	}
}
