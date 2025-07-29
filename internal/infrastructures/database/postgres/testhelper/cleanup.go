package testhelper

import (
	"gorm.io/gorm"
	"testing"
)

func CleanupTable(t *testing.T, db *gorm.DB, tableName string) {
	if err := db.Exec("TRUNCATE TABLE " + tableName + " RESTART IDENTITY CASCADE").Error; err != nil {
		t.Fatalf("failed to truncate table %s: %v", tableName, err)
	}
}
