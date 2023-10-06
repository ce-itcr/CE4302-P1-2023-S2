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

type AboutProcessingElementList [] AboutProcessingElement

// Struct to represent the time stamp of a Cache Controller **************************************************
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
	ID     int            	`json:"ID"`
	Status string         	`json:"Status"`
	MemoryAccesses int		`json:"MemoryAccesses"`
	Cache  CacheObjectList 	`json:"Cache"`
}

type AboutCacheControllerList [] AboutCacheController

// Struct to represent the time stamp of a Main Memory **************************************************
type BlockObject struct {
	Address			int    	`json:"Address"`
	Data    		int   	`json:"Data"`
}
type BlockObjectList [] BlockObject

type AboutMainMemory struct {
	Status      string    `json:"Status"`
	Blocks BlockObjectList `json:"Blocks"`
}


// Struct to represent the time stamp of the Interconnect **************************************************
type TransactionObject struct {
	Order				int    	`json:"Order"`
	Transaction    		string   `json:"Transaction"`
}
type TransactionObjectList [] TransactionObject

type LogObject struct {
	Order				int    	`json:"Order"`
	Log    				string   `json:"Log"`
}
type LogObjectList [] LogObject

type AboutInterconnect struct {
	Status      		string    	`json:"Status"`
	Logs LogObjectList		`json:"Logs"`
}

// Object Structure for data refresh
type MultiprocessingSystemState struct {
	PEs AboutProcessingElementList `json:"PEs"`
	CCs AboutCacheControllerList `json:"CCs"`
	IC AboutInterconnect `json:"IC"`
	MM AboutMainMemory
}

// Object Structure for executio results
type MultiprocessingSystemResults struct {
	Transactions 			TransactionObjectList 				`json:"Transactions"`
	PowerConsumption 		float64		`json:"PowerConsumption"`
	CacheMisses				int			`json:"CacheMisses"`
	CacheHits				int			`json:"CacheHits"`
	MemoryAccesses			int			`json:"MemoryAccesses"`
	MissRate				float64		`json:"MissRate"`
	HitRate					float64		`json:"HitRate"`
	ReadRequests			int			`json:"ReadRequests"`
	ReadExclusiveRequests	int			`json:"ReadExclusiveRequest"`
	DataResponses			int			`json:"DataResponses"`
	Invalidates				int			`json:"Invalidates"`
	MemoryReads				int			`json:"MemoryReads"`
	MemoryWrites			int			`json:"MemoryWrites"`
}