package main

import (
	"flag"
	"fmt"
	"log"

	"gitlab.com/uhn/vcfgoutils"
)

func usage() {
	fmt.Println("Usage: submit_vcf_to_nats [-s nats://<server>:<port>] [-f /path/to/vcf/file]")
	flag.PrintDefaults()
}

func main() {
	// setup the command line arguments
	urlPtr := flag.String("s", "", "The NATS server URLs (separated by comma) (required)")
	vcfPtr := flag.String("f", "", "The VCF file to transmit (required)")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()
	if *urlPtr == "" {
		fmt.Println("Missing argument -s")
		flag.Usage()
	} else if *vcfPtr == "" {
		fmt.Println("Missing argument -f")
		flag.Usage()
	} else {
		vcfgoutils.SendVcfToNatsAsJSON(urlPtr, vcfPtr)
	}
}
