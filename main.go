package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	dme "github.com/katrinesando/Distributed-Systems/tree/Handin4_DME/grpc"
	"google.golang.org/grpc"
)

var id = flag.String("id", "default", "id name")
var port = flag.Int("port", 5000, "port name")

func main() {
	flag.Parse()
	ownPort := int32(*port)
	if isFlagPassed(ownPort) {
		fmt.Printf("Port %v is already taken, please use another port", ownPort)
		//new port here
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &peer{
		id:      ownPort,
		lamport: 0,
		clients: make(map[int32]dme.AccesCriticalClient),
		ctx:     ctx,
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
		port := ownPort

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
	state   State
	//needs some sort of queue - maybe put it on clients
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
	if p.state != RELEASED {
		if p.state == HELD {
			p.lamport++
			return &dme.Reply{
				Answer:  false,
				Lamport: p.lamport,
			}, nil
		} else if p.lamport < otherLamport {
			p.lamport = otherLamport
			p.lamport++
			return &dme.Reply{
				Answer:  false,
				Lamport: p.lamport,
			}, nil
		} else if p.lamport == otherLamport && p.id > otherId {
			p.lamport++
			return &dme.Reply{
				Answer:  false,
				Lamport: p.lamport,
			}, nil
		}

	}
	p.lamport++
	return &dme.Reply{
		Answer:  true,
		Lamport: p.lamport,
	}, nil
}

func (p *peer) requestToAll() {
	p.state = WANTED
	p.lamport++
	log.Printf("%v is requesting access to critial section", p.id)
	request := &dme.Request{Id: p.id, Lamport: p.lamport} //needs to send lamport time stamp to all to others
	for id, client := range p.clients {
		reply, err := client.AttemptAcces(p.ctx, request)
		if err != nil {
			fmt.Println("something went wrong")
		}
		fmt.Printf("Got reply from id %v: %v\n", id, reply.Answer)
	}
}

func isFlagPassed(port int32) bool {
	found := false

	flag.Visit(func(f *flag.Flag) {
		value, err := strconv.Atoi(f.Name)
		if err != nil {
			fmt.Println("Error during String to Int")
			return
		}
		if int32(value) == port {
			found = true

		}
	})
	return found
}

//function to see who has priority

func CriticalSection(p peer) {
	log.Printf("Peer: %v has entered the Critical Section", p.id)
	time.Sleep(5)
	log.Printf("Peer: %v has left the Critical Section", p.id)
}
