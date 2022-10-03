package main

import (
	"fmt"
	"time"

	"github.com/eiannone/keyboard"
	"golang.org/x/sync/semaphore"
)

type Customer struct {
	ID int
}

type Barbershop struct {
	ChairCount int
	Chairs     chan Customer
	Semaphore  semaphore.Weighted
	Open       bool
}

// OpenShop initializes the barbershop
func (b *Barbershop) OpenShop() {
	fmt.Println("Opening barbershop")
	b.Chairs = make(chan Customer, b.ChairCount)
	b.Semaphore = *semaphore.NewWeighted(int64(b.ChairCount))
	b.Open = true
}

// CloseShop closes the channel and state of the barbershop
func (b *Barbershop) CloseShop() {
	if !b.Open {
		panic("Barbershop is already closed")
	}
	fmt.Println("Closing barbershop")
	close(b.Chairs)
	b.Open = false
}

// WalkIn adds a customer to the barbershop if there is an available chair
func (b *Barbershop) WalkIn(customer Customer) {
	if !b.Open {
		panic("Barbershop is closed")
	}
	fmt.Printf("%d Customer walks in\n", customer.ID)
	if b.Semaphore.TryAcquire(1) {
		fmt.Printf("%d Customer waiting for barber\n", customer.ID)
		b.Chairs <- customer
		fmt.Println(customer.ID, "Customer sits in chair")
	} else {
		fmt.Printf("%d No chair avaliable for customer; Leaving\n", customer.ID)
	}
}

// ServeCustomerOrSleep checks if a customer is in a chair and cuts their hair, or sleeps if no customer is waiting
func (b *Barbershop) ServeCustomerOrSleep() {
	if !b.Open {
		panic("Barbershop is closed")
	}
	noCustomers := false
	for b.Open {
		if len(b.Chairs) > 0 {
			noCustomers = false
			customer := <-b.Chairs
			fmt.Printf("%d Barber cuts hair\n", customer.ID)
			time.Sleep(1 * time.Second)
			fmt.Printf("%d Customer leaving\n", customer.ID)
			b.Semaphore.Release(1)
		} else {
			if !noCustomers {
				fmt.Println("Empty chairs; Barber sleeping")
			}
			noCustomers = true
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	barbershop := Barbershop{ChairCount: 3}
	barbershop.OpenShop()
	go barbershop.ServeCustomerOrSleep()

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press space bar to add a customer, or press ESC to exit")
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		// spacebar
		if key == keyboard.KeySpace {
			id := time.Now().UnixNano()
			go barbershop.WalkIn(Customer{ID: int(id)})
		}

		if key == keyboard.KeyEsc {
			barbershop.CloseShop()
			break
		}
	}
}
