package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	total := 111
	from := 3
	to := 18
	nFrom, nTo := reverse(total, from, to)
	assert.Equal(t, 94, nFrom)
	assert.Equal(t, 109, nTo)
}
