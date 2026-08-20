package protocol

import (
	"fmt"
	"strings"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

const dateLayout = "2006-01-02 15:04:05"

type RawReplay struct {
	Int1        uint64
	Int2        int64
	Floor       int
	Player1Char int
	Player2Char int
	Player1     RawPlayer
	Player2     RawPlayer
	Winner      int
	Date        string
	Int7        int64
	Views       uint64
	Int8        int64
	Likes       uint64
}

type RawPlayer struct {
	ID        string
	Name      string
	AccountID string
	OnlineID  string
	Int1      int64
	Int2      int64
	Int3      int64
	Int4      int64
}

func (r RawReplay) ParsedDate() (time.Time, error) {
	t, err := time.Parse(dateLayout, r.Date)
	if err != nil {
		return time.Time{}, fmt.Errorf("protocol: bad replay date %q: %w", r.Date, err)
	}
	return t, nil
}

type replayResponseEnvelope struct {
	_msgpack struct{} `msgpack:",as_array"`

	Header msgpack.RawMessage
	Body   replayResponseBody
}

type replayResponseBody struct {
	_msgpack struct{} `msgpack:",as_array"`

	Int1    int64
	Int2    int64
	Int3    int64
	Replays []RawReplay
}

func DecodeReplays(plaintext []byte) ([]RawReplay, error) {
	var env replayResponseEnvelope
	if err := msgpack.Unmarshal(plaintext, &env); err != nil {
		return nil, fmt.Errorf("protocol: decode replay response: %w", err)
	}
	return env.Body.Replays, nil
}

// --- login response ---

type LoginResult struct {
	Token    string
	Name     string
	SteamID  string
	PlayerID string
}

type loginRespHeader struct {
	_msgpack struct{} `msgpack:",as_array"`

	Token    string
	Int1     int64
	Date     string
	Version1 string
	Version2 string
	Version3 string
	String1  string
	String2  string
}

type loginRespInner struct {
	_msgpack struct{} `msgpack:",as_array"`

	PlayerID string
	Name     string
	SteamID  string
	SteamHex string
	Platform int64
}

type loginRespBody struct {
	_msgpack struct{} `msgpack:",as_array"`

	Int1 int64
	Data loginRespInner
}

type loginRespEnvelope struct {
	_msgpack struct{} `msgpack:",as_array"`

	Header loginRespHeader
	Body   loginRespBody
}

func DecodeLogin(plaintext []byte) (LoginResult, error) {
	var env loginRespEnvelope
	if err := msgpack.Unmarshal(plaintext, &env); err != nil {
		return LoginResult{}, fmt.Errorf("protocol: decode login response: %w", err)
	}
	if env.Header.Token == "" {
		return LoginResult{}, fmt.Errorf("protocol: login response has empty token")
	}
	return LoginResult{
		Token:    env.Header.Token,
		Name:     env.Body.Data.Name,
		SteamID:  env.Body.Data.SteamID,
		PlayerID: env.Body.Data.PlayerID,
	}, nil
}

// --- statistics response ---

func DecodeStatsJSON(plaintext []byte) (string, error) {
	var top []msgpack.RawMessage
	if err := msgpack.Unmarshal(plaintext, &top); err != nil {
		return "", fmt.Errorf("protocol: decode stats response: %w", err)
	}
	if len(top) < 2 {
		return "", fmt.Errorf("protocol: stats response has %d top-level elements, want >= 2", len(top))
	}

	var body []msgpack.RawMessage
	if err := msgpack.Unmarshal(top[1], &body); err != nil {
		return "", fmt.Errorf("protocol: decode stats body: %w", err)
	}

	for _, elem := range body {
		var str string
		if err := msgpack.Unmarshal(elem, &str); err != nil {
			continue
		}
		trimmed := strings.TrimSpace(str)
		if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
			return str, nil
		}
	}
	return "", fmt.Errorf("protocol: no JSON string found in stats body (%d elements)", len(body))
}
