package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

const TEST_LIMIT = 10

var testDataNormal = getTestDataFromJson()

var testDataError = []struct {
	data string
	want error
}{
	{
		data: "54422E0968050015070F110E11010115070F110E100100040000E4CCE202000108010B070B010000EB2343B83E9BFB0000",
		want: errInvalidPrefix,
	},
	{
		data: "5440200968050015070F110E11010115070F110E100100040000E4CCE202000108010B070B010000EB2343B83E9BFB0000",
		want: errInvalidSize,
	},
	{
		data: "5440CA0968050015070F110E16010215070F110E160100040000E4D1E202000108010A070B010080EB2343103F9BFB000000020000F301FFFF0400000030070601000700F3FF62002000070004002900B4FF6200210055007D0101000100010100000100170000030000000000A5140000182000010000000000B414000017200000000001000000000000000000000000000000000000000000000000000000000000000000000000000000680374047401600154014801D4014801CC0108024401000004FF0000000000000020FA",
		want: errInvalidSize,
	},
	{
		data: "54402E0968050015070F110F0C010115070F110F0C0100040000E503E3020001080109080C000080EB2343D83F9BFB00",
		want: io.ErrUnexpectedEOF,
	},
	{
		data: "54402E0968050015070F110F0C010115070F110F0C0100040000E503E3020001080109080C000080EB2343D83F9B",
		want: io.ErrUnexpectedEOF,
	},
	{
		data: "54402E0968050015070F110F0C010115070F110F0C0100040000E503E3",
		want: io.ErrUnexpectedEOF,
	},
}

func getTestDataFromJson() []string {
	jsonFile, err := os.Open("report_test_data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	var testData []string
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(byteValue, &testData); err != nil {
		log.Fatal(err)
	}
	return testData
}

func TestReport(t *testing.T) {
	type args struct {
		b []byte
	}
	type tester struct {
		name string
		args args
		want string
	}

	tests := make([]tester, len(testDataNormal))
	for i, d := range testDataNormal {
		tests[i].name = "data #" + strconv.Itoa(i)
		tests[i].args.b = hexToByte(d)
		tests[i].want = d
	}

	for i, tt := range tests {
		if tt.name == "" {
			continue
		}

		if i > TEST_LIMIT {
			break
		}

		t.Run(tt.name, func(t *testing.T) {
			if got, err := decodeReport(tt.args.b); err != nil {
				t.Errorf("error : %s", err)
			} else {
				// Validator
				if !got.ValidPrefix() {
					t.Errorf("Prefix is Not Valid")
				}

				if !got.ValidSize() {
					errString := fmt.Sprintf("Size is Not Valid. Got %d. want %d", got.Header.Size, got.Size())
					t.Errorf(errString)
				}

				// Encode Test
				// many case that can't be handled. cz float factorial
				// ex : 8.6 / 0.1 != 86.0
				encRes, err := encode(got)
				if err != nil {
					t.Errorf("encode error")
				}

				got2, _ := decodeReport(encRes)
				score := compareVar(got, got2)

				if score != 100 {
					errString := fmt.Sprintf("Not match. Score %d", score)
					t.Errorf(errString)
				}
			}
		})
	}
}

func TestReportErrorHandler(t *testing.T) {
	type args struct {
		b []byte
	}
	type tester struct {
		name string
		args args
		want error
	}

	tests := make([]tester, len(testDataError))
	for i, d := range testDataError {
		tests[i].name = "data #" + strconv.Itoa(i)
		tests[i].args.b = hexToByte(d.data)
		tests[i].want = d.want
	}

	for _, tt := range tests {
		if tt.name == "" {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			_, err := decodeReport(tt.args.b)
			if err == nil {
				errString := fmt.Sprintf("Success?? Packet should be Error (%s)", tt.want)
				t.Errorf(errString)
			} else if err != tt.want {
				errString := fmt.Sprintf("Gor error: %s. Packet should be Error (%s)", err, tt.want)
				t.Errorf(errString)
			}
		})
	}
}

// compare between 2 of any variabel
func compareVar(v1 interface{}, v2 interface{}) (score int) {
	rv1 := reflect.ValueOf(v1)
	rv2 := reflect.ValueOf(v2)

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
