package report

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
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
	jsonFile, err := os.Open("report_test_data.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &testData)
	jsonFile.Close()

	tests := make([]tester, len(testData))
	for i, d := range testData {
		src := []byte(strings.ToLower(d))
		dst := make([]byte, hex.DecodedLen(len(src)))
		n, err := hex.Decode(dst, src)
		if err != nil {
			log.Fatal(err)
		}
		tests[i].name = "Data #" + strconv.Itoa(i)
		tests[i].args.b = dst[:n]
	}

	for _, tt := range tests {
		if tt.name == "" {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("=== [%s] ===\n", tt.name)
			rr := New(tt.args.b)
			rr.reader.Len()
			if got, err := rr.Decode(); err != nil {
				t.Errorf("serve() = %v, want %v", &got, tt.want)
			} else {
				if rr.reader.Len() != 0 {
					t.Errorf("Some buffer not read")
				}
				// if got.Mcu != nil && !got.Mcu.Active {
				// 	fmt.Printf("=== [%s] ===\n", tt.name)
				// 	fmt.Println(got)
				// }
			}
		})
	}
}
