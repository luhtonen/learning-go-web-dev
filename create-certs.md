# Create self-signed certifications

## Generate private key (.key)
Key considerations for algorithm "RSA" ≥ 2048-bit
```
openssl genrsa -out server.key 2048
```

## Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)
```
openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
```

## Simple Golang HTTPS/TLS Server
```
package main

import (
    "io"
    "net/http"
    "log"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
    io.WriteString(w, "hello, world!\n")
}

func main() {
    http.HandleFunc("/hello", HelloServer)
    err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
```