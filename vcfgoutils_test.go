package vcfgoutils

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/brentp/vcfgo"
	"github.com/stretchr/testify/assert"
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
		expectedJSON := []string{
			"{\"Chromosome\":\"chr4\"",
			"\"Position\":1806181",
			"\"Reference\":\"C\"",
			"\"Alternate\":\"T\"",
			"\"Quality\":256",
			"\"Filter\":\"PASS\"",
			"\"AllelicDepth\":203",
			"\"Coverage\":25}"}
		expectedJSONString := strings.Join(expectedJSON, ",")
		assert.Equal(
			t,
			string(variantJSON),
			expectedJSONString,
			"The actual and expected JSON strings should be the same.")
	}
}
