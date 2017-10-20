package main

import (
    "os"
    "fmt"

    "github.com/brentp/vcfgo"
    "gitlab.com/uhn/vcfgoutils"
)

func main() {
    vcfFile := "testfiles/vcfgoutils_test.vcf"
    f, _ := os.Open(vcfFile)
    rdr, err := vcfgo.NewReader(f, false)
    if err != nil {
        panic(err)
    }
    for {
        variant := rdr.Read()
        if variant == nil {
            break
        }

        vcfgoutils.GetSampleIndex(variant.Header.SampleNames)
        fmt.Println(variant.Info)

        // fmt.Println(vcfgoutils.GetSampleIndex(variant.Header.SampleNames))
        // for index, sampleName := range variant.Header.SampleNames {
        //     fmt.Println(index, "\t", sampleName)
        // }
        // genotypeData := vcfgoutils.GetGenotypeFromSample(variant)
        // fmt.Println(genotypeData)
        // simpleVariant := vcfgoutils.ConvertVcfToJSON(variant)
        // fmt.Println(simpleVariant)
    }
}
