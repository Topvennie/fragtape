package steam

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var downloadClient *http.Client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}
var ErrDemoExpired = errors.New("demo expired")

// Download downloads a demo from the steam servers
func (s *steam) Download(ctx context.Context, demoURL string) ([]byte, error) {
	var lastErr error

	for range 3 {
		data, err := s.downloadOnce(ctx, demoURL)
		if err == nil {
			return data, nil
		}

		lastErr = err

		if errors.Is(err, ErrDemoExpired) {
			return nil, err
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		case <-time.After(1 * time.Second):
		}
	}

	return nil, fmt.Errorf("download demo: %w", lastErr)
}

func (s *steam) downloadOnce(ctx context.Context, demoURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, demoURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	resp, err := downloadClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadGateway {
			return nil, ErrDemoExpired
		}
		if resp.StatusCode >= 500 && resp.StatusCode <= 599 {
			return nil, fmt.Errorf("temporary upstream status: %s", resp.Status)
		}

		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	compressed, err := readCompressedBody(ctx, resp)
	if err != nil {
		return nil, err
	}

	data, err := decompressBZ2(ctx, compressed)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func readCompressedBody(ctx context.Context, resp *http.Response) ([]byte, error) {
	var compressed bytes.Buffer

	reader := bufio.NewReaderSize(resp.Body, 1<<20) // 1 MB
	copyBuf := make([]byte, 1<<20)                  // 1 MB

	n, err := io.CopyBuffer(&compressed, reader, copyBuf)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, fmt.Errorf("read compressed response eof: %w", err)
		}

		return nil, fmt.Errorf("read compressed response: %w", err)
	}

	if resp.ContentLength >= 0 && n != resp.ContentLength {
		return nil, fmt.Errorf(
			"incomplete compressed response: got %d bytes, expected %d bytes",
			n,
			resp.ContentLength,
		)
	}

	return compressed.Bytes(), nil
}

func decompressBZ2(ctx context.Context, compressed []byte) ([]byte, error) {
	bz2Reader := bzip2.NewReader(bytes.NewReader(compressed))

	var out bytes.Buffer
	out.Grow(220 << 20) // 220 MB

	copyBuf := make([]byte, 1<<20) // 1 MB
	if _, err := io.CopyBuffer(&out, bz2Reader, copyBuf); err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, fmt.Errorf("decompress bz2 stream: %w", err)
		}

		return nil, fmt.Errorf("decompress bz2 stream: %w", err)
	}

	return out.Bytes(), nil
}
