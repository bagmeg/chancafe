package main

import (
	"fmt"
	"time"
)

type Product struct {
	Name  string
	Maker string
}

type Order struct {
	Name string
}

type IBarista interface {
	MakeProduct()
}

type IWaiter interface {
	TakeOrder(<-chan Order)
}

type Barista struct {
	Name    string
	Worder  <-chan Order   // receive order from waiter
	Wresult chan<- Product // send product to waiter
}

type Waiter struct {
	Name    string
	Worder  chan<- Order   // send order to barista
	Wresult <-chan Product // receive product from barista
}

func NewWaiter(n <-chan string, or chan<- Order, re <-chan Product, res chan<- IWaiter) {
	name := <-n

	w := &Waiter{
		Name:    name,
		Worder:  or,
		Wresult: re,
	}

	res <- w
}

func NewBarista(n <-chan string, or <-chan Order, re chan<- Product, res chan<- IBarista) {
	name := <-n

	b := &Barista{
		Name:    name,
		Worder:  or,
		Wresult: re,
	}

	res <- b
}

func (b *Barista) MakeProduct() {
	for {
		o := <-b.Worder

		p := Product{Name: o.Name, Maker: b.Name}

		time.Sleep(333 * time.Millisecond)
		b.Wresult <- p
	}
}

func (w *Waiter) TakeOrder(o <-chan Order) {
	for {
		order := <-o

		w.Worder <- order

		p := <-w.Wresult

		fmt.Println("Waiter", w.Name, "serves", p.Name, "made by", p.Maker)
	}
}
