package encoding

import (
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strings"
)


func DecodeWindows1251(input string) (string, error) {
	decoder := charmap.Windows1251.NewDecoder()
	result, _, err := transform.String(decoder, input)
	if err != nil {
		return "", err
	}
	return result, nil
}


func EncodeWindows1251(input string) (string, error) {
	encoder := charmap.Windows1251.NewEncoder()
	result, _, err := transform.String(encoder, input)
	if err != nil {
		return "", err
	}
	return result, nil
}


func ConvertReaderFromWindows1251(reader *strings.Reader) ([]byte, error) {
	decoder := charmap.Windows1251.NewDecoder()
	decodedReader := transform.NewReader(reader, decoder)
	return ioutil.ReadAll(decodedReader)
}