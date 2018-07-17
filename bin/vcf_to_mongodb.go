package main

import (
	"encoding/json"
	"flag"
	"log"
	"runtime"

	"github.com/bio-core/vcfgoutils"
	"github.com/nats-io/go-nats"
	mgo "gopkg.in/mgo.v2"
)

// usage
// A function that returns the usage of the program.
func usage() {
	log.Fatalf("Usage: vcf_to_mongodb [-s nats://<servername>:<port> --mongodb <mongodb-ip>:<mongdb-port>] <subject> \n")
	flag.PrintDefaults()
}

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
}

func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var showTime = flag.Bool("t", false, "Display timestamps")
	var mongoDbPtr = flag.String("mongodb", "localhost:27017", "The MongoDB hostname/IP and port")
	var database = flag.String("database", "vcfdb", "The MongoDB database name to store data to")
	var collection = flag.String("collection", "mutations", "The collection name to store data to")

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
	session, err := mgo.Dial(*mongoDbPtr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("MongoDB connection established on", *mongoDbPtr, "...")
	session.SetMode(mgo.Monotonic, true)
	mongoDatabase := *database
	mongoCollection := *collection
	subj, i := args[0], 0
	var simpleSubMutation vcfgoutils.SimpleMutation
	nc.Subscribe(subj, func(msg *nats.Msg) {
		i++
		err := json.Unmarshal(msg.Data, &simpleSubMutation)
		if err == nil {
			vcfgoutils.InsertVCFIntoMongoDB(session, simpleSubMutation, mongoDatabase, mongoCollection)
		}
		printMsg(msg, i)
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
