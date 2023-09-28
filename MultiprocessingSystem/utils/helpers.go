package utils

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

// Response structure for the MainMemory - Interconnect communication
type RequestMemory struct {
	Status  bool
	Type    string // WRITE or READ operation
	Address int    // Address in memory to access
	Value   uint32 // value para escribir (si es una operación de escritura).
}

// Response structure for the MainMemory - Interconnect communication
type ResponseMemory struct {
	Status bool
	Type   string // WRITE or READ operation
	Value  uint32 // value para escribir (si es una operación de
	Time   int    // Time required to finish the WRITE or READ action
}
