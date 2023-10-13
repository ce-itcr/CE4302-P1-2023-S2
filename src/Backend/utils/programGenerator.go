package utils

import (
	"time"
	"math/rand"
	"os"
	"fmt"
)

// Instruction represents a core instruction
type Instruction struct {
	Type    string // INC, READ, WRITE
	Address int    // Memory address for READ and WRITE instructions
}

// GenerateRandomInstructions generates a set of random instructions for each core
func GenerateRandomInstructions(numCores, numInstructions int) [][]Instruction {
	rand.Seed(time.Now().UnixNano())

	// Memory has 16 entries
	memoryAddresses := make([]int, 16)
	for i := range memoryAddresses {
		memoryAddresses[i] = i
	}

	// Generate random instructions for each core
	instructionsPerCore := make([][]Instruction, numCores)

	for coreID := 0; coreID < numCores; coreID++ {
		for i := 0; i < numInstructions; i++ {
			var inst Instruction

			// Randomly choose the type of instruction
			switch rand.Intn(3) {
			case 0:
				inst.Type = "INC"
			case 1:
				inst.Type = "READ"
				inst.Address = getRandomMemoryAddress(memoryAddresses)
			case 2:
				inst.Type = "WRITE"
				inst.Address = getRandomMemoryAddress(memoryAddresses)
			}
			instructionsPerCore[coreID] = append(instructionsPerCore[coreID], inst)
		}
	}

	return instructionsPerCore
}

// getRandomMemoryAddress returns a random memory address from the given list
func getRandomMemoryAddress(addresses []int) int {
	return addresses[rand.Intn(len(addresses))]
}

// WriteInstructionsToFile writes instructions to a text file
func WriteInstructionsToFile(filename string, instructions []Instruction) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for i, inst := range instructions {
		line := fmt.Sprintf(inst.Type)
		if inst.Type != "INC" {
			line += fmt.Sprintf(" %d", inst.Address)
		}

		// Add newline if it's not the last line
		if i < len(instructions)-1 {
			line += "\n"
		}

		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}

	return nil
}
