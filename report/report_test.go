package report

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
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

	// it was. cz i paste from hex lib example
	tests := make([]tester, len(testData))
	for i, d := range testData {
		tests[i].name = "data #" + strconv.Itoa(i)
		tests[i].args.b = util.Hex2Byte(d)
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
				// ex : 8.6 / 0.1 != 86.0
				encRes, err := shared.Encode(got)
				if err != nil {
					t.Errorf("encode error")
				}

				hexRes := util.Byte2Hex(encRes)
				numMatch := 0
				notMatchIdx := 0
				for i, v := range hexRes {
					if rune(tt.want[i]) != v {
						notMatchIdx = i
					} else {
						numMatch++
					}
				}

				score := (numMatch * 100) / len(hexRes)

				if score < 70 {
					errString := fmt.Sprintf("not match in #%d. match (%d of %d) Score %d", notMatchIdx, numMatch, len(hexRes), score)
					t.Errorf(errString)
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

// openFileJSON open and decode json file
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
