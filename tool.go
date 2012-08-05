package m2go

import (
    "bytes"
    "fmt"
    "strconv"
)

func httpResponse(body, code, status string, headers map[string]string) string {
    http := "HTTP/1.1 %s %s\r\n%s\r\n%s"
    
    headers["Content-Length"] = strconv.Itoa(len(body))
    hd := ""
    for k, v := range headers {
        hd = fmt.Sprintf("%s: %s\r\n%s", k, v, hd)
    }
    
    return fmt.Sprintf(http, code, status, hd, body)
}

func parseNetstring(ns []byte) ([]byte, []byte) {
    lenS := bytes.Index(ns, []byte(":"))
    length, err := strconv.Atoi(string(ns[0:lenS]))
    if err != nil { return []byte(""), ns }
    
    text := ns[lenS + 1:length + lenS + 1]
    
    return text, ns[length + lenS + 2:]
}
