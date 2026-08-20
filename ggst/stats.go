package ggst

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/crocchetto/ggst-go/internal/protocol"
)

type PlayerStats struct {
	Raw    string
	Fields map[string]any
}

func (s PlayerStats) String() string { return s.Raw }

func (s PlayerStats) Int(key string) (value int, ok bool) {
	v, present := s.Fields[key]
	if !present {
		return 0, false
	}
	f, isNum := v.(float64)
	if !isNum {
		return 0, false
	}
	return int(f), true
}

func (s PlayerStats) Str(key string) (value string, ok bool) {
	v, present := s.Fields[key]
	if !present {
		return "", false
	}
	str, isStr := v.(string)
	if !isStr {
		return "", false
	}
	return str, true
}

func (c *Client) PlayerStats(ctx context.Context, targetPlayerID string) (PlayerStats, error) {
	if !c.authed() {
		return PlayerStats{}, fmt.Errorf("ggst: stats: not authenticated")
	}
	// empty targetPlayerID asks the server for the caller's own statistics
	// setting a custom player ID returns a reduced result
	payload, err := protocol.EncodePlayerStats(c.playerID, c.token, targetPlayerID)
	if err != nil {
		return PlayerStats{}, err
	}

	plaintext, err := c.transport.Do(ctx, protocol.PathStatistics, payload)
	if err != nil {
		return PlayerStats{}, fmt.Errorf("ggst: stats: %w", err)
	}

	rawJSON, err := protocol.DecodeStatsJSON(plaintext)
	if err != nil {
		return PlayerStats{}, fmt.Errorf("ggst: stats: %w", err)
	}

	fields := map[string]any{}
	if err := json.Unmarshal([]byte(rawJSON), &fields); err != nil {
		return PlayerStats{Raw: rawJSON}, fmt.Errorf("ggst: stats: parse JSON: %w", err)
	}

	return PlayerStats{Raw: rawJSON, Fields: fields}, nil
}
