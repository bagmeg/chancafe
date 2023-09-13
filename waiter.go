package main

import (
	"log"
	"sync"
)

type IWaiter interface {
	TakeOrder()
}

type Waiter struct {
	Name string

	orderFromCustomer <-chan Order
	orderToBarista    chan<- Order

	productFromBarista <-chan Product
	productToCustomer  chan<- Product

	stopWorking <-chan struct{}
	status      Status

	order   Order
	product Product

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

func (w *Waiter) TakeOrder() {
	var ok bool

	w.wg.Add(1)
	defer w.wg.Done()
	defer log.Printf("Oh!! Waiter %s is going home", w.Name)

	if w.status == WORKING {
		return
	}

	w.status = WORKING
	for {
		select {
		case w.order, ok = <-w.orderFromCustomer:
			if !ok {
				// log.Printf("Waiter %s: order channel is closed", w.Name)
				return
			}
			w.orderToBarista <- w.order
		case w.product, ok = <-w.productFromBarista:
			if !ok {
				// log.Printf("Waiter %s: product channel is closed", w.Name)
				return
			}
			w.productToCustomer <- w.product
		case <-w.stopWorking:
			w.stopWaiterWorking()
			return
		}
	}
}

func (w *Waiter) stopWaiterWorking() {
	w.status = NOT_WORKING

	// log.Printf("Waiter %s stopping", w.Name)
}
