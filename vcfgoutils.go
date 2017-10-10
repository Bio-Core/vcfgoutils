package vcfgoutils

import (
	"encoding/json"
	"log"
	"os"

	"github.com/brentp/vcfgo"
	"github.com/nats-io/nats"
)

// SimpleGermlineMutation
// A general struct containing the desired information from a VCF file.
type SimpleGermlineMutation struct {
	Chromosome   string
	Position     uint64
	Reference    string
	Alternate    string
	Quality      float32
	Filter       string
	AllelicDepth int
	Coverage     int
}

// ConvertVcfToJSON
// Convert a VCF data into JSON format for streaming using
// NATS
func ConvertVcfToJSON(variant *vcfgo.Variant) SimpleGermlineMutation {
	altdepths, _ := variant.Samples[0].AltDepths()

	// JSON format for variant call information
	simpleVariant := SimpleGermlineMutation{
		variant.Chromosome,
		variant.Pos,
		variant.Reference,
		variant.Alternate[0],
		variant.Quality,
		variant.Filter,
		variant.Samples[0].DP,
		altdepths[0],
	}

	return simpleVariant
}

// SendVcfToNatsAsJSON
// A wrapper function for opening a VCF file, converting to JSON format,
// and transmitting to a NATS server.
func SendVcfToNatsAsJSON(urls *string, vcfFile *string) {
	// setup the nats connection
	log.Println("Connecting to server: ", *urls)
	nc, err := nats.Connect(*urls)
	if err != nil {
		log.Fatalf("Cannot connect: %v\n", err)
	}
	defer nc.Close()

	// prepare the VCF file, read through each entry and parse out the relevant information
	// and stream the data to the NATS server
	f, _ := os.Open(*vcfFile)
	rdr, err := vcfgo.NewReader(f, false)
	if err != nil {
		panic(err)
	}
	for {
		variant := rdr.Read()
		if variant == nil {
			break
		}

		simpleVariant := ConvertVcfToJSON(variant)
		b, _ := json.Marshal(simpleVariant)

		// Send the data to the NATS server
		nc.Publish("queue1", []byte(b))
		nc.Flush()
	}
}
