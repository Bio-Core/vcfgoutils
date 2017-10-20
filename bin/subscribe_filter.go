package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"github.com/nats-io/go-nats"
)

// usage
// A function that returns the usage of the program.
func usage() {
	log.Fatalf("Usage: subscribe_filter [-s nats://<servername>:<port>] <subject> \n")
}

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
}

func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var showTime = flag.Bool("t", false, "Display timestamps")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	nc, err := nats.Connect(*urls)
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}

	subj, i := args[0], 0
	var byteMessage = ""
	nc.Subscribe(subj, func(msg *nats.Msg) {
		i++
		//printMsg(msg, i)
		byteMessage = string(msg.Data)
		fmt.Println(byteMessage)
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%s]\n", subj)
	if *showTime {
		log.SetFlags(log.LstdFlags)
	}

	runtime.Goexit()
}
