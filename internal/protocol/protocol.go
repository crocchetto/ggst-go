package protocol

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
)

// protocol constants
const (
	APIVersion  = "0.4.8"
	PlatformPC  = 3
	headerInt1  = 10
	loginBodyID = 1
)

type Header struct {
	_msgpack struct{} `msgpack:",as_array"`

	PlayerID string
	Token    string
	Int1     int
	Version  string
	Platform int
}

func newHeader(playerID, token string) Header {
	return Header{
		PlayerID: playerID,
		Token:    token,
		Int1:     headerInt1,
		Version:  APIVersion,
		Platform: PlatformPC,
	}
}

type loginBody struct {
	_msgpack struct{} `msgpack:",as_array"`

	Int1       int
	SteamID    string
	SteamHex   string
	Int2       int
	SteamToken string
}

type loginRequest struct {
	_msgpack struct{} `msgpack:",as_array"`

	Header Header
	Body   loginBody
}

func EncodeLogin(steamID, steamHex, steamToken string) ([]byte, error) {
	req := loginRequest{
		Header: newHeader("", ""),
		Body: loginBody{
			Int1:       loginBodyID,
			SteamID:    steamID,
			SteamHex:   steamHex,
			Int2:       256,
			SteamToken: steamToken,
		},
	}

	data, err := msgpack.Marshal(&req)
	if err != nil {
		return nil, fmt.Errorf("protocol: encode login: %w", err)
	}
	return data, nil
}
