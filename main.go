package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os/exec"
	"time"
)

type Player struct {
	Team string `json:"team"`
	Home string `json:"home"`
	From string `json:"from"`
}

type Message struct {
	Cmd        string `json:"cmd"`
	Hostname   string `json:"hostname"`
	Contestant Player `json:"player"`
}

type Info struct {
	team       string
	hostname   string
	prev       string
	next       string
	challenger *Player
}

var messages []string

func listen(hostname string, chInfo chan Info) {
	if ln, err := net.Listen("tcp", hostname); err == nil {
		defer ln.Close()
		fmt.Println("Listening...")
		for {
			if cn, err := ln.Accept(); err == nil {
				go handle(cn, chInfo)
			}
		}
	}
}

func handle(cn net.Conn, chInfo chan Info) {
	defer cn.Close()
	fmt.Printf("Connection accepted from %s\n", cn.RemoteAddr())
	msg := &Message{}
	dec := json.NewDecoder(cn)
	if err := dec.Decode(msg); err == nil {
		//fmt.Println(msg)
		switch msg.Cmd {
		case "jump":
			info := <-chInfo
			fmt.Println(info, msg)
			enc := json.NewEncoder(cn)
			if err := enc.Encode(Message{Cmd: "ok"}); err != nil {
				fmt.Printf("Can't encode OK REPLY\n%s\n", err)
			}
			player := msg.Contestant
			if info.challenger != nil {
				var loser Player
				if rand.Intn(100) >= 50 {
					loser = player
					player = *info.challenger
				} else {
					loser = *info.challenger
				}
				send(loser.Home, Message{Cmd: "send new", Hostname: info.hostname},
					func(cn net.Conn) {})
			}
			if info.next == "" || info.prev == "" {
				fmt.Printf("Ganaron los del equipo %s\n", player.Team)
				return
			}
			var remote string
			if player.From == info.prev {
				remote = info.next
			} else {
				remote = info.prev
			}
			player.From = info.hostname
			needToFreeInfo := true
			send(remote, Message{"jump", info.hostname, player}, func(cn2 net.Conn) {
				duration := time.Second * 3
				if err := cn2.SetReadDeadline(time.Now().Add(duration)); err != nil {
					fmt.Printf("SetReadDeadline failed:\n%s\n", err)
					panic("OMG!")
				}

				dec := json.NewDecoder(cn2)
				msg2 := &Message{}
				if err := dec.Decode(msg2); err == nil {
					fmt.Println("Se supone que recibimos OK")
				} else {
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						if msg.Hostname == info.next {
							info.challenger = &msg.Contestant
							fmt.Println("liberando ando info func send jump")
							needToFreeInfo = false
							chInfo <- info
						}
						fmt.Printf("read timeout:\n%s\n", err)
					} else {
						fmt.Printf("read error:\n%s\n", err)
					}
				}
			})
			fmt.Println("liberando ando info end of jump")
			if needToFreeInfo {
				chInfo <- info
			}
		case "send new":
			info := <-chInfo
			var remote string
			player := Player{Home: info.hostname, From: info.hostname}
			if info.prev == "" {
				remote = info.next
				player.Team = "Cobras"
			} else {
				remote = info.prev
				player.Team = "Leones"
			}
			fmt.Printf("Sending new player from %s\n", info.team)
			needToFreeInfo := true
			send(remote, Message{"jump", info.hostname, player}, func(cn2 net.Conn) {
				duration := time.Second
				if err := cn2.SetReadDeadline(time.Now().Add(duration)); err != nil {
					fmt.Printf("SetReadDeadline failed:\n%s\n", err)
					panic("OMG!")
				}

				dec := json.NewDecoder(cn2)
				msg2 := &Message{}
				if err := dec.Decode(msg2); err == nil {
					fmt.Println("Se supone que recibimos OK")
				} else {
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						if msg.Hostname == info.next {
							info.challenger = &msg.Contestant
							fmt.Println("liberando ando info func send send new")
							needToFreeInfo = false
							chInfo <- info
						}
						fmt.Printf("read timeout:\n%s\n", err)
					} else {
						fmt.Printf("read error:\n%s\n", err)
					}
				}
			})
			fmt.Println("liberando ando info end of send new")
			if needToFreeInfo {
				chInfo <- info
			}
		}
	} else {
		fmt.Printf("Couldn't decode: %s\n", err)
	}
}

func send(remote string, msg Message, f func(cn net.Conn)) {
	if cn, err := net.Dial("tcp", remote); err == nil {
		defer cn.Close()
		enc := json.NewEncoder(cn)
		if err := enc.Encode(msg); err == nil {
			f(cn)
		} else {
			fmt.Printf("Couldn't encode %s\n", err)
		}
	} else {
		fmt.Printf("Failed to send: %s\n", err)
	}
}

func printMessage(info *Info, msg string) {
	fmt.Printf("[%s] %s\n", info.hostname, msg)
	messages = append(messages, msg)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(messages)
}

func postMessage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var msg string
	err := decoder.Decode(&msg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	printMessage(&Info{}, msg)
	w.WriteHeader(http.StatusOK)
}

func main() {
	special := flag.Bool("s", false, "The special flag for testing stuff.")
	hostname := flag.String("h", "", "IP/Hostname:port to listen on.")
	prevRemote := flag.String("p", "", "Previous node. If empty, home team 1")
	nextRemote := flag.String("n", "", "Next node. If empty, home team 2")
	flag.Parse()

	if *special { // TODO assuming localhost only for testing purposes.
		chInfo := make(chan Info)
		if *hostname != "" {
			info := Info{
				hostname: *hostname,
				prev:     *prevRemote,
				next:     *nextRemote,
			}
			chInfo <- info
			go listen(*hostname, chInfo)
		}

		// Execute the provided commands
		commands := []string{
			"go run node.go -h localhost:8003 -p localhost:8002 -n localhost:8004",
			"go run node.go -h localhost:8004 -p localhost:8003 -n localhost:8005",
			"go run node.go -h localhost:8005 -p localhost:8004 -n localhost:8006",
			"go run node.go -h localhost:8006 -p localhost:8005 -n localhost:8007",
			"go run node.go -h localhost:8007 -p localhost:8006 -n localhost:8008",
			"go run node.go -h localhost:8008 -p localhost:8007 -n localhost:8009",
			"go run node.go -h localhost:8009 -p localhost:8008 -n localhost:8010",
			"go run node.go -h localhost:8010 -p localhost:8009",
		}

		for _, command := range commands {
			fmt.Println(command)
			args := []string{"run", "node.go"}
			args = append(args, command[8:]...)
			cmd := exec.Command("go", args...)
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Error executing command: %s\n", err)
			}
		}
	} else {
		fmt.Println("Please use the special flag to run the commands.")
	}

	// Start the HTTP server
	http.HandleFunc("/messages", getMessages)
	http.HandleFunc("/message", postMessage)
	fmt.Println("Starting HTTP server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
