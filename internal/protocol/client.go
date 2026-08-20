package protocol

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/crocchetto/ggst-go/internal/crypto"
)

type Transport struct {
	baseURL string
	http    *http.Client
}

func NewTransport(baseURL string) *Transport {
	if baseURL == "" {
		baseURL = BaseURL
	}
	return &Transport{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (t *Transport) Do(ctx context.Context, path string, payload []byte) ([]byte, error) {
	encrypted, err := crypto.Encrypt(payload)
	if err != nil {
		return nil, fmt.Errorf("transport: encrypt: %w", err)
	}
	form := url.Values{"data": {encrypted}}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		t.baseURL+path, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("transport: build request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Cache-Control", cacheControl)
	req.Header.Set("x-client-version", clientVersion)

	resp, err := t.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("transport: send: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("transport: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		snippet := body
		if len(snippet) > 256 {
			snippet = snippet[:256]
		}
		return nil, fmt.Errorf("transport: unexpected status %d, body: %q", resp.StatusCode, snippet)
	}

	plaintext, err := crypto.Decrypt(body)
	if err != nil {
		return nil, fmt.Errorf("transport: decrypt response: %w", err)
	}
	return plaintext, nil
}
