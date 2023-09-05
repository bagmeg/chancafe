package main

import (
	"fmt"
	"log"
	"time"
)

type Cafe struct {
	Name string

	waiter2customerChan chan Product
	customer2waiterChan chan Order

	waiter2baristaChan chan Order
	barista2waiterChan chan Product

	waiters  []IWaiter
	baristas []IBarista

	status bool
}

func NewCafe(name <-chan string, nums <-chan int, res chan<- *Cafe) {
	waiterCount := <-nums
	baristaCount := <-nums

	waiters := make([]IWaiter, 0, waiterCount)
	baristas := make([]IBarista, 0, baristaCount)

	waiter2customerChan := make(chan Product)
	customer2waiterChan := make(chan Order)

	waiter2baristaChan := make(chan Order)
	barista2waiterChan := make(chan Product)

	// make waiters
	waiterChan := make(chan IWaiter, 1)
	nameChan := make(chan string, 1)
	for i := 0; i < waiterCount; i++ {
		go NewWaiter(nameChan, customer2waiterChan, waiter2baristaChan, barista2waiterChan, waiter2customerChan, waiterChan)

		nameChan <- fmt.Sprintf("Waiter%d", i+1)

		waiter := <-waiterChan
		waiters = append(waiters, waiter)
	}

	baristaChan := make(chan IBarista, 1)
	for i := 0; i < baristaCount; i++ {
		go NewBarista(nameChan, waiter2baristaChan, barista2waiterChan, baristaChan)

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
	}
}

func (c *Cafe) Open() {
	// if cafe is opened, do nothing
	if c.status {
		return
	} else {
		c.status = true
	}

	for _, waiter := range c.waiters {
		go waiter.TakeOrder()
	}

	for _, barista := range c.baristas {
		go barista.MakeProduct()
	}

	log.Println("Cafe is opened")
}

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

func (c *Cafe) TakeOrder() {
	for {
		time.Sleep(777 * time.Millisecond)

		// order a random menu
		c.customer2waiterChan <- Order{
			Name: MenuList[time.Now().Nanosecond()%len(MenuList)],
		}

		product := <-c.waiter2customerChan

		log.Printf("Customer got %s, made by %s\n", product.Name, product.Maker)
	}
}
