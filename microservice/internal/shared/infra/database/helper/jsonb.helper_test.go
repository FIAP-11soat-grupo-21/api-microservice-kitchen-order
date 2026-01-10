package database_helper

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONB_Value_Success(t *testing.T) {
	// Arrange
	data := map[string]interface{}{
		"name": "test",
		"age":  30,
	}
	jsonb := JSONB{Data: data}

	// Act
	value, err := jsonb.Value()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, value)
	
	// Verifica se o valor pode ser deserializado de volta
	var result map[string]interface{}
	err = json.Unmarshal(value.([]byte), &result)
	assert.NoError(t, err)
	assert.Equal(t, "test", result["name"])
	assert.Equal(t, float64(30), result["age"]) // JSON unmarshals numbers as float64
}

func TestJSONB_Value_WithString(t *testing.T) {
	// Arrange
	jsonb := JSONB{Data: "simple string"}

	// Act
	value, err := jsonb.Value()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, value)
	
	// Verifica se o valor é uma string JSON válida
	var result string
	err = json.Unmarshal(value.([]byte), &result)
	assert.NoError(t, err)
	assert.Equal(t, "simple string", result)
}

func TestJSONB_Value_WithArray(t *testing.T) {
	// Arrange
	data := []interface{}{"item1", "item2", 123}
	jsonb := JSONB{Data: data}

	// Act
	value, err := jsonb.Value()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, value)
	
	// Verifica se o valor pode ser deserializado de volta
	var result []interface{}
	err = json.Unmarshal(value.([]byte), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "item1", result[0])
	assert.Equal(t, "item2", result[1])
	assert.Equal(t, float64(123), result[2])
}

func TestJSONB_Value_WithNil(t *testing.T) {
	// Arrange
	jsonb := JSONB{Data: nil}

	// Act
	value, err := jsonb.Value()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []byte("null"), value)
}

func TestJSONB_Scan_Success(t *testing.T) {
	// Arrange
	data := map[string]interface{}{
		"name": "test",
		"age":  30,
	}
	jsonData, _ := json.Marshal(data)
	
	var jsonb JSONB

	// Act
	err := jsonb.Scan(jsonData)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, jsonb.Data)
	
	// Verifica se os dados foram deserializados corretamente
	dataMap, ok := jsonb.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "test", dataMap["name"])
	assert.Equal(t, float64(30), dataMap["age"])
}

func TestJSONB_Scan_WithString(t *testing.T) {
	// Arrange
	jsonData := []byte(`"simple string"`)
	
	var jsonb JSONB

	// Act
	err := jsonb.Scan(jsonData)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "simple string", jsonb.Data)
}

func TestJSONB_Scan_WithArray(t *testing.T) {
	// Arrange
	jsonData := []byte(`["item1", "item2", 123]`)
	
	var jsonb JSONB

	// Act
	err := jsonb.Scan(jsonData)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, jsonb.Data)
	
	// Verifica se os dados foram deserializados corretamente
	dataArray, ok := jsonb.Data.([]interface{})
	assert.True(t, ok)
	assert.Len(t, dataArray, 3)
	assert.Equal(t, "item1", dataArray[0])
	assert.Equal(t, "item2", dataArray[1])
	assert.Equal(t, float64(123), dataArray[2])
}

func TestJSONB_Scan_InvalidType(t *testing.T) {
	// Arrange
	var jsonb JSONB

	// Act
	err := jsonb.Scan("not a byte slice")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type assertion to []byte failed")
}

func TestJSONB_Scan_InvalidJSON(t *testing.T) {
	// Arrange
	invalidJSON := []byte(`{"invalid": json}`)
	
	var jsonb JSONB

	// Act
	err := jsonb.Scan(invalidJSON)

	// Assert
	assert.Error(t, err)
}

func TestJSONBArray_Value_Success(t *testing.T) {
	// Arrange
	array := JSONBArray{"item1", "item2", 123}

	// Act
	value, err := array.Value()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, value)
	
	// Verifica se o valor pode ser deserializado de volta
	var result []interface{}
	err = json.Unmarshal(value.([]byte), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "item1", result[0])
	assert.Equal(t, "item2", result[1])
	assert.Equal(t, float64(123), result[2])
}

func TestJSONBArray_Value_Empty(t *testing.T) {
	// Arrange
	array := JSONBArray{}

	// Act
	value, err := array.Value()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []byte("[]"), value)
}

func TestJSONBArray_Scan_Success(t *testing.T) {
	// Arrange
	jsonData := []byte(`["item1", "item2", 123]`)
	
	var array JSONBArray

	// Act
	err := array.Scan(jsonData)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, array, 3)
	assert.Equal(t, "item1", array[0])
	assert.Equal(t, "item2", array[1])
	assert.Equal(t, float64(123), array[2])
}

func TestJSONBArray_Scan_EmptyArray(t *testing.T) {
	// Arrange
	jsonData := []byte(`[]`)
	
	var array JSONBArray

	// Act
	err := array.Scan(jsonData)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, array, 0)
}

func TestJSONBArray_Scan_InvalidType(t *testing.T) {
	// Arrange
	var array JSONBArray

	// Act
	err := array.Scan("not a byte slice")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type assertion to []byte failed")
}

func TestJSONBArray_Scan_InvalidJSON(t *testing.T) {
	// Arrange
	invalidJSON := []byte(`[invalid json]`)
	
	var array JSONBArray

	// Act
	err := array.Scan(invalidJSON)

	// Assert
	assert.Error(t, err)
}

func TestJSONB_RoundTrip(t *testing.T) {
	// Testa o ciclo completo: Value -> Scan
	
	// Arrange
	originalData := map[string]interface{}{
		"name":    "test",
		"age":     30,
		"active":  true,
		"scores":  []interface{}{10, 20, 30},
		"details": map[string]interface{}{"city": "São Paulo"},
	}
	
	jsonb1 := JSONB{Data: originalData}

	// Act - Value
	value, err := jsonb1.Value()
	assert.NoError(t, err)

	// Act - Scan
	var jsonb2 JSONB
	err = jsonb2.Scan(value)
	assert.NoError(t, err)

	// Assert
	dataMap, ok := jsonb2.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "test", dataMap["name"])
	assert.Equal(t, float64(30), dataMap["age"])
	assert.Equal(t, true, dataMap["active"])
}

func TestJSONBArray_RoundTrip(t *testing.T) {
	// Testa o ciclo completo: Value -> Scan
	
	// Arrange
	originalArray := JSONBArray{"item1", 123, true, map[string]interface{}{"key": "value"}}

	// Act - Value
	value, err := originalArray.Value()
	assert.NoError(t, err)

	// Act - Scan
	var newArray JSONBArray
	err = newArray.Scan(value)
	assert.NoError(t, err)

	// Assert
	assert.Len(t, newArray, 4)
	assert.Equal(t, "item1", newArray[0])
	assert.Equal(t, float64(123), newArray[1])
	assert.Equal(t, true, newArray[2])
	
	lastItem, ok := newArray[3].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "value", lastItem["key"])
}