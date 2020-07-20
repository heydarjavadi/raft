package main

import (
    "fmt"
    "math/rand"
    "time"
)

// RandomTimeout runs every timeout period with a lower and upper bound. These
// bounds can be set in the const section
func RandomTimeout(s *Server) bool {
	// Milliseconds, Min = timeout, max = timeout*2
	sleep := rand.Int() % (timeout*1000) + (timeout*1000)

	for i := sleep; i > 0; i-- {
		// Every 100 milliseconds, check for an update
		if i % 100 == 0 {
			// Check the channel for some input, but don't block on the channel
			select {
				case value := <-s.Hb: //if i received some heartbeat
	        		fmt.Printf("%v: heartbeat received from %v\n", s.ID, value)
					return true
				case <-s.VoteRequested: // if i received some request for vote
					return true
				default:
				// Do nothing for now
			}
		} else {
			time.Sleep(time.Millisecond)
		}
	}
	return false
}

// Heartbeat responds to a hearbeat request
func (s *Server) Heartbeat(message *AppendEntries, response *AppendEntries) error {

	fmt.Printf("%v: heartbeat received \n my state is: %V", s.ID,s.State)
	fmt.Println(message)

	response.Source = s.Port
	response.Destination = message.Source
	response.Term = s.Term

	if(message.Term >= s.Term){
		response.Success = true
		s.Term = message.Term
		s.State = Follower
		s.VotedFor = message.SourceID
		s.Hb <- message.SourceID

		fmt.Printf("Server Status updated\n",s)
		s.NumAliveServers = message.NumServers
		s.AliveServers = message.ServerStatus
	}else{
		response.Success = false
	}
	
	return nil
}

// Elect respond to an election vote, and wait for confirmation or timeout
// A server cannot vote for a leader who has a log index less than their own
func (s *Server) Elect(message *RequestVote, response *RequestVote) error {

	fmt.Printf("%v: elect message received from %d its term is %d\n", s.ID, message.SourceID, message.Term)
	fmt.Println("our voteFor is %d and term is %d", s.VotedFor, s.Term)

	response.SourceID = s.ID
	response.Source = s.Port
	response.Destination = message.Source
	response.Term = s.Term

	if(message.Term > s.Term){
		s.State = Follower
		s.Term = message.Term
		response.Vote = true
		s.VotedFor = message.SourceID
		s.VoteRequested <- true
	}else if(message.Term == s.Term){
		s.VoteRequested <- true
		if(s.VotedFor == -1 || s.VotedFor == message.SourceID){
			response.Vote = true
			s.VotedFor = message.SourceID
		} else {
			response.Vote = false
		}
	}else{
		response.Vote = false
	}

	fmt.Println("our vote is %t and term is %d", response.Vote, s.Term)
	
	return nil
}
