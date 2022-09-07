package main

import (
	"fmt"
	"math/rand"
	"time"
)

const eatThreshold = 0.80

func main() {

	allDoneEating := false
	checkChannels := []chan int{make(chan int), make(chan int), make(chan int), make(chan int), make(chan int)}
	forkChannels := []chan string{make(chan string), make(chan string), make(chan string), make(chan string), make(chan string)}

	for i := 0; i < 5; i++ {
		go forkFunc(forkChannels[i])
		if i == 4 {
			go philosopher(fmt.Sprint("phil ", i), checkChannels[i], forkChannels[0], forkChannels[i])
		} else {
			go philosopher(fmt.Sprint("phil ", i), checkChannels[i], forkChannels[i+1], forkChannels[i])
		}
	}

	for !allDoneEating {
		allDoneEating = true
		for i := 0; i < 5; i++ {
			if <-checkChannels[i] < 3 {
				allDoneEating = false
				break
			}
		}
		time.Sleep(1000)
	}

}

func forkFunc(channel chan string) {
	//implement a lock functionality without actual locks
	//implement unlock functionality -||-
	inUse := false
	for true {
		if <-channel == "pick up" && !inUse {

		}
	}
}

func philosopher(name string, checkChan chan int, leftFork chan string, rightFork chan string) {
	var timesEaten int
	var eatChanceBonus float32
	//loop/repeat this indefinately
	flip := rand.Float32() + eatChanceBonus
	if flip < eatThreshold {
		think(name)
		eatChanceBonus += 0.1
	} else {
		//request access to forks
		eat(name)
		timesEaten++
		checkChan <- timesEaten

	}
}

func eat(name string) {
	fmt.Println(name, " is eating, munch munch")
	time.Sleep(100 * time.Millisecond)
}

func think(name string) {
	fmt.Println(name, " is thinking, hmmmm.....")
	time.Sleep(100 * time.Millisecond)
}
