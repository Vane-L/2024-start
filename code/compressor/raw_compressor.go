package compressor

type RawCompressor struct{}

func (c RawCompressor) Zip(data []byte) ([]byte, error) {
	return data, nil
}

func (c RawCompressor) Unzip(data []byte) ([]byte, error) {
	return data, nil
}
