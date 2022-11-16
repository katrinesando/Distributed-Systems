package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"grpc"

	"google.golang.org/grpc"
	"handin5.dk/uni/grpc"
)

var id int32

func main() {
	LOG_FILE := "./txtLog"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithInsecure())

	client := grpc.NewAuctionClient()

	for i := 0; i < 3; i++ {
		port := int32(5000) + int32(i)

		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), opts...)
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}

		client.Connect(conn)
		defer conn.Close()
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		test := scanner.Text()
		if strings.Contains(test, "Bid") {
			bid, _ := strconv.Atoi(scanner.Text())
			go sendBid(ctx, client, bid)
		}
		if strings.Contains(test, "Result") {
			go getResult(ctx, client)
		}

	}
}

func sendBid(ctx context.Context, client grpc.AuctionClient, bidAmount int32) {
	stream, err := client.SendBid(ctx)
	if err != nil {
		log.Printf("Cannot send bid: error: %v", err)
	}

	msg := grpc.Bid{
		BidAmount: bidAmount,
		Id:        id,
	}
	stream.Send(&msg)

	ack, err := stream.CloseAndRecv()
	if err != nil {
		log.Print("Cannot send ack: %v", err)
		log.Print(ack)
	}

}

func getResult(ctx context.Context, client grpc.AuctionClient) {

	stream, err := client.GetResults(ctx)
	if err != nil {
		log.Printf("Cannot send bid: error: %v", err)
	}

	stream.Send()

	result, err := stream.CloseAndRecv()
	if err != nil {
		log.Print("Cannot send ack: %v", err)
		log.Print(result)
	}
	log.Printf(result)

}
