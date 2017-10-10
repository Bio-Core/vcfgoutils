package vcfgoutils

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/brentp/vcfgo"
	"gitlab.com/uhn/vcfgoutils"
)

func TestConvertVcfToJSON(t *testing.T) {
	vcfFile := "examples/vcfgoutils_test.vcf"
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
		simpleVariant := vcfgoutils.ConvertVcfToJSON(variant)
		b, _ := json.Marshal(simpleVariant)
		fmt.Println(b)
	}
}
