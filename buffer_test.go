package eventsender

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	buff := NewBuffer()

	n, err := buff.Write([]byte("1"))
	assert.NoError(t, err)
	assert.Equal(t, 1, n)

	n, err = buff.Write([]byte("2"))
	assert.NoError(t, err)
	assert.Equal(t, 1, n)

	restored := make([]byte, 32)

	n, err = buff.Read(restored)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.Equal(t, []byte("1"), restored[:1])

	n, err = buff.Read(restored)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.Equal(t, []byte("2"), restored[:1])

	assert.NoError(t, buff.Close())

	n, err = buff.Write([]byte("3"))
	assert.Empty(t, n)
	assert.Equal(t, "EOF", err.Error())

	n, err = buff.Read(restored)
	assert.Empty(t, n)
	assert.Equal(t, "EOF", err.Error())
}
