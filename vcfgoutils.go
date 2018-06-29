package vcfgoutils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/brentp/vcfgo"
	"github.com/nats-io/nats"
	mgo "gopkg.in/mgo.v2"
)

// SimpleMutation is a struct containing the desired information from a VCF
// file.
type SimpleMutation struct {
	Chromosome    string
	Position      uint64
	Reference     string
	Alternate     string
	Quality       float32
	Filter        string
	Genotype      string
	FilteredDepth int
	RefDepth      int
	AltDepth      int
}

// ConvertVcfToJSON converts VCF data into JSON format for streaming using
// NATS
func ConvertVcfToJSON(variant *vcfgo.Variant) SimpleMutation {
	altdepths, _ := variant.Samples[0].AltDepths()
	refdepth, _ := variant.Samples[0].RefDepth()
	genotypeString := GetGenotypeFromSample(variant)

	// JSON format for variant call information
	simpleVariant := SimpleMutation{
		variant.Chromosome,
		variant.Pos,
		variant.Reference,
		variant.Alternate[0],
		variant.Quality,
		variant.Filter,
		genotypeString,
		variant.Samples[0].DP,
		refdepth,
		altdepths[0],
	}

	return simpleVariant
}

// SendVcfToNatsAsJSON is a wrapper function for opening a VCF file,
// converting to JSON format and transmitting to a NATS server.
func SendVcfToNatsAsJSON(urls *string, vcfFile *string, subject *string) {
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
		log.Fatalln(err)
		panic(err)
	}
	log.Printf("Submitting VCF file \"%s\" to NATS server...\n", *vcfFile)
	for {
		variant := rdr.Read()
		if variant == nil {
			break
		}

		simpleVariant := ConvertVcfToJSON(variant)
		b, _ := json.Marshal(simpleVariant)

		// Send the data to the NATS server
		nc.Publish(*subject, []byte(b))
		nc.Flush()
	}
}

// ConvertJSONToVCF is a wrapper function for converting JSON streamed data
// back to a VCF entry. The function subscribes to a NATS server and reads the
// streamed byte data and returns the genomic data.
func ConvertJSONToVCF(variant *vcfgo.Variant) {
	fmt.Println("Converting JSON to genomic data...")
	//var subscribedVariant []SimpleMutation
	//return json.Unmarshal(variant, &subscribedVariant)
}

// GetGenotypeFromSample is a function that extracts the genotype listed as
// 0/1 or 1/1 indicating reference or alternate for each allele.
func GetGenotypeFromSample(variant *vcfgo.Variant) string {
	sample := variant.Samples[0]
	genotypeString := []string{strconv.Itoa(sample.GT[0]), strconv.Itoa(sample.GT[1])}
	return strings.Join(genotypeString, "/")
}

// GetSampleIndex is a function that returns the index of a sample given the
// sample name.
func GetSampleIndex(sampleList []string) {
	for sampleIndex, sampleName := range sampleList {
		fmt.Println(sampleIndex, "\t", sampleName)
	}
}

// InsertVCFIntoMongoDB is a wrapper function for injecting JSON data converted
// from VCF.
func InsertVCFIntoMongoDB(session *mgo.Session, data SimpleMutation, database string, collection string) {
	// validate arguments
	if database == "" {
		log.Fatalln("Must provide database name")
	}
	if collection == "" {
		log.Fatalln("Must provide a collection")
	}
	c := session.DB(database).C(collection)
	err := c.Insert(data)
	if err != nil {
		log.Fatal(err)
	}
}
