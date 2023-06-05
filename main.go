package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

type Command struct {
	Identifier int    `json:"Identifier"`
	Message    string `json:"Message"`
	Name       string `json:"Name"`
}

var rconIP, rconPort, rconPass string
var wsRcon bool
var stopReader chan struct{}
var reader *bufio.Reader

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	rconIP = getEnv("RCON_IP", "127.0.0.1")
	rconPort = getEnv("RCON_PORT", "28018")
	rconPass = getEnv("RCON_PASS", "")
	stopReader = make(chan struct{})
	reader = bufio.NewReader(os.Stdin)
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("an error occurred while working with the output stream: ", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("an error occurred while working with the input stream: ", err)
	}

	go handleOutput(stdout, stopReader)
	go handleOutput(stderr, stopReader)

	go func() {
		for {
			select {
			case <-stopReader:
				return
			default:
				text, _ := reader.ReadString('\n')
				text = strings.Trim(text, "\n")
				if text == "quit" {
					cmd.Process.Signal(syscall.SIGTERM)
				} else {
					fmt.Printf("Unable to run %s due to RCON not being connected yet\n", text)
				}
			}
		}
	}()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-exitSignal
		fmt.Println("Received request to stop the process, stopping the game...")
		cmd.Process.Signal(syscall.SIGTERM)
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	poll()

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Main game process exited with error: %v\n", err)
		os.Exit(1)
	}
}

func handleOutput(pipe io.ReadCloser, stop chan struct{}) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		select {
		case <-stop:
			return
		default:
			if wsRcon {
				line := scanner.Text()
				fmt.Println(line)
			}
		}
	}
}

func poll() {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/%s", rconIP, rconPort, rconPass), nil)
	if err != nil {
		fmt.Println("Waiting for RCON to come up...")
		time.Sleep(5 * time.Second)
		poll()
		return
	}

	fmt.Println("Connected to RCON. Generating the map now. Please wait until the server status switches to \"Running\".")
	sendStatus(reader, conn)
	wsRcon = true
	// close(stopReader)

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Connection to server closed.")
				wsRcon = false
				poll()
			}

			var command Command
			err = json.Unmarshal(msg, &command)
			if err != nil {
				fmt.Println("Error: Invalid JSON received")
				continue
			}

			if command.Message != "" {
				fmt.Println(command.Message)
			}
		}
	}()

	for {
		text, _ := reader.ReadString('\n')
		command := Command{
			Identifier: -1,
			Message:    strings.Trim(text, "\n"),
			Name:       "WebRcon",
		}
		jsonCommand, _ := json.Marshal(command)
		conn.WriteMessage(websocket.TextMessage, jsonCommand)
	}

}

func sendStatus(reader *bufio.Reader, conn *websocket.Conn) {
	command := Command{
		Identifier: -1,
		Message:    strings.Trim("status", "\n"),
		Name:       "WebRcon",
	}
	jsonCommand, _ := json.Marshal(command)
	conn.WriteMessage(websocket.TextMessage, jsonCommand)
}
