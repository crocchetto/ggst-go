package ggst

import (
	"context"
	"fmt"

	"github.com/crocchetto/ggst-go/internal/protocol"
)

type Client struct {
	transport *protocol.Transport

	token    string
	playerID string
	name     string
}

func NewClient(baseURL string) *Client {
	return &Client{
		transport: protocol.NewTransport(baseURL),
	}
}

type Session struct {
	Token    string
	PlayerID string
	Name     string
	SteamID  string
}

func (c *Client) Login(ctx context.Context, steamID, steamHex, ticket string) (Session, error) {
	payload, err := protocol.EncodeLogin(steamID, steamHex, ticket)
	if err != nil {
		return Session{}, err
	}

	plaintext, err := c.transport.Do(ctx, protocol.PathLogin, payload)
	if err != nil {
		return Session{}, fmt.Errorf("ggst: login: %w", err)
	}

	res, err := protocol.DecodeLogin(plaintext)
	if err != nil {
		return Session{}, fmt.Errorf("ggst: login: %w", err)
	}

	c.token = res.Token
	c.playerID = res.PlayerID
	c.name = res.Name

	return Session{
		Token:    res.Token,
		PlayerID: res.PlayerID,
		Name:     res.Name,
		SteamID:  res.SteamID,
	}, nil
}

func (c *Client) RestoreSession(token, playerID string) {
	c.token = token
	c.playerID = playerID
}

func (c *Client) PlayerID() string { return c.playerID }

func (c *Client) Name() string { return c.name }

func (c *Client) authed() bool { return c.token != "" }

func (c *Client) Replays(ctx context.Context, index, replaysPerPage int) ([]Replay, error) {
	if !c.authed() {
		return nil, fmt.Errorf("ggst: replays: not authenticated")
	}
	if replaysPerPage < 1 || replaysPerPage > 127 {
		return nil, fmt.Errorf("ggst: replays: replaysPerPage must be 1..127, got %d", replaysPerPage)
	}

	payload, err := protocol.EncodeReplays(c.playerID, c.token, index, replaysPerPage, c.playerID)
	if err != nil {
		return nil, err
	}

	plaintext, err := c.transport.Do(ctx, protocol.PathReplays, payload)
	if err != nil {
		return nil, fmt.Errorf("ggst: replays: %w", err)
	}

	raw, err := protocol.DecodeReplays(plaintext)
	if err != nil {
		return nil, fmt.Errorf("ggst: replays: %w", err)
	}

	return mapReplays(raw), nil
}

func mapReplays(raw []protocol.RawReplay) []Replay {
	out := make([]Replay, 0, len(raw))
	for _, r := range raw {
		date, _ := r.ParsedDate()
		out = append(out, Replay{
			Floor:       Floor(r.Floor),
			Player1Char: Character(r.Player1Char),
			Player2Char: Character(r.Player2Char),
			Player1:     mapPlayer(r.Player1),
			Player2:     mapPlayer(r.Player2),
			Winner:      r.Winner,
			Date:        date,
			Views:       r.Views,
			Likes:       r.Likes,
		})
	}
	return out
}

func mapPlayer(p protocol.RawPlayer) Player {
	return Player{
		ID:   p.ID,
		Name: p.Name,
	}
}
