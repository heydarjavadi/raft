package main

import(
    "fmt"
	"log"
    "net/rpc"
    "time"
)

// SendHeartbeatRequest creates and sends a message
func SendHeartbeatRequest(source *Server, destination *Server) {
    var mes = new(AppendEntries)
    mes.SourceID = source.ID
    mes.Source = source.Port
    mes.Destination = destination.Port
    mes.Term = source.Term
    mes.NumServers = source.NumAliveServers
    mes.ServerStatus = source.AliveServers

    // send response
    client, err := rpc.Dial("tcp", mes.Destination)
    if err != nil {
        // Fail silently
        //log.Print(err)
        fmt.Printf("No response from %v\n", destination.ID)
        if source.AliveServers[destination.ID] == true {
          source.AliveServers[destination.ID] = false
          source.NumAliveServers--
        }
        return
    }
    
    var reply = new(AppendEntries)
    err = client.Call("Server.Heartbeat", mes, reply)
    if err != nil {
        log.Print(err)
    }else{
        if(reply.Success){
            
            if source.AliveServers[destination.ID] == false {
                source.AliveServers[destination.ID] = true
                source.NumAliveServers++
            }

            fmt.Printf("Heartbeat from %d, Term %d, Num Servers %d\n",
            destination.ID, reply.Term, source.NumAliveServers)
        }
        if(reply.Term > source.Term){
            source.State = Follower
            return
        }
    }
}

// GetHeartbeats sends a heartbeat to all servers, and requests a heartbeat from
// all servers
func GetHeartbeats(s *Server) {
    for _, val := range s.Servers {
        if val.ID != s.ID {
            //fmt.Printf("%v, calling %v at %v\n",s.ID, val.ID, val.Port)
            go SendHeartbeatRequest(s, val)
        }
    }
    time.Sleep(time.Second)
}
