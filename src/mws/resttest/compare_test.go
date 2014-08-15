package resttest

import (
	"testing"
	"time"
)

func TestCompareDiffString(t *testing.T) {
	type testStruct struct {
		Str string
	}

	struct1 := testStruct{"A"}
	struct2 := testStruct{"BB"}
	_, results, err := compare(struct1, struct2)
	if err != nil {
		t.Errorf("Compare error: ", err)
	}
	if results != nil && len(results) == 1 {
		result := results[0]
		if result.FieldName != "Str" || result.Value1 != "A" || result.Value2 != "BB" {
			t.Errorf("Name: %v, Value1: %v, Value2: %v \n", result.FieldName, result.Value1, result.Value2)
		}
	} else {
		t.Errorf("Incorrect number of results returned")
	}
}

func TestCompareDiffInt(t *testing.T) {
	type testStruct struct {
		Integer int
	}

	struct1 := testStruct{1}
	struct2 := testStruct{2}
	_, results, err := compare(struct1, struct2)
	if err != nil {
		t.Errorf("Compare error: ", err)
	}
	if results != nil && len(results) == 1 {
		result := results[0]
		if result.FieldName != "Integer" || result.Value1 != 1 || result.Value2 != 2 {
			t.Errorf("Name: %v, Value1: %v, Value2: %v \n", result.FieldName, result.Value1, result.Value2)
		}
	} else {
		t.Errorf("Incorrect number of results returned")
	}
}

func TestCompareDiffFloat(t *testing.T) {
	type testStruct struct {
		F float64
	}

	struct1 := testStruct{1.01}
	struct2 := testStruct{2.02}
	_, results, err := compare(struct1, struct2)
	if err != nil {
		t.Errorf("Compare error: ", err)
	}
	if results != nil && len(results) == 1 {
		result := results[0]
		if result.FieldName != "F" || result.Value1 != 1.01 || result.Value2 != 2.02 {
			t.Errorf("Name: %v, Value1: %v, Value2: %v \n", result.FieldName, result.Value1, result.Value2)
		}
	} else {
		t.Errorf("Incorrect number of results returned")
	}
}

func TestCompareDiffBool(t *testing.T) {
	type testStruct struct {
		B bool
	}

	struct1 := testStruct{true}
	struct2 := testStruct{false}
	_, results, err := compare(struct1, struct2)
	if err != nil {
		t.Errorf("compare error: ", err)
	}
	if results != nil && len(results) == 1 {
		result := results[0]
		if result.FieldName != "B" || result.Value1 != true || result.Value2 != false {
			t.Errorf("Name: %v, Value1: %v, Value2: %v \n", result.FieldName, result.Value1, result.Value2)
		}
	} else {
		t.Errorf("Incorrect number of results returned")
	}
}

func TestCompareDiffTime(t *testing.T) {
	type testStruct struct {
		Ti time.Time
	}

	struct1 := testStruct{time.Now()}
	struct2 := testStruct{time.Now().AddDate(0, 0, 1)}
	_, results, err := compare(struct1, struct2)
	if err != nil {
		t.Errorf("Compare error: ", err)
	}
	if results != nil && len(results) == 1 {
		result := results[0]
		if result.FieldName != "Ti" {
			t.Errorf("Name: %v, Value1: %v, Value2: %v \n", result.FieldName, result.Value1, result.Value2)
		}
	} else {
		t.Errorf("Incorrect number of results returned")
	}
}

func TestCompareSameString(t *testing.T) {
	type testStruct struct {
		Str string
	}

	struct1 := testStruct{"Luke"}
	struct2 := testStruct{"Luke"}
	_, results, err := compare(struct1, struct2)
	if err != nil {
		t.Errorf("Compare error: ", err)
	}
	if results == nil || len(results) != 0 {
		t.Errorf("Incorrect number of results returned")
	}
}

