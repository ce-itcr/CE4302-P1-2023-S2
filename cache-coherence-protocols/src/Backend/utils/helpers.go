package utils

import ()

// Request structure for the PE - CacheController communication
type RequestProcessingElement struct {
    Type    string // WRITE or READ operation
    Address int    // The address to READ or WRITE from
    Data    int    // (Only for WRITE) The data to store
}

// Response structure for the PE - CacheController communication
type ResponseProcessingElement struct {
    Status bool   // Status to know if the request was successful
    Data   int    // (Only for READ) The data to store in the register
}

// Request structure for the CacheController - Interconnect communication
type RequestInterconnect struct {
    Type    string
    AR      string
    Address int    // The address to READ or WRITE from
    Data    int    // (Only for WRITE) The data to store
}

// Response structure for the CacheController - Interconnect communication
type ResponseInterconnect struct {
    Data    int    // (Only for READ) The data to store in the register
    NewStatus string
}

// Request structure for the Interconnect - Cache Controller
type RequestBroadcast struct {
    Type string 
    Address int    
}

// Response structure for the Interconnect - Cache Controller
type ResponseBroadcast struct {
    Match bool
    Status string
    Address string 
    Data int 
}

// Request structure for the Interconnect - Main Memory communication
type RequestMainMemory struct {
	Type    string // WRITE or READ operation
	Address int    // Address in memory to access
	Value   uint32 // value para escribir (si es una operación de escritura).
}

// Response structure for the Iterconnect - Main Memory communication
type ResponseMainMemory struct {
	Status bool
	Type   string // WRITE or READ operation
	Value  uint32 // value para escribir (si es una operación de
	Time   int    // Time required to finish the WRITE or READ action
}
