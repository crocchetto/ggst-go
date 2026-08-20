package protocol

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
)

// --- player statistics: statistics/get ---

type StatsType int

const StatsTypeCharacterXP StatsType = 7

type statsBody struct {
	_msgpack struct{} `msgpack:",as_array"`

	PlayerID  string
	StatsType StatsType
	Sel1      int
	Sel2      int
	Sel3      int
	Sel4      int
	Trailer   int
}

type statsRequest struct {
	_msgpack struct{} `msgpack:",as_array"`

	Header Header
	Body   statsBody
}

func EncodePlayerStats(sessionPlayerID, token, targetPlayerID string) ([]byte, error) {
	req := statsRequest{
		Header: newHeader(sessionPlayerID, token),
		Body: statsBody{
			PlayerID:  targetPlayerID,
			StatsType: StatsTypeCharacterXP,
			Sel1:      -1,
			Sel2:      -1,
			Sel3:      -1,
			Sel4:      -1,
			Trailer:   1,
		},
	}
	data, err := msgpack.Marshal(&req)
	if err != nil {
		return nil, fmt.Errorf("protocol: encode player stats: %w", err)
	}
	return data, nil
}

// --- replays: catalog/get_replay ---

type replayQuery struct {
	_msgpack struct{} `msgpack:",as_array"`

	Int1     int
	Int2     int
	MinRank  int
	MaxRank  int
	MinFloor int
	MaxFloor int
	Seq      []any
	Char1    int
	Char2    int
	Winner   int
	BestBout int
	Mode     int
}

type replayBody struct {
	_msgpack struct{} `msgpack:",as_array"`

	Int1           int
	Index          int
	ReplaysPerPage int
	Query          replayQuery
	Platforms      int
}

type replayRequest struct {
	_msgpack struct{} `msgpack:",as_array"`

	Header Header
	Body   replayBody
}

// replay query defaults
const (
	replayMaxRank    = 22
	replayModeRanked = 7
	replayPlatforms  = 6
)

func EncodeReplays(sessionPlayerID, token string, index, replaysPerPage int, searchPlayerID string) ([]byte, error) {
	seq := []any{}
	if searchPlayerID != "" {
		seq = []any{searchPlayerID}
	}
	req := replayRequest{
		Header: newHeader(sessionPlayerID, token),
		Body: replayBody{
			Int1:           1,
			Index:          index,
			ReplaysPerPage: replaysPerPage,
			Query: replayQuery{
				Int1:     -1,
				Int2:     0,
				MinRank:  0,
				MaxRank:  replayMaxRank,
				MinFloor: 0,
				MaxFloor: 99,
				Seq:      seq,
				Char1:    -1,
				Char2:    -1,
				Winner:   0,
				BestBout: 1,
				Mode:     -1,
			},
			Platforms: replayPlatforms,
		},
	}
	data, err := msgpack.Marshal(&req)
	if err != nil {
		return nil, fmt.Errorf("protocol: encode replays: %w", err)
	}
	return data, nil
}
