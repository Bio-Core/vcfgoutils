package main

import (
	"flag"
	"log"

	"github.com/rdeborja/vcfgoutils"
)

func usage() {
	log.Fatalf("Usage: submit_vcf_to_nats [-s nats://<server>:<port>] [-f /path/to/vcf/file]")
}

func main() {
	// setup the command line arguments
	urlPtr := flag.String("s", "nats://localhost:4222", "The NATS server URLs (separated by comma) (default: nats://localhost:4222)")
	vcfPtr := flag.String("f", "", "The VCF file to transmit (required)")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()
	vcfgoutils.SendVcfToNatsAsJson(urlPtr, vcfPtr)
}
