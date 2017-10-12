package vcfgoutils

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/brentp/vcfgo"
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
		simpleVariant := ConvertVcfToJSON(variant)
		variantJSON, _ := json.Marshal(simpleVariant)
		expectedString := "{\"Chromosome\":\"chr4\",\"Position\":1806181,\"Reference\":\"C\",\"Alternate\":\"T\",\"Quality\":256,\"Filter\":\"PASS\",\"AllelicDepth\":203,\"Coverage\":25}"
		// Compare JSON strings to see if they are the same
		if expectedString != string(variantJSON) {
			t.Fatalf("Expected JSON string does not equal actual JSON string...")
		}
	}
}
