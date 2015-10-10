package main

import (
	"bufio"
	"flag"
	"github.com/stianeikeland/go-rpio"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	pin        = rpio.Pin(4)
	serverMode = false
	listen     string
)

func main() {
	flag.StringVar(&listen, "listen", "/tmp/thermostat.sock", "socket/host")
	flag.Parse()

	args := flag.Args()
	if len(args) > 0 {
		command := args[0]
		if command == "server" {
			runServer()
		} else if command == "on" {
			runSetState(true)
		} else if command == "off" {
			runSetState(false)
		} else {
			log.Fatal("Unrecognised command: %s", command)
		}
	} else {
		runGetState()
	}
}

func runSetState(state bool) {
	connection := newConnection()
	connection.Write([]byte("set:on" + "\n"))
	response, _ := bufio.NewReader(connection).ReadString('\n')
	if response == "on" {
		log.Println("on")
	} else {
		log.Println("off")
	}
}

func runGetState() {
	connection := newConnection()
	connection.Write([]byte("get" + "\n"))
	log.Println("sent get")
	reader := bufio.NewReader(connection)
	response, _ := reader.ReadString('\n')
	log.Println("waiting")
	log.Println("Received response", response)
	if response == "on" {
		log.Println("on")
	} else {
		log.Println("off")
	}
}

func newConnection() net.Conn {
	connection, err := net.DialTimeout("unix", listen, time.Millisecond*500)
	if err != nil {
		log.Fatal(err)
	}

	return connection
}

func runServer() {
	if err := rpio.Open(); err != nil {
		log.Fatal(err)
	}

	defer rpio.Close()
	pin.Output()

	boiler := NewBoiler(pin)
	defer boiler.Stop()

	syscall.Umask(0000)

	l, e := net.Listen("unix", listen)
	if e != nil {
		log.Fatal(e)
	}

	go boiler.RunLoop()
	go func() {
		for {
			conn, _ := l.Accept()
			command, _ := bufio.NewReader(conn).ReadString('\n')
			log.Println("Received command", command)
			if command == "get" {
				_, err := conn.Write([]byte("on" + "\n"))
				if err != nil {
					log.Fatal(err)
				}
				log.Println("sent get state back")
			} else if command == "set" {
				conn.Write([]byte("on" + "\n"))
				log.Println("set to on")
			}
		}
	}()

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	signal := <-ch
	log.Println("Received signal", signal)

	l.Close()
	boiler.Stop()
}
