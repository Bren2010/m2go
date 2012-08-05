m2go
====
m2go is a [Mongrel2](http://mongrel2.org/) handler for [Google's Go](http://golang.org/).

How To Install
--------------
```bash
go get github.com/alecthomas/gozmq
go get github.com/Bren2010/m2go
```

Example
-------
```go
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
        packet, _ := conn.Read()
        if packet.IsDisconnect() { continue }
        
        packet.Respond("Hello, World!")
    }
}
```

License
-------
Same as Mongrel2.
