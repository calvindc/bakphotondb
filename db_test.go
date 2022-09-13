package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupDB(t *testing.T) {
	os.Remove(filepath.Join(".", idFileName))
	SetupDB(".")
	assert.EqualValues(t, NextID(), 1)
	assert.EqualValues(t, NextID(), 2)
	idFile.Close()
	id = 0
	SetupDB(".")
	assert.EqualValues(t, NextID(), 3)
	assert.EqualValues(t, NextID(), 4)
	idFile.Close()
}
