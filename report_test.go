package sdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
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
		tests[i].args.b = hexToByte(d)
		tests[i].want = d
	}

	for i, tt := range tests {
		if tt.name == "" {
			continue
		}

		// limit test
		if i > 10 {
			break
		}

		t.Run(tt.name, func(t *testing.T) {
			// fmt.Printf("=== [%s] ===\n", tt.name)
			rr := newReport(tt.args.b)
			if got, err := rr.decode(); err != nil {
				t.Errorf("got = %v, want %v", &got, tt.want)
			} else {
				if rr.reader.Len() != 0 {
					t.Errorf("some buffer not read")
				}

				// many case that can't be handled. cz float factorial
				// ex : 8.6 / 0.1 != 86.0
				encRes, err := encode(got)
				if err != nil {
					t.Errorf("encode error")
				}

				got2, _ := newReport(encRes).decode()
				score := compareVar(got, got2)

				if score != 100 {
					errString := fmt.Sprintf("Not match. Score %d", score)
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

// compare between 2 of any variabel
func compareVar(v1 interface{}, v2 interface{}) (score int) {

	rv1 := reflect.ValueOf(v1)
	rv2 := reflect.ValueOf(v2)

	// compare kind
	if rv1.Kind() != rv2.Kind() {
		return 0
	}

	score = 0
	switch rk := rv1.Kind(); rk {

	case reflect.Ptr:
		// if one is nil, both of rv1 and rv2 must be nil
		if rv1.IsNil() || rv2.IsNil() {
			if rv1.IsNil() && rv2.IsNil() {
				return 100
			}
			return 0
		}
		rv1 = rv1.Elem()
		rv2 = rv2.Elem()
		score = compareVar(rv1.Interface(), rv2.Interface())

	case reflect.Struct:
		if rv1.Type() != rv2.Type() {
			return 0
		}
		if rv1.Type() == typeOfTime {
			t1 := rv1.Interface().(time.Time)
			t2 := rv2.Interface().(time.Time)
			if t1.Unix() == t2.Unix() {
				score = 100
			}

		} else {
			totalScore := 0
			numFiled := rv1.NumField()
			for i := 0; i < numFiled; i++ {
				rvField1 := rv1.Field(i)
				rvField2 := rv2.Field(i)

				tmpScore := compareVar(rvField1.Interface(), rvField2.Interface())
				totalScore += tmpScore
				// fmt.Printf("%d(%s) ", tmpScore, rvField1.Type())
			}
			score = (totalScore) / numFiled
			// fmt.Printf(" = %d | avg = %d\n", totalScore, score)
		}

	case reflect.Array:
		if rv1.Len() != rv2.Len() {
			return 0
		}

		totalScore := 0
		arrLen := rv1.Len()
		for i := 0; i < arrLen; i++ {
			totalScore += compareVar(rv1.Index(i).Interface(), rv2.Index(i).Interface())
		}
		score = (totalScore) / arrLen

	case reflect.String:
		if rv1.String() == rv2.String() {
			score = 100
		}

	case reflect.Bool:
		if rv1.Bool() == rv2.Bool() {
			score = 100
		}

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		if rv1.Uint() == rv2.Uint() {
			score = 100
		}

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		if rv1.Int() == rv2.Int() {
			score = 100
		}

	case reflect.Float32, reflect.Float64:
		rvFloat1 := math.Round(rv1.Float()*100) / 100
		rvFloat2 := math.Round(rv2.Float()*100) / 100
		if rvFloat1 == rvFloat2 {
			score = 100
		}

	default:
		score = 100
	}
	return score
}
