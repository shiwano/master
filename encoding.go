package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	"io"
	"regexp"
)

var (
	utf8BOM            = []byte{239, 187, 191}
	utf8CharsetPattern = regexp.MustCompile("(?i)utf-?8$")
)

func isUTF8Charset(charsetName string) bool {
	return utf8CharsetPattern.MatchString(charsetName)
}

func addBOM(data []byte) []byte {
	return append(utf8BOM, data...)
}

func stripBOM(data []byte) []byte {
	return bytes.TrimPrefix(data, utf8BOM)
}

func encode(data []byte, charsetName string) ([]byte, error) {
	encoding, _ := charset.Lookup(charsetName)
	if encoding == nil {
		return nil, fmt.Errorf("Unsupported charset: %v", charsetName)
	}

	reader := bytes.NewReader(data)
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	encodeWriter := transform.NewWriter(writer, encoding.NewEncoder())
	if _, err := io.Copy(encodeWriter, reader); err != nil {
		return nil, err
	}
	if err := writer.Flush(); err != nil {
		return nil, err
	}

	if isUTF8Charset(charsetName) {
		return addBOM(b.Bytes()), nil
	}
	return b.Bytes(), nil
}

func decode(data []byte, charsetName string) ([]byte, error) {
	encoding, _ := charset.Lookup(charsetName)
	if encoding == nil {
		return nil, fmt.Errorf("Unsupported charset: %v", charsetName)
	}

	reader := bytes.NewReader(data)
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	decodeReader := transform.NewReader(reader, encoding.NewDecoder())
	if _, err := io.Copy(writer, decodeReader); err != nil {
		return nil, err
	}
	if err := writer.Flush(); err != nil {
		return nil, err
	}

	if isUTF8Charset(charsetName) {
		return stripBOM(b.Bytes()), nil
	}
	return b.Bytes(), nil
}

func detectEncoding(data []byte) (string, error) {
	detector := chardet.NewTextDetector()
	detected, err := detector.DetectBest(data)
	if err != nil {
		return "", err
	}
	return detected.Charset, nil
}
