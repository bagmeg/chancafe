package main

type Cafe struct {
	Name string
}

func NewCafe(n <-chan string, res chan<- *Cafe) {
	name := <-n

	c := &Cafe{
		Name: name,
	}

	res <- c
}
