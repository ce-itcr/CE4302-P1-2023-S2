package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	restfulapi "Backend/components/RESTfulAPI"
)

func trimNewline(s string) string {
	return s[:len(s)-1]
}

func main() {
	go restfulapi.Restfulapi()
	fmt.Println("Enter 'C ' (c and space) to end the server")

	reader := bufio.NewReader(os.Stdin)
PELoop:
	for {
		fmt.Print("\nEnter a command: ")
		command, _ := reader.ReadString('\n')
		command = trimNewline(command)

		args := strings.Split(command, " ")
		if len(args) < 1 {
			fmt.Println("Invalid command. Please enter 'C'.")
			continue
		}

		switch args[0] {
		case "C":
			fmt.Println("C pressed")
			restfulapi.Close()
			fmt.Println("Terminate server")
			break PELoop
		}
	}
}
