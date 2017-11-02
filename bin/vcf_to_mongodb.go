package main

import (
	"encoding/json"
	"flag"
	"log"
	"runtime"
	"strings"

	"github.com/nats-io/go-nats"
	"gitlab.com/uhn/vcfgoutils"
	mgo "gopkg.in/mgo.v2"
)

// usage
// A function that returns the usage of the program.
func usage() {
	log.Fatalf("Usage: vcf_to_mongodb [-s nats://<servername>:<port> --mongodb <ip-address> --mongoport <mongodb port>] <subject> \n")
}

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
}

func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var showTime = flag.Bool("t", false, "Display timestamps")
	var mongoDbPtr = flag.String("mongodb", "", "The MongoDB server name/IP")
	var mongoPortPtr = flag.String("mongoport", "27017", "The MongoDB port (default: 27017)")

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

	// open a connection to the MongoDB
	mongoHostArray := []string{*mongoDbPtr, *mongoPortPtr}
	mongoHost := strings.Join(mongoHostArray, ":")
	session, err := mgo.Dial(mongoHost)
	log.Println("MongoDB connection established on", mongoHost, "...")
	if err != nil {
		log.Fatal(err)
	}
	session.SetMode(mgo.Monotonic, true)
	mongoDatabase := "test"
	mongoCollection := "simplemutation"
	subj, i := args[0], 0
	var simpleSubMutation vcfgoutils.SimpleMutation
	nc.Subscribe(subj, func(msg *nats.Msg) {
		i++
		err := json.Unmarshal(msg.Data, &simpleSubMutation)
		if err == nil {
			vcfgoutils.InsertVCFIntoMongoDB(session, simpleSubMutation, mongoDatabase, mongoCollection)
		}
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
