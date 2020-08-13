package service

import (
	"go-gen/internal/pkg/database"
	"testing"
)

func TestCreateDBEntity(t *testing.T) {
	format := []string{"json", "gorm"}
	database.SetUp()
	CreateDBEntity(format)
}
