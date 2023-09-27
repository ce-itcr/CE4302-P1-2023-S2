package utils

// Structure for read/write requests
type Request struct {
    Type string                     // WRITE or READ operation
    Address int                     // The address to READ or WRITE from
    Data int                        // (Only for WRITE) The data to store
}

type Response struct {
    Status bool                     // Status to know if the request was successfull
    Type string                     // WRITE or READ operation
    Data int                        // (Only for READ) The data to store in the register
}

// Structure for read/write requests
type RequestMem struct {
    Type string                     // WRITE or READ operation
    Address int                     // The address to READ or WRITE from
	Data int
}

type ResponseMem struct {
    Status bool                     // Status to know if the request was successfull
    Data int                        // (Only for READ) The data to store in the register
	Address int
	StatusData string
}

