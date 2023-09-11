package main

import (
	"log"
	"sync"
	"time"
)

type Status int

const (
	WORKING = iota
	NOT_WORKING
)

type Product struct {
	Name  Menu
	Maker string
}

type Order struct {
	Name Menu
}

type IBarista interface {
	MakeProduct()
}

type IWaiter interface {
	TakeOrder()
}

type Barista struct {
	Name string

	orderFromWaiter <-chan Order
	productToWaiter chan<- Product

	stopWorking <-chan struct{}
	status      Status

	wg *sync.WaitGroup
}

type Waiter struct {
	Name string

	orderFromCustomer <-chan Order
	orderToBarista    chan<- Order

	productFromBarista <-chan Product
	productToCustomer  chan<- Product

	stopWorking <-chan struct{}
	status      Status

	wg *sync.WaitGroup
}

func NewWaiter(
	stopWorking <-chan struct{},
	wg *sync.WaitGroup,
	n <-chan string,
	orderFromCustomer <-chan Order,
	orderToBarista chan<- Order,
	productFromBarista <-chan Product,
	productToCustomer chan<- Product,
	res chan<- IWaiter,
) {
	name := <-n

	w := &Waiter{
		Name:               name,
		orderFromCustomer:  orderFromCustomer,
		orderToBarista:     orderToBarista,
		productFromBarista: productFromBarista,
		productToCustomer:  productToCustomer,
		stopWorking:        stopWorking,
		wg:                 wg,
	}

	res <- w
}

func NewBarista(
	stopWorking <-chan struct{},
	wg *sync.WaitGroup,
	n <-chan string,
	orderFromWaiter <-chan Order,
	productToWaiter chan<- Product,
	res chan<- IBarista,
) {
	name := <-n

	b := &Barista{
		Name:            name,
		orderFromWaiter: orderFromWaiter,
		productToWaiter: productToWaiter,
		stopWorking:     stopWorking,
		wg:              wg,
	}

	res <- b
}

func (b *Barista) MakeProduct() {
	b.wg.Add(1)
	defer b.wg.Done()
	defer log.Printf("Oh!! Barista %s is going home", b.Name)

	go b.stopBaristWorking()

	b.status = WORKING
	for {
		if b.status == NOT_WORKING {
			return
		}

		order, ok := <-b.orderFromWaiter
		if !ok {
			log.Println("Barista: order channel is closed")
			return
		}

		p := Product{Name: order.Name, Maker: b.Name}

		// takes some time to make product
		time.Sleep(813 * time.Millisecond)

		b.productToWaiter <- p
	}
}

func (b *Barista) stopBaristWorking() {
	<-b.stopWorking
	b.status = NOT_WORKING
}

func (w *Waiter) TakeOrder() {
	w.wg.Add(1)
	defer w.wg.Done()
	defer log.Printf("Oh!! Waiter %s is going home", w.Name)

	go w.stopWaiterWorking()

	w.status = WORKING
	for {
		if w.status == NOT_WORKING {
			return
		}

		order, ok := <-w.orderFromCustomer
		if !ok {
			log.Println("Waiter: order channel is closed")
			return
		}

		// takes some time to get order
		time.Sleep(187 * time.Millisecond)

		w.orderToBarista <- order

		product := <-w.productFromBarista

		w.productToCustomer <- product
	}
}

func (w *Waiter) stopWaiterWorking() {
	<-w.stopWorking
	w.status = NOT_WORKING
}
