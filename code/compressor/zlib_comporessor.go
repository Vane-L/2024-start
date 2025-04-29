package compressor

import (
	"bytes"
	"compress/zlib"
	"io"
)

type ZlibCompressor struct{}

func (c ZlibCompressor) Zip(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	writer := zlib.NewWriter(buf)
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
func (c ZlibCompressor) Unzip(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	gzipReader, err := zlib.NewReader(buf)
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
