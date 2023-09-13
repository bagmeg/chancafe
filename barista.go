package main

import (
	"log"
	"sync"
	"time"
)

type IBarista interface {
	MakeProduct()
}

type Barista struct {
	Name string

	orderFromWaiter <-chan Order
	productToWaiter chan<- Product

	stopWorking <-chan struct{}
	status      Status

	order Order

	wg *sync.WaitGroup
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
	var ok bool

	b.wg.Add(1)
	defer b.wg.Done()
	defer log.Printf("Oh!! Barista %s is going home", b.Name)

	if b.status == WORKING {
		return
	}

	b.status = WORKING
	for {
		select {
		case b.order, ok = <-b.orderFromWaiter:
			if !ok {
				// log.Println("Barista: order channel is closed")
				return
			}
			go func() {
				p := Product{Name: b.order.Name, Maker: b.Name}
				time.Sleep(813 * time.Millisecond)
				if b.status == WORKING {
					b.productToWaiter <- p
				}
			}()
		case <-b.stopWorking:
			b.stopBaristWorking()
			return
		}
	}
}

func (b *Barista) stopBaristWorking() {
	b.status = NOT_WORKING

	// log.Println("Barista: order channel is closed")
}
