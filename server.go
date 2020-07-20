package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"time"
	"log"
)

// Server has three possible states:
//   0. Follower
//   1. Candidate
//   2. Leader
type Server struct {
	ID              int
	Term            int    // Life span of the current leader
	Log             []Log  // List of all logs
	Port            string // tcp Port
	State           int    // Current mode of the server
	Servers         []*Server
	AliveServers    []bool
	Hb              chan int
	VoteRequested   chan bool
	VotedFor           int
	TotalVotes      []bool
	NumAliveServers int
}

// Server Struct toString
func (s Server) String() string {
    return fmt.Sprintf(
    	"{ID:%d,  Term: %d, Log:%v, Port:%s \n State:%d,VotedFor: %d, NumAliveServers: %d,TotalVotes: %v }",
    	s.ID,s.Term,s.Log, s.Port, s.State, s.VotedFor, s.NumAliveServers, s.TotalVotes, 
    )
}

// CreateServer makes it easy to quickly create a server
func CreateServer(id int, port string, startState int) *Server {
	server := new(Server)
	server.ID = id
	server.Term = 0
	server.Port = port
	server.State = startState
	server.Hb = make(chan int)
	server.VotedFor = -1
	server.VoteRequested = make(chan bool)
	server.TotalVotes = []bool{false, false, false, false, false}
	server.AliveServers = []bool{true, true, true, true, true}
	server.NumAliveServers = numServers
	return server
}

// Run is the main loop of the server, which starts by activating the server,
// and looping it through timeouts.
func Run(s *Server) {
	// Seed random for timeout
	rand.Seed(time.Now().UTC().UnixNano())

	fmt.Printf("%v: server run", s.ID)

	go func() {
		address, err := net.ResolveTCPAddr("tcp", s.Port)
		
		if err != nil {
			log.Print(err)
		}

		inbound, err := net.ListenTCP("tcp", address)
		if err != nil {
			log.Print(err)
		}
		rpc.Register(s)
		rpc.Accept(inbound)

		fmt.Printf("server %v added \n", s.ID)
	}()

	// main loop
	for {

		fmt.Printf("%v: in for  state:%d \n", s.ID , s.State)

		switch s.State {
			// Server is a follower
			case Follower:
				// Wait for heartbeat request
				if !RandomTimeout(s) {
					// Switch State
					s.State = Candidate
				}

			// Server is a candidate for leader
			case Candidate:
				// Start an election
				StartElection(s)

			// Server is a leader
			case Leader:
				fmt.Printf("%v: is leader\n", s.ID)
				// Get heartbeat from all servers
				GetHeartbeats(s)
		}
	}
}
