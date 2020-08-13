package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateDBEntity(t *testing.T) {

	e := EntityReq{
		EntityPkg: "damn",
	}

	assert.Nil(t, GenerateDBEntity(&e))
}
