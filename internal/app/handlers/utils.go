package handlers

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func decompress(data []byte) ([]byte, error) {
	var res bytes.Buffer
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	_, err = res.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	err = r.Close()
	if err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}

func getRequestBody(r *http.Request) ([]byte, error) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}
	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		decompressBody, err := decompress(b)
		if err != nil {
			panic(err)
		}
		return decompressBody, nil
	}
	return b, nil
}
