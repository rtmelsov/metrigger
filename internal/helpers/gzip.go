package helpers

import (
	"bytes"
	"compress/gzip"
	"io"
)

// Функция для упаковки данных в gzip
func CompressData(data []byte) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, err := gzipWriter.Write(data)
	if err != nil {
		return nil, err
	}
	// Закрываем gzipWriter, чтобы завершить запись
	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}
	return &buf, nil
}

// Функция для распаковки данных из gzip
func DecompressData(data []byte) (*bytes.Buffer, error) {
	buf := bytes.NewReader(data)
	gzipReader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	var result bytes.Buffer
	_, err = io.Copy(&result, gzipReader)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
