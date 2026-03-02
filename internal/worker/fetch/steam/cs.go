package steam

import (
	"compress/bzip2"
	"context"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func (s *steam) downloadDemo(ctx context.Context, demoURL string) ([]byte, error) {
	zap.S().Debug("Downloading demo")
	req, err := http.NewRequestWithContext(ctx, "GET", demoURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("new request %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadGateway {
			// Demo file expired
			return nil, nil
		}
		return nil, fmt.Errorf("unexpected status code %s", resp.Status)
	}

	unzippedReader := bzip2.NewReader(resp.Body)

	bytes, err := io.ReadAll(unzippedReader)
	if err != nil {
		return nil, fmt.Errorf("read body %w", err)
	}

	return bytes, nil
}
