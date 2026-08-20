package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"

	steamworks "github.com/crocchetto/go-steamworks"
)

const GGSTAppID = 1384160

const ggstIdentity = "ggst-game.guiltygear.com"

const defaultTicketTimeout = 10 * time.Second

type Credentials struct {
	SteamID  string
	SteamHex string
	Ticket   string
}

func Init() error {
	if err := os.Setenv("SteamAppId", strconv.Itoa(GGSTAppID)); err != nil {
		return fmt.Errorf("auth: set SteamAppId: %w", err)
	}
	if err := steamworks.Init(); err != nil {
		return fmt.Errorf("auth: steam init failed: %w", err)
	}
	return nil
}

func Acquire() (Credentials, error) {
	return AcquireWithTimeout(defaultTicketTimeout)
}

func AcquireWithTimeout(timeout time.Duration) (Credentials, error) {
	id := uint64(steamworks.SteamUser().GetSteamID())
	if id == 0 {
		return Credentials{}, fmt.Errorf("auth: no logged-in Steam user (SteamID is 0)")
	}

	ticket, err := steamworks.GetAuthTicketForWebApiBlocking(ggstIdentity, timeout)
	if err != nil {
		return Credentials{}, fmt.Errorf("auth: %w", err)
	}

	return Credentials{
		SteamID:  strconv.FormatUint(id, 10),
		SteamHex: strconv.FormatUint(id, 16),
		Ticket:   ticket,
	}, nil
}
