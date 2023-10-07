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

// Queue Struct for int ******************************************************************************************************************
type Queue struct {
	items []int
}

// Enqueue adds an item to the end of the queue.
func (q *Queue) Enqueue(item int) {
	q.items = append(q.items, item)
}

// Dequeue removes and returns the item from the front of the queue.
func (q *Queue) Dequeue() int {
	if len(q.items) == 0 {
		return 0
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

// IsEmpty returns true if the queue is empty.
func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}

// Size returns the number of items in the queue.
func (q *Queue) Size() int {
	return len(q.items)
}

// Queue Struct for string ******************************************************************************************************************
type QueueS struct {
	Items []string
}

// Enqueue adds an item to the end of the queue.
func (q *QueueS) Enqueue(item string) {
	q.Items = append(q.Items, item)
}

// Dequeue removes and returns the item from the front of the queue.
func (q *QueueS) Dequeue() string {
	if len(q.Items) == 0 {
		return ""
	}
	item := q.Items[0]
	q.Items = q.Items[1:]
	return item
}

// IsEmpty returns true if the queue is empty.
func (q *QueueS) IsEmpty() bool {
	return len(q.Items) == 0
}

// Size returns the number of items in the queue.
func (q *QueueS) Size() int {
	return len(q.Items)
}

// Queue Struct for int ******************************************************************************************************************
