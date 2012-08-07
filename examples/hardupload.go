package main

import (
    "fmt"
    "io/ioutil"
    m2go "github.com/Bren2010/m2go"
)

var (
    senderId = "82209006-86FF-4982-B5EA-D1E29E55D481"
    recvSpec = "tcp://127.0.0.1:9999"
    sendSpec = "tcp://127.0.0.1:9998"
)

func main() {
    conn, _ := m2go.Connect(senderId, recvSpec, sendSpec)
    defer conn.Close()
    
    response := ""
    for true {
        req, _ := conn.Read()
        if req.IsDisconnect() { fmt.Println("Disconnect.");continue }
        
        done, dOk := req.Headers["x-mongrel2-upload-done"]
        start, sOk := req.Headers["x-mongrel2-upload-start"]
        
        if dOk {
            if start != done {
                fmt.Println("Got the wrong target file:  ", start, done)
                continue
            }
            
            body, _ := ioutil.ReadFile(done)
            fmt.Printf("Done:  Body is %d long, content-length is %s\n",
                len(body), req.Headers["content-length"])
            
            response = "Good!"
        } else if sOk {
            fmt.Println("Upload starting, don't reply yet.")
            fmt.Println("Will read file from", start)
            continue
        } else {
            response = "Hello, there!"
        }
        
        req.Respond(response + "\r\n")
    }
}
