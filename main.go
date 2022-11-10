package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	dme "github.com/katrinesando/Distributed-Systems/tree/Handin4_DME/grpc"
	"google.golang.org/grpc"
)

var id = flag.Int("id", 1, "id name")
var port = flag.Int("port", 5000, "port name")

func main() {
	flag.Parse()
	ownPort := 5000 + int32(*id)
	if isFlagPassed(int32(*id)) {
		fmt.Printf("id %v is already taken, please use another id", *id)
		//new port here
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &peer{
		id:      ownPort,
		lamport: 0,
		clients: make(map[int32]dme.AccesCriticalClient),
		ctx:     ctx,
		state:   RELEASED,
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
		log.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		c := dme.NewAccesCriticalClient(conn)
		p.clients[port] = c
	}

	rand.Seed(time.Now().UnixNano())
	for {
		flip := rand.Float32()
		if flip > 0.7 {
			p.requestToAll()
		} else {
			p.lamport++
			go p.internalWork()
		}
	}
}

type peer struct {
	dme.UnimplementedAccesCriticalServer
	id      int32
	lamport int32
	clients map[int32]dme.AccesCriticalClient
	ctx     context.Context
	state   State
}

type State int32

const (
	RELEASED State = iota
	WANTED
	HELD
)

func (p *peer) ReplyAccessAttempt(ctx context.Context, req *dme.Request) (*dme.Reply, error) {
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
	request := &dme.Request{Id: p.id, Lamport: p.lamport}
	numValidReplies := 0
	for numValidReplies < len(p.clients) {

		for id, client := range p.clients {
			reply, err := client.ReplyAccessAttempt(p.ctx, request)
			if err != nil {
				log.Println("something went wrong")
			}
			if reply.Answer {
				numValidReplies++
			}
			log.Printf("Got reply from id %v: %v\n", id, reply.Answer)
		}
		if numValidReplies == len(p.clients) {
			p.state = HELD
			p.criticalSection()
			break
		}
		numValidReplies = 0
		time.Sleep(100)
	}
}

func isFlagPassed(port int32) bool {
	found := false

	flag.Visit(func(f *flag.Flag) {
		value, err := strconv.Atoi(f.Name)
		if err != nil {
			log.Println("Error during String to Int")
			return
		}
		if int32(value) == port {
			found = true
		}
	})
	return found
}

func (p *peer) internalWork() {
	log.Printf("%v is performing internal work", p.id)
	time.Sleep(5)
}

func (p *peer) criticalSection() {
	log.Printf("Peer: %v has entered the Critical Section", p.id)
	time.Sleep(5)
	log.Printf("Peer: %v has left the Critical Section", p.id)
	p.lamport++

	p.state = RELEASED
}
