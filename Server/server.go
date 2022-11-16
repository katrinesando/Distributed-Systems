package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	handin "handin5.dk/uni/grpc"
)

var port = flag.String("server", ":9100", "Tcp server")

func (s *server) SendBid(ctx context.Context, b *handin.Bid) (*handin.Ack, error) {
	for _, test := range s.AuctionBids {
		if proto.Equal(test.BidAmount, b) {
			ack := handin.Ack{Outcome: 1}
			return &ack, nil
		}
	}
	// No ack was found, return an unnamed ack with 0
	return &handin.Ack{Outcome: 0}, nil
}

func (s *server) GetResults(p handin.AuctionServer) error {
	return proto.ErrNil
}

func main() {
	flag.Parse()
	//log to different txt Log file
	LOG_FILE := "./txtLog"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	lis, err := net.Listen("tcp", "localhost%v", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	handin.RegisterAuctionServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}

func newServer() *server {
	s := &server{
		AuctionBids: make(map[int32]*handin.Bid),
	}
	return s
}

type server struct {
	handin.UnimplementedAuctionServer
	AuctionBids map[int32]*handin.Bid
}
