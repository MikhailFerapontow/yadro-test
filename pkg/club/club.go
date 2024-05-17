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
		startHour: models.NewTime(workHours[0], workHours[1]),
		endHour:   models.NewTime(workHours[2], workHours[3]),
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

		time := models.NewTime(tokens[0], tokens[1])
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

// client arrives to club
func (c *Club) clientArrive(time models.Time, client string) {
	//check working hours
	if time.Cmp(c.startHour) == -1 || time.Cmp(c.endHour) == 1 {
		fmt.Printf("%s 13 NotOpenYet\n", time.String())
		return
	}

	//check if client is in club
	val, ok := c.client[client]
	if !ok {
		c.client[client] = models.Client{InClub: true, TableNum: 0}
		return
	}

	if val.InClub {
		fmt.Printf("%s 13 YouShallNotPass\n", time.String())
	} else {
		c.client[client] = models.Client{InClub: true, TableNum: 0}
	}
}

// client takes empty table
func (c *Club) clientTakeTable(time models.Time, client string, tableId int) {
	curClient, ok := c.client[client]
	if !ok || !curClient.InClub {
		fmt.Printf("%s 13 ClientUnknown\n", time.String())
		return
	}

	curTable := c.table[tableId-1]
	if curTable.Occupied {
		fmt.Printf("%s 13 PlaceIsBusy\n", time.String())
		return
	}

	c.client[client] = models.Client{InClub: true, TableNum: tableId}
	c.table[tableId-1] = models.Table{Occupied: true, StartUse: time}
}

// client waits to take table
func (c *Club) clientWait(time models.Time, client string) {
	//check if free tables exist
	for _, t := range c.table {
		if !t.Occupied {
			fmt.Printf("%s 13 ICanWaitNoLonger!\n", time.String())
			return
		}
	}

	// fill client queue
	select {
	case c.queue <- client:
	default:
		//if queue is full
		fmt.Printf("%s 11 %s\n", time.String(), client)
	}
}

// client leaves the club
func (c *Club) clientLeave(time models.Time, client string) {
	curClient, ok := c.client[client]
	if !ok || !curClient.InClub {
		fmt.Printf("%s 13 ClientUnknown\n", time.String())
		return
	}

	tableIdx := curClient.TableNum - 1
	t, ok := c.table[tableIdx]
	if !ok {
		c.client[client] = models.Client{InClub: false, TableNum: 0}
		return
	}
	t.StopUsage(time)

	c.client[client] = models.Client{InClub: false, TableNum: 0}

	// take first client from queue
	select {
	case queueClient := <-c.queue:
		c.client[queueClient] = models.Client{InClub: true, TableNum: tableIdx + 1}
		c.table[tableIdx] = models.Table{Occupied: true, StartUse: time, InUse: t.InUse}
		fmt.Printf("%s 12 %s %d\n", time.String(), queueClient, tableIdx+1)
	default:
		c.table[tableIdx] = t
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
		fmt.Printf("%s 11 %s\n", c.endHour.String(), client)
	}
}

func (c *Club) TableInfo() {

	// Порядок нам не важен судя по заданию
	for i, t := range c.table {
		fmt.Printf("%d %d %s\n", i+1, t.CalculateProfit(c.tariff), t.InUse.String())
	}
}
