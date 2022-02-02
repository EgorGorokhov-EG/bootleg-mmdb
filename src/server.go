package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	PORT      = "8080"
	HOST      = "localhost"
	CONN_TYPE = "tcp"
)

// Global cache to store key-value pairs
var cache Cache

// Aof to handle persistence
var aof Aof

func main() {
	logger := log.New(os.Stdout, "server-logger", log.LstdFlags)
	listener, err := net.Listen(CONN_TYPE, HOST+":"+PORT)
	defer listener.Close()

	if err != nil {
		logger.Println("Can't start listener on port", PORT, ".", err)
		os.Exit(1)
	}

	// Initialize needed structs
	cache.init()
	aof.init()

	logger.Println("Server started on port:", PORT)

	for {
		connection, err := listener.Accept()
		if err != nil {
			logger.Println("Can't accept a connection:", err)
		}
		handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
	defer connection.Close()
	// TODO: Make it possible to receive queries larger than 1024 bytes
	buffer := make([]byte, 1024)

	_, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Can't read a request:", err)
		return
	}

	query := strings.Trim(string(buffer), "\n")
	queryArray := strings.Split(query, " ")

	response := processQuery(queryArray)

	connection.Write([]byte(response))
}

// returns respond for the given command
func processQuery(query []string) string {

	command := query[0]

	switch command {
	case "set":
		key, value := query[1], query[2]
		fmt.Println("SET for", key, "->", value)
		cache.set(key, value)
		aof.appendLog(strings.Join(query, " "))
		return key + "->" + value + " is set!"

	case "get":
		key := query[1]
		fmt.Println("GET for", key)
		value, ok := cache.get(key)
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

func startPersistence() {
	for {
		time.Sleep(time.Second)
		aof.flushAof()
	}
}

// Aof stores info needed to implement AOF persistence
type Aof struct {
	path    string
	aofBuf  string
	aofLock *sync.Mutex
	logger  *log.Logger
}

func (aof *Aof) init() {
	aof.path = "/tmp/bootleg-mmdb-persistence/aof.pers"
	aof.aofBuf = ""
	aof.aofLock = new(sync.Mutex)
	aof.logger = log.New(os.Stdout, "aof-logger", log.LstdFlags)

	aof.logger.Println("Aof is ready!")
}

func (aof *Aof) appendLog(query string) {
	aof.aofLock.Lock()
	aof.aofBuf = aof.aofBuf + query + "\n"
	aof.aofLock.Unlock()
}

func (aof *Aof) flushAof() {

	// TODO: Somehow handle these file related errors
	aof.aofLock.Lock()
	tempAofBuf := aof.aofBuf
	aof.aofBuf = ""
	aof.aofLock.Unlock()

	aofFile, err := os.OpenFile(aof.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panic(err)
	}

	if _, err := aofFile.Write([]byte(tempAofBuf)); err != nil {
		aofFile.Close()
		log.Panic(err)
	}

	if err := aofFile.Close(); err != nil {
		log.Panic(err)
	}
}
