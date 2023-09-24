package utils

import (
    "fmt"
)

func SayHello(name string) {
    fmt.Printf("Hello, %s!\n", name)
}


type Request struct {
    Type    string // WRITE or READ operation
    Address int    // The address to READ or WRITE from
    Data    int    // (Only for WRITE) The data to store
}

type Response struct {
    Status bool   // Status to know if the request was successful
    Type   string // WRITE or READ operation
    Data   int    // (Only for READ) The data to store in the register
}
