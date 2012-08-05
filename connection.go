package m2go

import (
    "bytes"
    "encoding/json"
    "strconv"
    "strings"
    zmq "github.com/alecthomas/gozmq"
)

type Connection struct {
    senderId string
    context zmq.Context
    recv zmq.Socket
    send zmq.Socket
}

// Open the sockets to handle Mongrel 2 requests.
func Connect(senderId, recvSpec, sendSpec string) (*Connection, error) {
    context, err := zmq.NewContext()
    if err != nil { return nil, err}
    
    recv, err := context.NewSocket(zmq.PULL)
    if err != nil { return nil, err }
    recv.Connect(recvSpec)
    
    send, err := context.NewSocket(zmq.PUB)
    if err != nil { return nil, err }
    send.Connect(sendSpec)
    send.SetSockOptString(zmq.IDENTITY, senderId)
    
    return &Connection{senderId, context, recv, send}, nil
}

// Reads a new request.
func (c *Connection) Read() (*Request, error) {
    data, err := c.recv.Recv(0)
    if err != nil { return &Request{}, err}
    
    // Extract the easy bits.
    parts := bytes.SplitN(data, []byte(" "), 4)
    
    // Extract headers
    dirtyHeaders := make(map[string] interface{})
    
    hJson, rest := parseNetstring(parts[3])
    err = json.Unmarshal(hJson, &dirtyHeaders)
    
    if err != nil { return &Request{}, err }
    
    // Force headers to be a map[string] string
    headers := make(map[string] string)
    
    for k, v := range dirtyHeaders {
        headers[k] = v.(string)
    }
    
    // Extract body
    body, _ := parseNetstring(rest)
    
    // Extract data if this is Json
    jsonData := make(map[string] interface{})
    if headers["METHOD"] == "JSON" {
        err = json.Unmarshal(body, &jsonData)
        if err != nil { return &Request{}, err }
    }
    
    // Return extracted data
    return &Request{
        connection: c,
        UUID: string(parts[0]),
        ID: string(parts[1]),
        Path: string(parts[2]),
        Headers: headers,
        Body: string(body),
        Data: jsonData,
    }, nil
}

// Deliver a message to a list of ids.
func (c *Connection) Deliver(uuid string, ids []string, body string) {
    if len(ids) == 0 { return } // Prevent silliness.
    
    // Limit number of ids to 128.
    rest := []string{}
    if len(ids) > 128 {
        ids = ids[0:128]
        rest = ids[128:]
    }
    
    // Format the netstring.
    id := strings.Join(ids, " ")
    size := strconv.Itoa(len(id))
    resp := uuid + " " + size + ":" + id + ", " + body
    
    c.send.Send([]byte(resp), 0)
    
    if len(rest) != 0 { c.Deliver(uuid, rest, body) }
}

// Deliver a JSON message to a list of ids.
func (c *Connection) DeliverJson(uuid string, ids []string, data interface{}) {
    body, err := json.Marshal(data)
    if err != nil { return }
    
    c.Deliver(uuid, ids, string(body))
}

// Deliver a HTTP message to a list of ids.
func (c *Connection) DeliverHTTP(uuid string, ids []string, body, code, status string, headers map[string]string) {
    c.Deliver(uuid, ids, httpResponse(body, code, status, headers))
}

// Close the connection.
func (c *Connection) Close() {
    c.context.Close()
    c.recv.Close()
    c.send.Close()
}
