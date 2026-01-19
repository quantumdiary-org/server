package encoding

import (
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strings"
)

// DecodeWindows1251 decodes a Windows-1251 encoded string to UTF-8
func DecodeWindows1251(input string) (string, error) {
	decoder := charmap.Windows1251.NewDecoder()
	result, _, err := transform.String(decoder, input)
	if err != nil {
		return "", err
	}
	return result, nil
}

// EncodeWindows1251 encodes a UTF-8 string to Windows-1251
func EncodeWindows1251(input string) (string, error) {
	encoder := charmap.Windows1251.NewEncoder()
	result, _, err := transform.String(encoder, input)
	if err != nil {
		return "", err
	}
	return result, nil
}

// ConvertReaderFromWindows1251 converts a reader with Windows-1251 content to UTF-8
func ConvertReaderFromWindows1251(reader *strings.Reader) ([]byte, error) {
	decoder := charmap.Windows1251.NewDecoder()
	decodedReader := transform.NewReader(reader, decoder)
	return ioutil.ReadAll(decodedReader)
}