func TestCompareSameInt(t *testing.T) {
	type testStruct struct {
		Integer int
	}

	struct1 := testStruct{1}
	struct2 := testStruct{1}
	_, results, err := compare(struct1, struct2)
	if err != nil {
		t.Errorf("Compare error: ", err)
	}
	if results == nil || len(results) != 0 {
		t.Errorf("Incorrect number of results returned")
	}
}

func TestCompareSameBool(t *testing.T) {
	type testStruct struct {
		b bool
	}

	struct1 := testStruct{true}
	struct2 := testStruct{true}
	_, results, err := compare(struct1, struct2)
	if err != nil {
		t.Errorf("Compare error: ", err)
	}
	if results == nil || len(results) != 0 {
		t.Errorf("Incorrect number of results returned")
	}
}

func TestCompareMultiple(t *testing.T) {
	type testStruct struct {
		Str     string
		Integer int
		F       float64
		B       bool
	}

	struct1 := testStruct{"a", 1, 1.01, true}
	struct2 := testStruct{"b", 2, 2.02, false}
	_, results, err := compare(struct1, struct2)
	if err != nil {
		t.Errorf("Compare error: ", err)
	}
	if results == nil || len(results) != 4 {
		t.Errorf("Incorrect number of results returned")
	}
}

func TestCompareDiffStringPtr(t *testing.T) {
	type testStruct struct {
		Str *string
	}

	struct1 := testStruct{SPtr("A")}
	struct2 := testStruct{SPtr("BB")}
	_, results, err := compare(struct1, struct2)
	if err != nil {
		t.Errorf("Compare error: ", err)
	}
	if results != nil && len(results) == 1 {
		result := results[0]
		if result.FieldName != "Str" || result.Value1 != "A" || result.Value2 != "BB" {
			t.Errorf("Name: %v, Value1: %v, Value2: %v \n", result.FieldName, result.Value1, result.Value2)
		}
	} else {
		t.Errorf("Incorrect number of results returned")
	}
}

func TestCompareDiffStringNilPtr(t *testing.T) {
	type testStruct struct {
		Str *string
	}

	struct1 := testStruct{nil}
	struct2 := testStruct{SPtr("BB")}
	_, results, err := compare(struct1, struct2)
	if err != nil {
		t.Errorf("Compare error: ", err)
	}
	if results != nil && len(results) == 1 {
		result := results[0]
		if result.FieldName != "Str" || result.Value1 != "" || result.Value2 != "BB" {
			t.Errorf("Name: %v, Value1: %v, Value2: %v \n", result.FieldName, result.Value1, result.Value2)
		}
	} else {
		t.Errorf("Incorrect number of results returned")
	}
}

func TestCompareSameStringPtr(t *testing.T) {
	type testStruct struct {
		Str *string
	}

	struct1 := testStruct{SPtr("AA")}
	struct2 := testStruct{SPtr("AA")}
	_, results, err := compare(struct1, struct2)
	if err != nil {
		t.Errorf("Compare error: ", err)
	}
	if results == nil || len(results) != 0 {
		t.Errorf("Incorrect number of results returned")
	}
}

func TestCompareDiffIntPtr(t *testing.T) {
	type testStruct struct {
		In *int
	}

	i1 := 1
	i2 := 2
	struct1 := testStruct{&i1}
	struct2 := testStruct{&i2}
	_, results, err := compare(struct1, struct2)
	if err != nil {
		t.Errorf("Compare error: ", err)
	}
	if results != nil && len(results) == 1 {
		result := results[0]
		if result.FieldName != "In" || result.Value1 != 1 || result.Value2 != 2 {
			t.Errorf("Name: %v, Value1: %v, Value2: %v \n", result.FieldName, result.Value1, result.Value2)
		}
	} else {
		t.Errorf("Incorrect number of results returned")
	}
}

