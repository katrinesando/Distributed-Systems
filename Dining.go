/*Explanation for why this system should avoid deadlocks
Firstly having a channel for each direction of communication per goroutine avoids
sends and recieves blocking each other.
Secondly reading from channels into variables lets us avoid situations where a single channel blocks the rest.
Lastly the philosopher method contains a randomized float which the philosphers deciesions are based upon.
This lets philosphers "choose" what to do somewhat randomly every iteration of the loop

*/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const eatThreshold = 0.60

func main() {

	fork1 := []chan string{make(chan string), make(chan string)}
	fork2 := []chan string{make(chan string), make(chan string)}
	fork3 := []chan string{make(chan string), make(chan string)}
	fork4 := []chan string{make(chan string), make(chan string)}
	fork5 := []chan string{make(chan string), make(chan string)}

	go forkFunc(fork1)
	go forkFunc(fork2)
	go forkFunc(fork3)
	go forkFunc(fork4)
	go forkFunc(fork5)

	go philosopher("Socrates", fork1, fork2)
	go philosopher("Nietzche", fork2, fork3)
	go philosopher("Sartre", fork3, fork4)
	go philosopher("Plato", fork4, fork5)
	go philosopher("de Beauvoir", fork5, fork1)

	for {

	}

}

func forkFunc(channels []chan string) {
	inUse := false
	for true {
		mes := <-channels[0]
		if mes == "pick up" && !inUse {
			inUse = true
			channels[1] <- "use"
		} else if mes == "pick up" && inUse {
			channels[1] <- "no"
		} else if mes == "finished" {
			inUse = false
		}
	}
}

func philosopher(name string, leftFork []chan string, rightFork []chan string) {
	var timesEaten int
	for true {
		flip := rand.Float32()
		if flip < eatThreshold {

			think(name)
		} else {
			fmt.Println(name, " attempting to eat")
			leftFork[0] <- "pick up"
			rightFork[0] <- "pick up"
			time.Sleep(10)
			gotLeftFork := <-leftFork[1] == "use"
			gotRightFork := <-rightFork[1] == "use"
			if gotLeftFork && gotRightFork {

				timesEaten++
				leftFork[0] <- "finished"
				rightFork[0] <- "finished"
				eat(name)
			} else {
				fmt.Println(name, " got no fork so sad")
				leftFork[0] <- "finished"
				rightFork[0] <- "finished"
			}

		}
	}
}

func eat(name string) {
	fmt.Println(name, " is eating, munch munch")

}

func think(name string) {
	fmt.Println(name, " is thinking, hmmmm.....")

}
