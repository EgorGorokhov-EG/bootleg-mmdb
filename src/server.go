package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	PORT      = "8080"
	HOST      = "localhost"
	CONN_TYPE = "tcp"
)

var cache map[string]string

func main() {
	listener, err := net.Listen(CONN_TYPE, HOST+":"+PORT)
	defer listener.Close()

	if err != nil {
		fmt.Println("Can't start listener on port", PORT, ".", err)
		os.Exit(1)
	}

	initCache()

	fmt.Println("Server started on port:", PORT)

	for {
		connection, err := listener.Accept()

		if err != nil {
			fmt.Println("Can't accept a connection:", err)
		}

		handleConnection(connection)
	}
}

func initCache() {
	cache = make(map[string]string)

	fmt.Println("The cache created!")
}

func handleConnection(connection net.Conn) {
	defer connection.Close()
	buffer := make([]byte, 1024)

	_, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Can't read a request:", err)
		return
	}

	command := strings.Trim(string(buffer), "\n")
	splitCommand := strings.Split(command, " ")

	response := handleCommand(splitCommand)

	connection.Write([]byte(response))
}

// returns respond to the given command
func handleCommand(arguments []string) string {

	command := arguments[0]

	switch command {
	case "set":
		key, value := arguments[1], arguments[2]
		fmt.Println("SET for", key, "->", value)
		cache[key] = value
		return key + "->" + value + " is set!"

	case "get":
		key := arguments[1]
		fmt.Println("GET for", key)
		value, ok := cache[key]
		if ok == false {
			fmt.Println("Can't find the value for the given key!")
			return "null"
		} else {
			return value
		}
	default:
		return "Unknown command!"
	}
}
