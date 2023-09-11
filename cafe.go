package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Menu string

const (
	Americano            Menu = "Americano"
	Latte                Menu = "Latte"
	Cappuccino           Menu = "Cappuccino"
	IcedAmericano        Menu = "IcedAmericano"
	StrawberryLatte      Menu = "StrawberryLatte"
	SpagehttiWithKetchup Menu = "SpagehttiWithKetchup"
	PineapplePizza       Menu = "PineapplePizza"
	GimchiPata           Menu = "GimchiPata"
)

var (
	MenuList = []Menu{
		Americano,
		Latte,
		Cappuccino,
		IcedAmericano,
		StrawberryLatte,
	}
)

type Cafe struct {
	Name string

	waiter2customerChan chan Product
	customer2waiterChan chan Order

	waiter2baristaChan chan Order
	barista2waiterChan chan Product

	stopWorkingChan chan struct{}

	wg *sync.WaitGroup

	waiters  []IWaiter
	baristas []IBarista

	totalWorkers int

	opened bool
}

func NewCafe(name <-chan string, nums <-chan int, res chan<- *Cafe) {
	waiterCount := <-nums
	baristaCount := <-nums

	totalWorkers := waiterCount + baristaCount

	waiters := make([]IWaiter, 0, waiterCount)
	baristas := make([]IBarista, 0, baristaCount)

	waiter2customerChan := make(chan Product, 1)
	customer2waiterChan := make(chan Order, 1)

	waiter2baristaChan := make(chan Order, 3)
	barista2waiterChan := make(chan Product, 3)

	stopWorkingChan := make(chan struct{}, totalWorkers)
	wg := sync.WaitGroup{}

	// make waiters
	waiterChan := make(chan IWaiter, 1)
	nameChan := make(chan string, 1)
	for i := 0; i < waiterCount; i++ {
		go NewWaiter(stopWorkingChan, &wg, nameChan, customer2waiterChan, waiter2baristaChan, barista2waiterChan, waiter2customerChan, waiterChan)

		nameChan <- fmt.Sprintf("Waiter%d", i+1)

		waiter := <-waiterChan
		waiters = append(waiters, waiter)
	}

	baristaChan := make(chan IBarista, 1)
	for i := 0; i < baristaCount; i++ {
		go NewBarista(stopWorkingChan, &wg, nameChan, waiter2baristaChan, barista2waiterChan, baristaChan)

		nameChan <- fmt.Sprintf("Barista%d", i+1)

		barista := <-baristaChan
		baristas = append(baristas, barista)
	}

	res <- &Cafe{
		Name: <-name,

		waiters:  waiters,
		baristas: baristas,

		waiter2customerChan: waiter2customerChan,
		customer2waiterChan: customer2waiterChan,

		waiter2baristaChan: waiter2baristaChan,
		barista2waiterChan: barista2waiterChan,

		stopWorkingChan: stopWorkingChan,
		wg:              &wg,

		totalWorkers: totalWorkers,
	}
}

func (c *Cafe) Open() {
	// if cafe is opened, do nothing
	if c.opened {
		return
	} else {
		c.opened = true
	}

	for _, waiter := range c.waiters {
		go waiter.TakeOrder()
	}

	for _, barista := range c.baristas {
		go barista.MakeProduct()
	}

	log.Printf("Cafe %s is opened", c.Name)
}

func (c *Cafe) Close() {
	// if cafe is closed, do nothing
	if !c.opened {
		return
	}

	c.opened = false

	for i := 0; i < c.totalWorkers; i++ {
		c.stopWorkingChan <- struct{}{}
	}

	close(c.barista2waiterChan)
	close(c.waiter2baristaChan)
	close(c.waiter2customerChan)
	close(c.customer2waiterChan)

	c.wg.Wait()

	log.Printf("Cafe %s is closed", c.Name)
}

func (c *Cafe) TakeOrder() {
	for {
		if !c.opened {
			break
		}

		time.Sleep(777 * time.Millisecond)

		// order a random menu
		c.customer2waiterChan <- Order{
			Name: MenuList[time.Now().Nanosecond()%len(MenuList)],
		}

		product, ok := <-c.waiter2customerChan
		if !ok {
			log.Println("Cafe: waiter2customerChan is closed")
			break
		}

		log.Printf("Customer got %s, made by %s\n", product.Name, product.Maker)
	}
}
