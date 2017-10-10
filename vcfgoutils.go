package vcfgoutils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/brentp/vcfgo"
	"github.com/nats-io/nats"
)

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

func ConvertVcfToJson(variant *vcfgo.Variant) SimpleGermlineMutation {
	altdepths, _ := variant.Samples[0].AltDepths()

	// JSON format for variant call information
	simple_variant := SimpleGermlineMutation{
		variant.Chromosome,
		variant.Pos,
		variant.Reference,
		variant.Alternate[0],
		variant.Quality,
		variant.Filter,
		variant.Samples[0].DP,
		altdepths[0],
	}

	return simple_variant
}

func SendVcfToNatsAsJson(urls *string, vcf_file *string) {
	// setup the nats connection
	fmt.Println("Connecting to server: ", *urls)
	nc, err := nats.Connect(*urls)
	if err != nil {
		log.Fatalf("Cannot connect: %v\n", err)
	}
	defer nc.Close()

	// prepare the VCF file, read through each entry and parse out the relevant information
	// and stream the data to the NATS server
	f, _ := os.Open(*vcf_file)
	rdr, err := vcfgo.NewReader(f, false)
	if err != nil {
		panic(err)
	}
	for {
		variant := rdr.Read()
		if variant == nil {
			break
		}

		simple_variant := ConvertVcfToJson(variant)
		b, _ := json.Marshal(simple_variant)

		// Send the data to the NATS server
		nc.Publish("queue1", []byte(b))
		nc.Flush()
	}
}
