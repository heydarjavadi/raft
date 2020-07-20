package main

import (
	"fmt"
)
// Message contains the value of the previous event in the log, and the index of
// the previous index.
type RequestVote struct {
	SourceID int
	Source      string
	Destination string
	Term   int
	Vote	bool
}

// Struct toString
func (m RequestVote) String() string {

    return fmt.Sprintf(
    	"{sourceID:%d, Source:%s, Destination:%s \n Term: %d, Vote: %t}",
    	m.SourceID,m.Source,m.Destination, m.Term, m.Vote, 
    )
}

// Message contains the value of the previous event in the log, and the index of
// the previous index.
type AppendEntries struct {
	SourceID int
	Source      string
	Destination string
	Term   int
	Success	bool
	NumServers int
	ServerStatus []bool
}

// Message Struct toString
func (m AppendEntries) String() string {

    return fmt.Sprintf(
    	"{sourceID:%d, Source:%s, Destination:%s \n Term: %d, Success: %t, NumServers: %d\n,ServerStatus: %v }",
    	m.SourceID,m.Source,m.Destination, m.Term, m.Success, m.NumServers, m.ServerStatus, 
    )
}