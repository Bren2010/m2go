package m2go

import (
    "encoding/json"
    "strconv"
)

type Request struct {
    connection *Connection
    
    UUID string
    ID string
    Path string
    Headers map[string] string
    Body string
    Data map[string] interface{}
}

// Send data to the request.  This will leave the transaction hanging!
func (r *Request) Send(body string) {
    size := strconv.Itoa(len(r.ID))
    resp := r.UUID + " " + size + ":" + r.ID + ", " + body
    r.connection.send.Send([]byte(resp), 0)
}

// Respond to a request and kill the connection.
func (r *Request) Respond(body string) {
    r.Send(body)
    r.Kill()
}

// Respond to a request with Json
func (r *Request) RespondJson(data interface{}) {
    body, err := json.Marshal(data)
    if err != nil { return }
    
    r.Respond(string(body))
}

// Respond to a request with valid HTTP response.
func (r *Request) RespondHTTP(body, code, status string, headers map[string]string) {
    r.Send(httpResponse(body, code, status, headers))
}

// Terminate a connection.
func (r *Request) Kill() {
    r.Send("")
}

// Returns true if the packet is a disconnect.
func (r *Request) IsDisconnect() bool {
    if r.Headers["METHOD"] != "JSON" { return false }
    return r.Data["type"] == "disconnect"
}
