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

func logMessage(source, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s][%s] %s\n", timestamp, source, message)
}

type Command struct {
	Identifier int    `json:"Identifier"`
	Message    string `json:"Message"`
	Name       string `json:"Name"`
}

var rconIP, rconPort, rconPass string
var wsRcon bool
var stopReader chan struct{}
var conn *websocket.Conn

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
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logMessage("SYSTEM", fmt.Sprintf("an error occurred while working with the output stream: %v", err))
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		logMessage("SYSTEM", fmt.Sprintf("an error occurred while working with the input stream: %v", err))
	}

	go handleOutput(stdout, stopReader, "STDOUT")
	go handleOutput(stderr, stopReader, "STDERR")

	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			text = strings.Trim(text, "\n")
			if !wsRcon {
				if text == "quit" {
					cmd.Process.Signal(syscall.SIGTERM)
					os.Exit(1)
				} else {
					logMessage("INPUT", fmt.Sprintf("Unable to run %s due to RCON not being connected yet", text))
				}
			} else {
				sendRconCommand(conn, text)
			}

		}
	}()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-exitSignal
		logMessage("SYSTEM", "Received request to stop the process, stopping the game...")
		cmd.Process.Signal(syscall.SIGTERM)
	}()

	err = cmd.Start()
	if err != nil {
		logMessage("SYSTEM", fmt.Sprintf("failed to start main game process: %v", err))
		return
	}

	poll()

	err = cmd.Wait()
	if err != nil {
		logMessage("SYSTEM", fmt.Sprintf("Main game process exited with error: %v", err))
		os.Exit(1)
	}
}

func handleOutput(pipe io.ReadCloser, stop chan struct{}, source string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		if !wsRcon {
			line := scanner.Text()
			if line == "" {
				continue
			}

			logMessage(source, line)
		}
	}
}

func poll() {
	var err error
	conn, _, err = websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/%s", rconIP, rconPort, rconPass), nil)
	if err != nil {
		logMessage("RCON", "Waiting for RCON to come up...")
		time.Sleep(5 * time.Second)
		poll()
		return
	}

	logMessage("RCON", "Connected to RCON. Generating the map now. Please wait until the server status switches to \"Running\".")
	close(stopReader)
	sendRconCommand(conn, "status")
	wsRcon = true

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				logMessage("RCON", "Connection to server closed.")
				os.Exit(1)
			}

			var command Command
			err = json.Unmarshal(msg, &command)
			if err != nil {
				logMessage("RCON", "Error: Invalid JSON received")
				continue
			}

			if command.Message != "" {
				logMessage("RCON", command.Message)
			}
		}
	}()
}

func sendRconCommand(conn *websocket.Conn, cmd string) {
	command := Command{
		Identifier: -1,
		Message:    strings.Trim(cmd, "\n"),
		Name:       "WebRcon",
	}
	jsonCommand, _ := json.Marshal(command)
	conn.WriteMessage(websocket.TextMessage, jsonCommand)
}
