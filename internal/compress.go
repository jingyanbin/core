package internal

import (
	"bytes"
	"compress/gzip"
	"io"
)

var Compress compress

type compress struct {
}

func (*compress) GZip(data []byte) []byte {
	var in bytes.Buffer
	defer in.Reset()
	w, err := gzip.NewWriterLevel(&in, 1)
	if err != nil {
		panic(err)
	}
	w.Write(data)
	w.Close()
	return in.Bytes()
}

func (*compress) UnGZip(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	unData, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return unData, nil
}
