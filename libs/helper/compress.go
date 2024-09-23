package helper

import (
	"bytes"
	"compress/gzip"
	"io"
)

// Fungsi untuk kompres data dengan gzip
func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Fungsi untuk mendekompres data dengan gzip
func Decompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var res bytes.Buffer
	if _, err := io.Copy(&res, r); err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}