func TestCompareDiffIntNilPtr(t *testing.T) {
	type testStruct struct {
		In *int
	}

	var i1 *int
	i1 = nil
	i2 := 2
	struct1 := testStruct{i1}
	struct2 := testStruct{&i2}
	_, results, err := compare(struct1, struct2)
	if err != nil {
		t.Errorf("Compare error: ", err)
	}
	if results != nil && len(results) == 1 {
		result := results[0]
		if result.FieldName != "In" || result.Value1 != 0 || result.Value2 != 2 {
			t.Errorf("Name: %v, Value1: %v, Value2: %v \n", result.FieldName, result.Value1, result.Value2)
		}
	} else {
		t.Errorf("Incorrect number of results returned")
	}
}

func TestCompareSameIntPtr(t *testing.T) {
	type testStruct struct {
		In *int
	}

	i1 := 1
	i2 := 1
	struct1 := testStruct{&i1}
	struct2 := testStruct{&i2}
	_, results, err := compare(struct1, struct2)
	if err != nil {
		t.Errorf("Compare error: ", err)
	}
	if results == nil || len(results) != 0 {
		t.Errorf("Incorrect number of results returned: %v", len(results))
	}
}

func TestCompareDiffBoolPtr(t *testing.T) {
	type testStruct struct {
		B *bool
	}

	b1 := true
	b2 := false
	struct1 := testStruct{&b1}
	struct2 := testStruct{&b2}
	_, results, err := compare(struct1, struct2)
	if err != nil {
		t.Errorf("Compare error: ", err)
	}
	if results != nil && len(results) == 1 {
		result := results[0]
		if result.FieldName != "B" || result.Value1 != true || result.Value2 != false {
			t.Errorf("Name: %v, Value1: %v, Value2: %v \n", result.FieldName, result.Value1, result.Value2)
		}
	} else {
		t.Errorf("Incorrect number of results returned")
	}
}

func TestCompareSameBoolPtr(t *testing.T) {
	type testStruct struct {
		B *bool
	}

	b1 := true
	b2 := true
	struct1 := testStruct{&b1}
	struct2 := testStruct{&b2}
	_, results, err := compare(struct1, struct2)
	if err != nil {
		t.Errorf("Compare error: ", err)
	}
	if results == nil || len(results) != 0 {
		t.Errorf("Incorrect number of results returned")
	}
}

func TestCompareDiffTimePtr(t *testing.T) {
	type testStruct struct {
		Ti *time.Time
	}

	time1 := time.Now()
	time2 := time.Now().AddDate(0, 0, 1)

	struct1 := testStruct{&time1}
	struct2 := testStruct{&time2}
	_, results, err := compare(struct1, struct2)
	if err != nil {
		t.Errorf("Compare error: ", err)
	}
	if results != nil && len(results) == 1 {
		result := results[0]
		if result.FieldName != "Ti" {
			t.Errorf("Name: %v, Value1: %v, Value2: %v \n", result.FieldName, result.Value1, result.Value2)
		}
	} else {
		t.Errorf("Incorrect number of results returned")
	}
}

func TestComparePtrStruct(t *testing.T) {
	type testStruct struct {
		Str string
	}

	struct1 := &testStruct{"A"}
	struct2 := &testStruct{"BB"}
	_, results, err := compare(struct1, struct2)
	if err != nil {
		t.Errorf("Compare error: ", err)
	}
	if results != nil && len(results) == 1 {
		result := results[0]
		if result.FieldName != "Str" || result.Value1 != "A" || result.Value2 != "BB" {
			t.Errorf("Name: %v, Value1: %v, Value2: %v \n", result.FieldName, result.Value1, result.Value2)
		}
	} else {
		t.Errorf("Incorrect number of results returned: %v", len(results))
	}
}

func BenchmarkCompareMultiple(b *testing.B) {
	b.StopTimer()

	type testStruct struct {
		Str     string
		Integer int
		F       float64
		B       bool
	}

	struct1 := &testStruct{"a", 1, 1.01, true}
	struct2 := &testStruct{"b", 2, 2.02, false}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = compare(struct1, struct2)
	}
}

func SPtr(str string) *string {
	return &str

}
