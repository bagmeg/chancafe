package main

import "time"

func main() {
	orderChan := make(chan Order)
	productChan := make(chan Product)

	baristaChan := make(chan IBarista)
	waiterChan := make(chan IWaiter)

	baristas := make([]IBarista, 0)
	waiters := make([]IWaiter, 0)

	nameChan := make(chan string)

	// make barista James
	go NewBarista(nameChan, orderChan, productChan, baristaChan)
	nameChan <- "James"

	baristas = append(baristas, <-baristaChan)

	// make barista John
	go NewBarista(nameChan, orderChan, productChan, baristaChan)
	nameChan <- "John"

	baristas = append(baristas, <-baristaChan)

	// make waiter Brand
	go NewWaiter(nameChan, orderChan, productChan, waiterChan)
	nameChan <- "Brand"

	waiters = append(waiters, <-waiterChan)

	// make waiter Summer
	go NewWaiter(nameChan, orderChan, productChan, waiterChan)
	nameChan <- "Summer"

	waiters = append(waiters, <-waiterChan)

	counterChan := make(chan Order)
	for _, w := range waiters {
		go w.TakeOrder(counterChan)
	}

	for _, b := range baristas {
		go b.MakeProduct()
	}

	for {
		order := Order{
			Name: "Cappuccino",
		}

		counterChan <- order
		time.Sleep(77 * time.Millisecond)
	}
}
