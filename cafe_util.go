package main

type Status int

const (
	NOT_WORKING = iota
	WORKING
)

type Product struct {
	Name  Menu
	Maker string
}

type Order struct {
	Name Menu
}
