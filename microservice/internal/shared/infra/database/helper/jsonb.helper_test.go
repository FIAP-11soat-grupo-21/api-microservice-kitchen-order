package database_helper

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func verifyJSONBValue(t *testing.T, value interface{}, expectedData interface{}) {
	assert.NotNil(t, value)
	var result interface{}
	err := json.Unmarshal(value.([]byte), &result)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)
}

func TestJSONB_Value_Success(t *testing.T) {
	data := map[string]interface{}{"name": "test", "age": 30}
	jsonb := JSONB{Data: data}

	value, err := jsonb.Value()

	assert.NoError(t, err)
	var result map[string]interface{}
	err = json.Unmarshal(value.([]byte), &result)
	assert.NoError(t, err)
	assert.Equal(t, "test", result["name"])
	assert.Equal(t, float64(30), result["age"])
}

func TestJSONB_Value_WithString(t *testing.T) {
	jsonb := JSONB{Data: "simple string"}

	value, err := jsonb.Value()

	assert.NoError(t, err)
	var result string
	err = json.Unmarshal(value.([]byte), &result)
	assert.NoError(t, err)
	assert.Equal(t, "simple string", result)
}

func TestJSONB_Value_WithArray(t *testing.T) {
	data := []interface{}{"item1", "item2", 123}
	jsonb := JSONB{Data: data}

	value, err := jsonb.Value()

	assert.NoError(t, err)
	var result []interface{}
	err = json.Unmarshal(value.([]byte), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "item1", result[0])
	assert.Equal(t, "item2", result[1])
	assert.Equal(t, float64(123), result[2])
}

func TestJSONB_Value_WithNil(t *testing.T) {
	jsonb := JSONB{Data: nil}

	value, err := jsonb.Value()

	assert.NoError(t, err)
	assert.Equal(t, []byte("null"), value)
}

func TestJSONB_Scan_Success(t *testing.T) {
	data := map[string]interface{}{"name": "test", "age": 30}
	jsonData, _ := json.Marshal(data)
	var jsonb JSONB

	err := jsonb.Scan(jsonData)

	assert.NoError(t, err)
	dataMap, ok := jsonb.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "test", dataMap["name"])
	assert.Equal(t, float64(30), dataMap["age"])
}

func TestJSONB_Scan_WithString(t *testing.T) {
	jsonData := []byte(`"simple string"`)
	var jsonb JSONB

	err := jsonb.Scan(jsonData)

	assert.NoError(t, err)
	assert.Equal(t, "simple string", jsonb.Data)
}

func TestJSONB_Scan_WithArray(t *testing.T) {
	jsonData := []byte(`["item1", "item2", 123]`)
	var jsonb JSONB

	err := jsonb.Scan(jsonData)

	assert.NoError(t, err)
	dataArray, ok := jsonb.Data.([]interface{})
	assert.True(t, ok)
	assert.Len(t, dataArray, 3)
	assert.Equal(t, "item1", dataArray[0])
	assert.Equal(t, "item2", dataArray[1])
	assert.Equal(t, float64(123), dataArray[2])
}

func TestJSONB_Scan_InvalidType(t *testing.T) {
	var jsonb JSONB

	err := jsonb.Scan("not a byte slice")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type assertion to []byte failed")
}

func TestJSONB_Scan_InvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{"invalid": json}`)
	var jsonb JSONB

	err := jsonb.Scan(invalidJSON)

	assert.Error(t, err)
}

func TestJSONB_RoundTrip(t *testing.T) {
	originalData := map[string]interface{}{
		"name":    "test",
		"age":     30,
		"active":  true,
		"scores":  []interface{}{10, 20, 30},
		"details": map[string]interface{}{"city": "São Paulo"},
	}
	jsonb1 := JSONB{Data: originalData}

	value, err := jsonb1.Value()
	assert.NoError(t, err)

	var jsonb2 JSONB
	err = jsonb2.Scan(value)
	assert.NoError(t, err)

	dataMap, ok := jsonb2.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "test", dataMap["name"])
	assert.Equal(t, float64(30), dataMap["age"])
	assert.Equal(t, true, dataMap["active"])
}

func TestJSONBArray_Value_Success(t *testing.T) {
	array := JSONBArray{"item1", "item2", 123}

	value, err := array.Value()

	assert.NoError(t, err)
	var result []interface{}
	err = json.Unmarshal(value.([]byte), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "item1", result[0])
	assert.Equal(t, "item2", result[1])
	assert.Equal(t, float64(123), result[2])
}

func TestJSONBArray_Value_Empty(t *testing.T) {
	array := JSONBArray{}

	value, err := array.Value()

	assert.NoError(t, err)
	assert.Equal(t, []byte("[]"), value)
}

func TestJSONBArray_Scan_Success(t *testing.T) {
	jsonData := []byte(`["item1", "item2", 123]`)
	var array JSONBArray

	err := array.Scan(jsonData)

	assert.NoError(t, err)
	assert.Len(t, array, 3)
	assert.Equal(t, "item1", array[0])
	assert.Equal(t, "item2", array[1])
	assert.Equal(t, float64(123), array[2])
}

func TestJSONBArray_Scan_EmptyArray(t *testing.T) {
	jsonData := []byte(`[]`)
	var array JSONBArray

	err := array.Scan(jsonData)

	assert.NoError(t, err)
	assert.Len(t, array, 0)
}

func TestJSONBArray_Scan_InvalidType(t *testing.T) {
	var array JSONBArray

	err := array.Scan("not a byte slice")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type assertion to []byte failed")
}

func TestJSONBArray_Scan_InvalidJSON(t *testing.T) {
	invalidJSON := []byte(`[invalid json]`)
	var array JSONBArray

	err := array.Scan(invalidJSON)

	assert.Error(t, err)
}

func TestJSONBArray_RoundTrip(t *testing.T) {
	originalArray := JSONBArray{"item1", 123, true, map[string]interface{}{"key": "value"}}

	value, err := originalArray.Value()
	assert.NoError(t, err)

	var newArray JSONBArray
	err = newArray.Scan(value)
	assert.NoError(t, err)

	assert.Len(t, newArray, 4)
	assert.Equal(t, "item1", newArray[0])
	assert.Equal(t, float64(123), newArray[1])
	assert.Equal(t, true, newArray[2])

	lastItem, ok := newArray[3].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "value", lastItem["key"])
}
