package main

import "sync"

type Manager struct {
	// customer -> waiter
	OrderChan chan Order
	// waiter -> barista
	OrderToBaristachan chan Order
	// barista -> waiter
	ProductToWaiterChan chan Product
	// waiter -> customer
	ProductChan chan Product

	stopWorkingChan chan struct{}
	wg              *sync.WaitGroup
}

func NewManager() *Manager {
	return &Manager{
		OrderChan:           make(chan Order),
		OrderToBaristachan:  make(chan Order),
		ProductToWaiterChan: make(chan Product),
		ProductChan:         make(chan Product),
		stopWorkingChan:     make(chan struct{}),
		wg:                  &sync.WaitGroup{},
	}
}

func (m *Manager) Start() {

}

func (m *Manager) Close() {
	close(m.OrderChan)
	close(m.OrderToBaristachan)
	close(m.ProductToWaiterChan)
	close(m.ProductChan)
	close(m.stopWorkingChan)
}

func (m *Manager) GetOrderChan() chan Order {
	return m.OrderChan
}

func (m *Manager) GetOrderToBaristachan() chan Order {
	return m.OrderToBaristachan
}

func (m *Manager) GetProductToWaiterChan() chan Product {
	return m.ProductToWaiterChan
}

func (m *Manager) GetProductChan() chan Product {
	return m.ProductChan
}

func (m *Manager) GetStopWorkingChan() chan struct{} {
	return m.stopWorkingChan
}
