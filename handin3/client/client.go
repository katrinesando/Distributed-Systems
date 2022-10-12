package main

import (
	"log"

	"google.golang.org/grpc"
)

type Vector struct {
	clock []int
}

var id int
var clock Vector

func main() {
	// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()

}

func UpdateClock(recievedClock Vector) {
	sameLen := len(clock.clock) == len(recievedClock.clock)
	if sameLen {
		for i := 0; i < len(recievedClock.clock); i++ {

		}
	}
}
