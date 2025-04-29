package compressor

import (
	"testing"
)

func TestRaw(t *testing.T) {
	raw := Compressors[Raw]
	rawResp, err := raw.Zip([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(rawResp))
	rawResp, err = raw.Unzip(rawResp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(rawResp))
}

func TestGzip(t *testing.T) {
	gzip := Compressors[Gzip]
	gzipResp, err := gzip.Zip([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	gzipResp, err = gzip.Unzip(gzipResp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(gzipResp))
}

func TestSnappy(t *testing.T) {
	snappy := Compressors[Snappy]
	snappyResp, err := snappy.Zip([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	snappyResp, err = snappy.Unzip(snappyResp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(snappyResp))
}

func TestZlib(t *testing.T) {
	zlib := Compressors[Zlib]
	zlibResp, err := zlib.Zip([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	zlibResp, err = zlib.Unzip(zlibResp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(zlibResp))
}
