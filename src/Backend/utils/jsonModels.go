package utils


// Struct to represent the time stamp of a Processing Element **************************************************
type InstructionObject struct {
	Position		int    	`json:"Position"`
	Instruction    	string    `json:"Instruction"`
}
type InstructionObjectList [] InstructionObject

type AboutProcessingElement struct {
	ID          int       `json:"ID"`
	Register    int       `json:"Register"`
	Status      string    `json:"Status"`
	Instructions InstructionObjectList `json:"Instructions"`
}

// CacheObject represents the structure of the cache objects in the JSON ***************************************
type CacheObject struct {
	Block   int    `json:"Block"`
	Address int    `json:"Address"`
	Data    int    `json:"Data"`
	State   string `json:"State"`
}

// CacheObjectList represents the list of cache objects in the JSON
type CacheObjectList []CacheObject

// MainObject represents the main structure in the JSON
type AboutCacheController struct {
	ID     int            `json:"ID"`
	Status string         `json:"Status"`
	Cache  CacheObjectList `json:"Cache"`
}

// Struct to represent the time stamp of a Processing Element **************************************************
type BlockObject struct {
	Address			int    	`json:"Address"`
	Data    		int   	`json:"Data"`
}
type BlockObjectList [] BlockObject

type AboutMainMemory struct {
	Status      string    `json:"Status"`
	Blocks BlockObjectList `json:"Blocks"`
}