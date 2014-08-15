// harness.go
package resttest

import (
	"bytes"
	//"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mws/util"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

/* Internal */
var (
	count util.TestCounter
)

/*
 Deprecated infavor of GetCallerName()
*/
func GetFunctionName(i interface{}) string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

/*
  Comprises a RESTful test and expected results
*/
type RESTTestContainer struct {
	Desc       string
	Handler    func(http.ResponseWriter, *http.Request)
	Path       string
	Sub        map[string]string
	Method     string
	JSONParams string
	Params     url.Values
	Status     int
	Match      map[string]bool
	//add match interface and create interface and match the list of properties in the
	// match map against the returned object
	MatchVal        interface{}
	MatchFields     []string
	PreAuthId       int64
	OwnThreshold    int64
	GlobalThreshold int64
}

/*
Execute a collection of RESTful tests and match against the results, if specified
*/
func RunTestSet(t *testing.T, tests []RESTTestContainer) {
	log.Printf("-----Starting REST TestSet ( %s )---------------------", util.GetCallerName())
	//
	// Execute the tests
	for _, test := range tests {
		log.Printf("----------Starting REST Test ( %s )---------------------", test.Desc)
		//
		//Substitue path values
		failedTest := false
		for k, subStr := range test.Sub {
			//run replace on the string for each item in submap
			re, error := regexp.Compile(k)
			if error != nil {
				t.Errorf("%s: Substitute = %s, want %s", test.Path, k, subStr)
				failedTest = true
			}
			test.Path = re.ReplaceAllString(test.Path, subStr)
		}
		log.Printf("---------------- Testing REST Path ( %s )---------------------", test.Path)
		//
		// Initialize Response Recorder
		record := httptest.NewRecorder()
		buf := []byte(test.JSONParams)
		body := bytes.NewBuffer(buf)
		req, _ := http.NewRequest(test.Method, test.Path, body)
		req.Header.Set("Content-Type", "application/json")
		//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		//req.Header.Set("Authorization", "Basic "+base64.URLEncoding.EncodeToString([]byte("admin:admin")))
		//
		// Set the pre auth header for all tests by default,
		// so this user must exist in the database
		if test.PreAuthId >= 0 {
			req.Header.Set("X-meritwiki-client-preauth-id", strconv.FormatInt(test.PreAuthId, 10))
		}
		if test.OwnThreshold >= 0 {
			req.Header.Set("X-meritwiki-owner-threshold", strconv.FormatInt(test.OwnThreshold, 10))
		}
		if test.GlobalThreshold >= 0 {
			req.Header.Set("X-meritwiki-global-threshold", strconv.FormatInt(test.GlobalThreshold, 10))
		}
		test.Handler(record, req)
		//
		// Call Handler
		bodyStr := string(record.Body.Bytes()[:])
		log.Printf("[%s] response body  =[%s]", test.Desc, bodyStr)
		//
		// Test the response code
		if got, want := record.Code, test.Status; got != want {
			t.Errorf("%s: response code = %d, want %d", test.Desc, got, want)
			failedTest = true
			goto final
		}
		log.Printf("%s: response code = %d, want %d", test.Desc, record.Code, test.Status)

		if test.MatchVal != nil {
			i := reflect.New(reflect.TypeOf(test.MatchVal)).Interface()
			buf1 := new(bytes.Buffer)
			io.Copy(buf1, record.Body)
			json.Unmarshal(buf1.Bytes(), i)
			//log.Printf("Unmarshalled %v", i)
			pass, results, err := comparePartial(i, test.MatchVal, test.MatchFields)
			if err != nil {
				t.Errorf("Compare error: ", err)
				failedTest = true
			}
			//log.Printf("Unmatched results=%d", len(results))

			if !pass && results != nil && len(results) > 0 {
				for _, r := range results {
					//result := results[0]
					log.Printf("Name: %v, Value1: %v, Value2: %v \n", r.FieldName, r.Value1, r.Value2)
					t.Errorf("Unmatching results returned")
					failedTest = true
				}
			} else {

			}
		} else {
			//
			// Test the response against expected results
			for re, _ := range test.Match {
				//re = strings.Replace(re, " ", "", -1)
				//re = strings.TrimSpace(re)
				re := stripchars(re, "\n")
				newbody := stripchars(record.Body.String(), "\n")
				//body = strings.Replace(body, " ", "", -1)
				//body := strings.Replace(record.Body.String(), "\r\n", "~", 0)
				//body = strings.TrimSpace(body)
				//if got, want := body, re; got != want {

				got, msg, err := MatchJSON(newbody, re)
				if !got || err != nil {
					t.Errorf("Problem Matching JSON %s, %s", msg, err)
					//t.Errorf("%s:\r\n Returned %s\r\n Expected %s", test.Desc, body, re)
					failedTest = true
				}
				log.Printf("%s:\r\n Returned %s\r\n Expected %s", test.Desc, newbody, re)
			}
		}
	final:
		if failedTest {
			count.FailCount++
			log.Printf("----------Finished REST Test ( %s ) FAILED --------------------", test.Desc)
		} else {
			count.SuccessCount++
			log.Printf("----------Finished REST Test ( %s ) PASSED --------------------", test.Desc)
		}
		log.Printf("-------------------------------------------------------------------------------")

	}
	log.Printf("-----Finished REST TestSet ( %s )---------------------", util.GetCallerName())
	log.Printf("Success Count = %d Fail Count %d", count.SuccessCount, count.FailCount)
	log.Printf("-------------------------------------------------------------------------------")
}

func stripchars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}

func MatchJSON(actual, want interface{}) (success bool, message string, err error) {
	actualString, aok := toString(actual)
	expectedString, eok := toString(want)

	if aok && eok {
		abuf := new(bytes.Buffer)
		ebuf := new(bytes.Buffer)

		if err := json.Indent(abuf, []byte(actualString), "", "  "); err != nil {
			return false, "", err
		}

		if err := json.Indent(ebuf, []byte(expectedString), "", "  "); err != nil {
			return false, "", err
		}

		var aval interface{}
		var eval interface{}

		json.Unmarshal([]byte(actualString), &aval)
		json.Unmarshal([]byte(expectedString), &eval)

		if reflect.DeepEqual(aval, eval) {
			return true, fmt.Sprintf("%s not to match JSON of %s", abuf.String(), ebuf.String()), nil
		} else {
			return false, fmt.Sprintf("%s  to match JSON of %s", abuf.String(), ebuf.String()), nil
		}
	} else {
		return false, "", fmt.Errorf("MatchJSONMatcher matcher requires a string or stringer.  Got:\n%#v", actual)
	}
	return false, "", nil
}

func toString(a interface{}) (string, bool) {
	aString, isString := a.(string)
	if isString {
		return aString, true
	}

	aBytes, isBytes := a.([]byte)
	if isBytes {
		return string(aBytes), true
	}

	aStringer, isStringer := a.(fmt.Stringer)
	if isStringer {
		return aStringer.String(), true
	}

	return "", false
}

func MatchJSON2Interface(wantJSON string, got interface{}) (err error) {

	gotJSON, err := json.Marshal(got)
	if err != nil {
		log.Printf("Failure to match JSON due to %s", err)
		return fmt.Errorf("Failure to match JSON due to %s", err)
	}
	//
	if wantJSON != string(gotJSON) {
		log.Printf("Wanted JSON = \r\n%s\r\n Got JSON:\r\n %s", wantJSON, gotJSON)
		return fmt.Errorf("Failure to match JSON due to %s", err)
	}
	log.Printf("Successfully matched  %s", string(gotJSON))
	return nil
}
