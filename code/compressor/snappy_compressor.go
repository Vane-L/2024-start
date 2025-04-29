package compressor

import (
	"bytes"
	"io"

	"github.com/golang/snappy"
)

type SnappyCompressor struct{}

func (c SnappyCompressor) Zip(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	writer := snappy.NewBufferedWriter(buf)
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

func (c SnappyCompressor) Unzip(data []byte) ([]byte, error) {
	reader := snappy.NewReader(bytes.NewBuffer(data))
	data, err := io.ReadAll(reader)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, err
	}
	return data, nil
}
