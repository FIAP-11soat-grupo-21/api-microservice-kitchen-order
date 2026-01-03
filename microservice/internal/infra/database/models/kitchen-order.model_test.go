package models

import (
	"testing"
)

func TestKitchenOrderModel_TableName(t *testing.T) {
	// Arrange
	model := KitchenOrderModel{}

	// Act
	tableName := model.TableName()

	// Assert
	expectedTableName := "kitchen_order"
	if tableName != expectedTableName {
		t.Errorf("Expected table name '%s', got '%s'", expectedTableName, tableName)
	}
}