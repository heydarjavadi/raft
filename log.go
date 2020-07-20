package main

import ("fmt")

// Log records a request
type Log struct {
	Value     int
	Term     int
	Committed bool
}

// Log Struct toString
func (l Log) String() string {
    return fmt.Sprintf("{Value:%d, Term: %d, Committed: %t}",l.Value,l.Term, l.Committed)
}

// Commit sets a log to become committed
func (l *Log) Commit() {
	l.Committed = true
}
