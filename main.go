package main

func main() {
	nameChan := make(chan string, 1)
	numChan := make(chan int, 2)
	cafeChan := make(chan *Cafe, 1)

	go NewCafe(nameChan, numChan, cafeChan)

	numChan <- 2
	numChan <- 3

	nameChan <- "Starbucks"

	cafe := <-cafeChan

	cafe.Open()

	// start taking order from customers
	cafe.TakeOrder()
}
