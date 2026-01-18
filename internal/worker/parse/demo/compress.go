package demo

import (
	"bytes"
	"fmt"
	"io"

	"github.com/klauspost/compress/zstd"
	"github.com/vmihailenco/msgpack/v5"
)

func (m *Match) Compress() ([]byte, error) {
	if m == nil {
		return nil, fmt.Errorf("match is nil")
	}

	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	enc.UseCompactInts(true)

	if err := enc.Encode(m); err != nil {
		return nil, fmt.Errorf("serialize match %w", err)
	}

	var compressed bytes.Buffer
	writer, err := zstd.NewWriter(&compressed, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	if err != nil {
		return nil, fmt.Errorf("create zstd writer %w", err)
	}

	if _, err := writer.Write(buf.Bytes()); err != nil {
		return nil, fmt.Errorf("write compressed data %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close zstd writer %w", err)
	}

	return compressed.Bytes(), nil
}

func Decompress(data []byte) (*Match, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("compressed data is empty")
	}

	reader, err := zstd.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create zstd reader %w", err)
	}
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("decompress data %w", err)
	}

	var m Match
	if err := msgpack.Unmarshal(decompressed, &m); err != nil {
		return nil, fmt.Errorf("deserialize match %w", err)
	}

	return &m, nil
}
