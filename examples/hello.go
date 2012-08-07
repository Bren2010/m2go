package main

import (
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
    
    for true {
        req, _ := conn.Read()
        if req.IsDisconnect() { continue }
        
        req.Respond("Hello, World!")
    }
}
