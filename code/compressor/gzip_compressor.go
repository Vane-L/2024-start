package compressor

import (
	"bytes"
	"compress/gzip"
	"io"
)

type GzipCompressor struct{}

func (c GzipCompressor) Zip(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	writer := gzip.NewWriter(buf)
	defer func() {
		_ = writer.Close()
	}()
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = writer.Flush()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c GzipCompressor) Unzip(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	gzipReader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = gzipReader.Close()
	}()
	data, err = io.ReadAll(gzipReader)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, err
	}
	return data, nil
}
