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

	fmt.Println("Listening on port:", PORT)

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

	requestMsg := strings.Trim(string(buffer), "\n")
	splitRequest := strings.Split(requestMsg, " ")

	response := handleCommand(splitRequest)

	connection.Write([]byte(response))
}

// returns respond to the given command
func handleCommand(arguments []string) string {
	switch arguments[0] {
	case "set":
		fmt.Println("SET for", arguments[1], ":", arguments[2])
		cache[arguments[1]] = arguments[2]
		fmt.Println(cache)
		return arguments[1] + ":" + arguments[2] + " is set!"
	case "get":
		fmt.Println("GET for", arguments[1])
		fmt.Println(cache[arguments[1]])
		val, ok := cache[arguments[1]]

		if ok == false {
			fmt.Println("Can't find the value for the given key!")
			return "null"
		} else {
			return val
		}
	default:
		return "The unknown command"
	}
}
