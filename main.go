package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	dme "github.com/katrinesando/Distributed-Systems/tree/Handin4_DME/grpc"
	"google.golang.org/grpc"
)

func main() {
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int32(arg1) + 5000

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &peer{
		id:            ownPort,
		amountOfPings: make(map[int32]int32),
		clients:       make(map[int32]dme.AccesCriticalClient),
		ctx:           ctx,
	}

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	grpcServer := grpc.NewServer()
	dme.RegisterAccesCriticalServer(grpcServer, p)

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()

	for i := 0; i < 3; i++ {
		port := int32(5000) + int32(i)

		if port == ownPort {
			continue
		}

		var conn *grpc.ClientConn
		fmt.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		c := dme.NewAccesCriticalClient(conn)
		p.clients[port] = c
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

	}
}

type peer struct {
	dme.UnimplementedAccesCriticalServer
	id      int32
	lamport int32
	clients map[int32]dme.AccesCriticalClient
	ctx     context.Context
}

type State int32

const (
	RELEASED State = iota
	WANTED
	HELD
)

func (p *peer) AttemptAcces(ctx context.Context, req *dme.Request) (*dme.Reply, error) {
	otherId := req.Id
	otherLamport := req.Lamport
}

func CriticalSection(p peer) {
	log.Printf("Peer: ", p.id, "has entered the Critical Section at Lamport", p.lamport)
	time.Sleep(5)
	p.lamport++
	log.Printf("Peer: ", p.id, "has left the Critical Section at Lamport", p.lamport)
	p.state = RELEASED
}
