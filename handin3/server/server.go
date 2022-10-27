package main

import (
	"fmt"
	h "handin3"
	"io"
	"log"
	"net"
	"os"
	"sync"

	chatpb "handin3/chatpb"

	"google.golang.org/grpc"
)

var clock h.Vector
var id int32

type server struct {
	chatpb.UnimplementedChatServiceServer
	mu      sync.Mutex
	channel map[string][]chan *chatpb.Message
}

func (s *server) JoinChannel(ch *chatpb.Channel, msgStream chatpb.ChatService_JoinChannelServer) error {
	msgChannel := make(chan *chatpb.Message)
	s.channel[ch.Name] = append(s.channel[ch.Name], msgChannel)
	log.Printf("SERVER: Participant %v joined Chitty-Chat at Lamport time %v\n", ch.SendersName, ch.Clock)
	fmt.Printf("SERVER: Participant %v joined Chitty-Chat at Lamport time %v\n", ch.SendersName, ch.Clock)
	id++
	broadcastMsg := chatpb.Message{
		Channel: &chatpb.Channel{
			Name:        ch.Name,
			SendersName: ch.SendersName,
			Clock:       clock.Clock},
		Message: "1111",
		Sender:  ch.SendersName,
		Id:      id,
	}
	go func() {
		streams := s.channel[broadcastMsg.Channel.Name]
		for _, msgChan := range streams {
			msgChan <- &broadcastMsg
		}
	}()

	// doing this never closes the stream
	for {
		select {
		case <-msgStream.Context().Done():
			log.Printf("SERVER: Participant %v left Chitty-Chat at Lamport time %v\n", ch.SendersName, ch.Clock)
			fmt.Printf("SERVER: Participant %v left Chitty-Chat at Lamport time %v\n", ch.SendersName, ch.Clock)
			for i, element := range s.channel[ch.Name] {
				if element == msgChannel {
					s.channel[ch.Name] = append(s.channel[ch.Name][:i], s.channel[ch.Name][i+1:]...)
					break
				}
			}

			clock.Clock[0]++
			msg := chatpb.Message{
				Channel: &chatpb.Channel{
					Name:        ch.Name,
					SendersName: ch.SendersName,
					Clock:       clock.Clock},
				Message: "4040",
				Sender:  ch.SendersName,
			}

			go func() {
				streams := s.channel[msg.Channel.Name]
				for _, msgChan := range streams {
					msgChan <- &msg
				}
			}()

			return nil
		case msg := <-msgChannel:
			log.Printf("SERVER: Sent %v from: %v at Lamport time %v", msg.Message, msg.Sender, msg.Channel.Clock)
			fmt.Printf("SERVER: Sent %v from: %v at Lamport time %v", msg.Message, msg.Sender, msg.Channel.Clock)
			recClock := h.Vector{
				Clock: msg.Channel.Clock,
			}

			UpdateClock(recClock)

			go msgStream.Send(msg)

		}

	}
}

func (s *server) SendMessage(msgStream chatpb.ChatService_SendMessageServer) error {
	msg, err := msgStream.Recv() //receive message

	if err == io.EOF {
		return nil
	}

	if err != nil {
		return err
	}

	recClock := h.Vector{
		Clock: msg.Channel.Clock,
	}

	UpdateClock(recClock)

	ack := chatpb.MessageAck{Status: "SENT"}
	msgStream.SendAndClose(&ack) //sends back it is acknowledged - only used for log right now

	go func() {
		streams := s.channel[msg.Channel.Name]
		for _, msgChan := range streams {
			msgChan <- msg
		}
	}()

	return nil
}

func newServer() *server {
	s := &server{
		channel: make(map[string][]chan *chatpb.Message),
	}
	return s
}

func main() {
	//log to different txt Log file
	LOG_FILE := "./txtLog"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	id = 0
	clock.Clock = make([]int32, 1, 1)
	fmt.Println("--- SERVER APP ---")
	lis, err := net.Listen("tcp", "localhost:9100")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	chatpb.RegisterChatServiceServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}

func UpdateClock(recievedClock h.Vector) {
	clock = h.AdjustToOtherClock(clock, recievedClock)
	clock.Clock[0]++
}
