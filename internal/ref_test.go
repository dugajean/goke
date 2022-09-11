package internal

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRef(t *testing.T) {
	myRef := NewRef("hello", nil)
	assert.NotNil(t, myRef)
	assert.Nil(t, myRef.Error())
}

func TestRefEquals(t *testing.T) {
	myRef := NewRef(true, nil)
	assert.True(t, myRef.Value())
}

func TestRefError(t *testing.T) {
	myRef := NewRef("", errors.New("foo"))
	assert.Equal(t, "foo", myRef.Error().Error())
}
