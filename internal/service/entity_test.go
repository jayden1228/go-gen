package service

import (
	"go-gen/internal/pkg/database"
	"os"
	"testing"
)

func TestCreateDBEntity(t *testing.T) {
	// format
	format := []string{"json", "gorm"}
	// path
	oPath, _ := os.Getwd()
	// module name
	projectModule := "go-gen"
	database.SetUp()
	CreateDBEntity(format, oPath, projectModule)
}
