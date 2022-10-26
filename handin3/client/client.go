package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	h "handin3"
	"io"
	"log"
	"os"

	chatpb "handin3/chatpb"

	"google.golang.org/grpc"
)

var channelName = flag.String("channel", "default", "Channel name for chatting")
var senderName = flag.String("sender", "default", "Senders name")
var tcpServer = flag.String("server", ":9100", "Tcp server")
var id int32

func joinChannel(ctx context.Context, client chatpb.ChatServiceClient) {

	IncrementClock()

	channel := chatpb.Channel{Name: *channelName, SendersName: *senderName, Clock: clock.Clock}

	stream, err := client.JoinChannel(ctx, &channel)
	if err != nil {
		log.Fatalf("couldnt connect to client: %v", err)
	}

	waitc := make(chan struct{}) //wait for something to finish, requires almost no memory

	go func() {
		for {
			in, err := stream.Recv()
			if err != nil {
				log.Fatalf("Failed to receive message from channel joining. \nErr: %v", err)
			}
			if err == io.EOF { //done sending message
				close(waitc)
				return
			}

			mes := fmt.Sprintf("%v", in.Message)
			if mes == "4040" {
				fmt.Printf("--- %v Left the chat ---\n", in.Sender)
			} else if mes == "1111" {
				fmt.Printf("--- %v Joined the Chat ---\n", in.Sender)
				id = in.Id
			} else if *senderName != in.Sender {
				fmt.Printf("(%v --ID: %v) : %v \n", in.Sender, in.Id, in.Message)
			}
			recClock := h.Vector{
				Clock: in.Channel.Clock,
			}

			UpdateClock(recClock)
		}
	}()
	<-waitc
}

func sendMessage(ctx context.Context, client chatpb.ChatServiceClient, message string) {
	IncrementClock()
	stream, err := client.SendMessage(ctx)
	if err != nil {
		log.Printf("Cannot send message: error: %v", err)
	}
	msg := chatpb.Message{
		Channel: &chatpb.Channel{
			Name:        *channelName,
			SendersName: *senderName,
			Clock:       clock.Clock},
		Message: message,
		Sender:  *senderName,

		Id: id,
	}
	stream.Send(&msg)

	//sends message to console
	ack, err := stream.CloseAndRecv()
	if err != nil {
		log.Print("Cannot send ack: %v", err)
		log.Print(ack)
	}

}

var clock h.Vector

func main() {

	clock.Clock = make([]int32, 1, 1)

	flag.Parse()

	fmt.Println("┌─────────────────────────────┐")
	fmt.Println("│            		      │")
	fmt.Println("│ Welcome to the Chat Service │")
	fmt.Println("│            		      │")
	fmt.Println("└─────────────────────────────┘")
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithInsecure())

	conn, err := grpc.Dial(*tcpServer, opts...)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	defer conn.Close()

	ctx := context.Background()
	client := chatpb.NewChatServiceClient(conn)

	go joinChannel(ctx, client)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		go sendMessage(ctx, client, scanner.Text())

	}

}

func UpdateClock(recievedClock h.Vector) {
	clock = h.AdjustToOtherClock(clock, recievedClock)
	IncrementClock()
}

func IncrementClock() {
	correctLen := id >= int32(len(clock.Clock))
	for correctLen {
		clock.Clock = append(clock.Clock, 0)
		correctLen = id >= int32(len(clock.Clock))
	}
	clock.Clock[id]++
}
