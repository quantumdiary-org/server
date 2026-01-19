package encoding_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"netschool-proxy/api/api/internal/pkg/encoding"
)

func TestEncodeDecodeWindows1251(t *testing.T) {
	originalText := "Привет, мир!"
	
	encoded, err := encoding.EncodeWindows1251(originalText)
	assert.NoError(t, err)
	assert.NotEmpty(t, encoded)
	
	decoded, err := encoding.DecodeWindows1251(encoded)
	assert.NoError(t, err)
	assert.Equal(t, originalText, decoded)
}

func TestEncodeDecodeEnglishText(t *testing.T) {
	originalText := "Hello, world!"
	
	encoded, err := encoding.EncodeWindows1251(originalText)
	assert.NoError(t, err)
	assert.Equal(t, originalText, encoded) // For ASCII characters, encoding should be the same
	
	decoded, err := encoding.DecodeWindows1251(encoded)
	assert.NoError(t, err)
	assert.Equal(t, originalText, decoded)
}