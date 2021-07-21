package report

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

func Test_report(t *testing.T) {
	type args struct {
		b []byte
	}
	type tester struct {
		name string
		args args
		want string
	}

	var testData []string
	if err := openFileJSON("report_test_data.json", &testData); err != nil {
		log.Fatal(err)
	}
	// jsonFile, err := os.Open("report_test_data.json")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// byteValue, _ := ioutil.ReadAll(jsonFile)
	// json.Unmarshal(byteValue, &testData)
	// jsonFile.Close()

	// can it be simpler?
	// it was. cz i paste from hex lib example
	tests := make([]tester, len(testData))
	for i, d := range testData {
		tests[i].name = "data #" + strconv.Itoa(i)
		tests[i].args.b = hexToBytes(d)
		tests[i].want = d
	}

	for _, tt := range tests {
		if tt.name == "" {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			// fmt.Printf("=== [%s] ===\n", tt.name)
			rr := New(tt.args.b)
			rr.reader.Len()
			if got, err := rr.Decode(); err != nil {
				t.Errorf("got = %v, want %v", &got, tt.want)
			} else {
				if rr.reader.Len() != 0 {
					t.Errorf("some buffer not read")
				}

				// many case that can't be handled. cz float factorial
				// ex : 8.6 / 0.1 != 06.0
				encRes, err := shared.Encode(got)
				if err != nil {
					t.Errorf("encode error")
				}

				hexRes := bytesToHex(encRes)
				isMatch := true
				notMatchIdx := 0
				for i, v := range hexRes {
					if rune(tt.want[i]) != v {
						notMatchIdx = i
						isMatch = false
					}
				}

				if !isMatch {
					t.Errorf("encode not match in index" + strconv.Itoa(notMatchIdx))
					// fmt.Println("============", notMatchIdx/2)
					// fmt.Println(tt.want)
					// fmt.Println(hexRes[:notMatchIdx+1])
					// fmt.Println(tt.args.b)
					// fmt.Println(encRes)
					// fmt.Println(got)
				}
				// if got.Mcu != nil && !got.Mcu.Active {
				// 	fmt.Printf("=== [%s] ===\n", tt.name)
				// 	fmt.Println(got)
				// }
			}
		})
	}
}

func openFileJSON(filename string, testData *[]string) error {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, testData)
	return nil
}

// it's to convert src:hexstring to dst:bytes
func hexToBytes(hexString string) []byte {
	src := []byte(strings.ToLower(hexString))
	dst := make([]byte, hex.DecodedLen(len(src)))
	n, err := hex.Decode(dst, src)
	if err != nil {
		log.Fatal(err)
	}
	return dst[:n]
}

// it's to convert dst:bytes to src:hexstring
func bytesToHex(b []byte) string {
	dst := make([]byte, hex.EncodedLen(len(b)))
	n := hex.Encode(dst, b)
	hexString := string(dst[:n])
	return strings.ToUpper(hexString)
}
