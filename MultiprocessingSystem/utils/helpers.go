package utils

import ()

// Request structure for the PE - CacheController communication
type RequestM1 struct {
    Type    string // WRITE or READ operation
    Address int    // The address to READ or WRITE from
    Data    int    // (Only for WRITE) The data to store
}

// Response structure for the PE - CacheController communication
type ResponseM1 struct {
    Status bool   // Status to know if the request was successful
    Type   string // WRITE or READ operation
    Data   int    // (Only for READ) The data to store in the register
}

// Request structure for the CacheController - Interconnect communication
type RequestM2 struct {
    Type    string // WRITE or READ operation
    Address int    // The address to READ or WRITE from
    Data    int    // (Only for WRITE) The data to store
}

// Response structure for the CacheController - Interconnect communication
type ResponseM2 struct {
    Status bool   // Status to know if the request was successful
    Type   string // WRITE or READ operation
    Data   int    // (Only for READ) The data to store in the register
}
