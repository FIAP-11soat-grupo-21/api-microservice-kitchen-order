package models

import (
	"testing"
)

func TestOrderStatusModel_TableName(t *testing.T) {
	// Arrange
	model := OrderStatusModel{}

	// Act
	tableName := model.TableName()

	// Assert
	expectedTableName := "order_status"
	if tableName != expectedTableName {
		t.Errorf("Expected table name '%s', got '%s'", expectedTableName, tableName)
	}
}
