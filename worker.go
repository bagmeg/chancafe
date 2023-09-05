package main

import (
	"time"
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
}

type Waiter struct {
	Name string

	orderFromCustomer <-chan Order
	orderToBarista    chan<- Order

	productFromBarista <-chan Product
	productToCustomer  chan<- Product
}

func NewWaiter(
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
	}

	res <- w
}

func NewBarista(
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
	}

	res <- b
}

func (b *Barista) MakeProduct() {
	for {
		order := <-b.orderFromWaiter

		p := Product{Name: order.Name, Maker: b.Name}

		// takes some time to make product
		time.Sleep(813 * time.Millisecond)

		b.productToWaiter <- p
	}
}

func (w *Waiter) TakeOrder() {
	for {
		order := <-w.orderFromCustomer

		// takes some time to get order
		time.Sleep(187 * time.Millisecond)

		w.orderToBarista <- order

		product := <-w.productFromBarista

		w.productToCustomer <- product
	}
}
