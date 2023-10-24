package main

import (
	"fmt"
	"sync"
	"time"
)

type Tunnel struct {
	tunnelName string
	trainQueue chan int
}

func NewTunnel(tunnelName string) *Tunnel {
	return &Tunnel{
		tunnelName: tunnelName,
		trainQueue: make(chan int, 10),
	}
}

func (t *Tunnel) enterTunnel(trainNumber int) {
	t.trainQueue <- trainNumber
	fmt.Printf("Потяг %d очікує на в'їзд у %s\n", trainNumber, t.tunnelName)

	oppositeTunnelName := "Тунель 2"
	if t.tunnelName == "Тунель 2" {
		oppositeTunnelName = "Тунель 1"
	}

	oppositeTunnel := getTunnel(oppositeTunnelName)

	if len(oppositeTunnel.trainQueue) > 0 {
		oppositeTrain := <-oppositeTunnel.trainQueue
		fmt.Printf("Потяг %d перевіряє чи є Потяг %d в протилежному тунелі.\n", trainNumber, oppositeTrain)

		if hasExceededWaitTime(oppositeTrain) {
			fmt.Printf("Потяг %d змінює рух до %s\n", trainNumber, t.tunnelName)
			t.trainQueue <- oppositeTrain
		} else {
			oppositeTunnel.trainQueue <- oppositeTrain
		}
	}

	fmt.Printf("Потяг %d проїжджає %s тунель.\n", trainNumber, t.tunnelName)
	time.Sleep(10 * time.Second)
	fmt.Printf("Потяг %d проїхав %s тунель.\n", trainNumber, t.tunnelName)
	<-t.trainQueue
	incrementCounter(t.tunnelName)
}

type Train struct {
	trainNumber int
	tunnel      *Tunnel
}

func NewTrain(trainNumber int, tunnel *Tunnel) *Train {
	return &Train{
		trainNumber: trainNumber,
		tunnel:      tunnel,
	}
}

func (t *Train) run() {
	t.tunnel.enterTunnel(t.trainNumber)
}

var tunnel1 = NewTunnel("Тунель 1")
var tunnel2 = NewTunnel("Тунель 2")
var trainsPassedTunnel1 = 0
var trainsPassedTunnel2 = 0
var trainArrivalTimes = make(map[int]time.Time)
var mutex = &sync.Mutex{}
var wg sync.WaitGroup

func getTunnel(tunnelName string) *Tunnel {
	if tunnelName == "Тунель 1" {
		return tunnel1
	} else {
		return tunnel2
	}
}

func incrementCounter(tunnelName string) {
	mutex.Lock()
	defer mutex.Unlock()

	if tunnelName == "Тунель 1" {
		trainsPassedTunnel1++
	} else if tunnelName == "Тунель 2" {
		trainsPassedTunnel2++
	}

	totalTrains := trainsPassedTunnel1 + trainsPassedTunnel2
	if totalTrains == 20 {
		fmt.Println("Колії вільні")
	}
}

func hasExceededWaitTime(trainNumber int) bool {
	mutex.Lock()
	defer mutex.Unlock()

	arrivalTime, exists := trainArrivalTimes[trainNumber]
	if !exists {
		trainArrivalTimes[trainNumber] = time.Now()
		return false
	}

	currentTime := time.Now()
	waitingTime := currentTime.Sub(arrivalTime)
	return waitingTime > 60*time.Second
}

func main() {
	var totalTrains int
	fmt.Println("Яку кількість потягів очікуємо?")
	fmt.Scanln(&totalTrains)

	for i := 1; i <= totalTrains; i++ {
		selectedTunnel := tunnel1
		if i%2 == 0 {
			selectedTunnel = tunnel2
		}
		wg.Add(1)
		go func(trainNumber int, tunnel *Tunnel) {
			train := NewTrain(trainNumber, tunnel)
			train.run()
			wg.Done()
		}(i, selectedTunnel)
	}

	wg.Wait()
}
