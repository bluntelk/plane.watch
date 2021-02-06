package main

import (
	"bufio"
	"log"
	"net"
	"time"
)

type Dump1090Reader struct {
	commandChan chan int
	messageChan chan string
	host, port  string

	connected   bool

	handler     func(string)
}

func NewDump1090Reader(host, port string) *Dump1090Reader {
	d := new(Dump1090Reader)
	d.host = host
	d.port = port
	d.commandChan = make(chan int)
	d.messageChan = make(chan string)
	return d
}

func (d *Dump1090Reader) SetHandler(f func(string)) {
	d.handler = f
}

func (d *Dump1090Reader) isConnected() bool {
	return d.connected
}

func (d *Dump1090Reader) Stop() {
	if d.isConnected() {
		d.commandChan <- 1
	}
}

func (d *Dump1090Reader) Connect() error {
	log.Printf("Connecting to dump1090 @ %s:%s", d.host, d.port)
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(d.host, d.port), 5 * time.Second)
	if nil != err {
		return err
	}
	d.connected = true

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			d.messageChan <- scanner.Text()
		}
	}()

	go func() {
		var txt string
		for {
			select {
			case <-d.commandChan:
				d.connected = false
				conn.Close()
				return
			case txt = <-d.messageChan:
			// we have a message we need to process
				if d.handler != nil {
					d.handler(txt)
				}
			}
		}
	}()

	return nil
}
