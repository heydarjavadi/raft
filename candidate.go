package main

import (
    "fmt"
    "net/rpc"
	"log"
	"time"
)

// StartElection is called when a server times out.
func StartElection(s *Server) {
	fmt.Printf("%v: election started\n", s.ID)
	s.Term = s.Term + 1
	s.VotedFor = s.ID
	s.TotalVotes = []bool{false, false, false, false, false}
	s.TotalVotes[s.ID] = true
	for _, val := range s.Servers {
		// Let's assume he votes for himself
		if val.ID != s.ID {
			RequestForVote(s, val)
		}
	}

	time.Sleep(time.Second)

	if (s.State == Candidate){
		fmt.Printf("%v alive servers\n", s.NumAliveServers)

		if x := CheckVotes(s); x > s.NumAliveServers/2 {
			s.State = Leader
			return
		}else{
			fmt.Println("%v i cant take majority votes %v",s.ID, x)
		}	
	}else{
		fmt.Println("%v state changed to %v",s.ID, s.State)
	}
}

// CheckVotes for election win
func CheckVotes(s *Server) int {
	votes := 0
	fmt.Println("candidate checkVotes %v", s.TotalVotes)
	for _, val := range s.TotalVotes {
		if val {
			votes++
		}
	}
	return votes
}

// RequestForVote sends a request for a vote from destination
func RequestForVote(source *Server, destination *Server) {
	var mes = new(RequestVote)
	mes.SourceID = source.ID
	mes.Source = source.Port
	mes.Destination = destination.Port
	mes.Term = source.Term

	// send response
	client, err := rpc.Dial("tcp", mes.Destination)
	if err != nil {
		fmt.Printf("Cannot connect to %v for vote\n",destination.ID)
		log.Print(err)
		if source.AliveServers[destination.ID] == true {
			source.AliveServers[destination.ID] = false
			source.NumAliveServers--
		}
		return
	}

	var reply = new(RequestVote)
	err = client.Call("Server.Elect", mes, reply)
	if err != nil {
		fmt.Println("vote reply error ",err)
	}else{
		
		if source.AliveServers[destination.ID] == false {
			source.AliveServers[destination.ID] = true
			source.NumAliveServers++
		}

		if(reply.Term > source.Term){
			source.State = Follower
			return
		}

		source.TotalVotes[destination.ID] = reply.Vote
	   	fmt.Printf("%d: %t vote from %d   reply: %v \n", source.ID,reply.Vote, destination.ID,reply)
	}

}
