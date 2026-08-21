# ggst-go

the API is private and undocumented; this client was built by observing the game's own traffic

## install

```
go get github.com/crocchetto/ggst-go
```

## how it works

every request and response is a MessagePack array encrypted with AES-256-GCM and sent form-encoded to the API,
the authentication uses a short-lived Steam web-API auth ticket, which the server exchanges for a session during login

the library is split:

- **`ggst`** — the core client
- **`auth`** — obtains the Steam credentials (SteamID and web-API ticket) needed for a fresh login, requires a running Steam client

## usage

the typical flow is first acquire Steam credentials, log in, and then call the data endpoints

here's an example:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/crocchetto/ggst-go/auth"
	"github.com/crocchetto/ggst-go/ggst"
)

func main() {
	// initialize Steam
	if err := auth.Init(); err != nil {
		log.Fatal(err)
	}
	creds, err := auth.Acquire()
	if err != nil {
		log.Fatal(err)
	}

	// log in to the GGST API
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := ggst.NewClient("")
	sess, err := client.Login(ctx, creds.SteamID, creds.SteamHex, creds.Ticket)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Logged in as %s (player ID %s)\n", sess.Name, sess.PlayerID)

	// fetch your own statistics
	stats, err := client.PlayerStats(ctx, "")
	if err != nil {
		log.Fatal(err)
	}
	if rating, ok := stats.Int("CHP_RankMatchRating"); ok {
		fmt.Println("Chipp rank-match rating:", rating)
	}

	// fetch recent replays
	replays, err := client.Replays(ctx, 0, 10)
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range replays {
		fmt.Printf("%s (%s) vs %s (%s) on %s\n",
			r.Player1.Name, r.Player1Char,
			r.Player2.Name, r.Player2Char, r.Date.Format("2006-01-02"))
	}
}
```

### reusing a session

`Login` returns a `Session` containing a token and player ID, if you cache these, you can skip a fresh login with `RestoreSession`:

```go
client := ggst.NewClient("")
client.RestoreSession(token, playerID)
```

## API overview

### package `ggst`

- `NewClient(baseURL string) *Client` — create a client. Pass `""` for the real API
- `(*Client) Login(ctx, steamID, steamHex, ticket) (Session, error)`
- `(*Client) RestoreSession(token, playerID string)`
- `(*Client) PlayerStats(ctx, targetPlayerID string) (PlayerStats, error)` — pass `""` for your own full statistics
- `(*Client) Replays(ctx, index, replaysPerPage int) ([]Replay, error)`
- `PlayerStats` — holds the raw stats JSON plus `Int(key)` / `Str(key)` accessors for reading fields
- `Character`, `Floor` — typed enums with `Name()` / `IsValid()`
- `Replay`, `Player` — decoded replay data, with `Player1Won()` and `WinnerPlayer()` helpers

### package `auth`

- `Init() error` — initialize Steam for GGST
- `Acquire() (Credentials, error)` — get SteamID and a fresh web-API ticket
- `AcquireWithTimeout(timeout) (Credentials, error)`

## statistics fields

`PlayerStats` gives you two ways to read the server's statistics document

### typed accessors (recommended)

profile-level fields and per-character fields have typed accessors, so you don't need to know the raw JSON key strings

each returns a value and an `ok` flag that is `false` when the field is missing:

```go
stats, _ := client.PlayerStats(ctx, "")

// profile fields
name, _ := stats.NickName()
matches, _ := stats.TotalRankMatch()
dollars, _ := stats.WorldDollar()
fmt.Printf("%s — %d ranked matches, %d World$\n", name, matches, dollars)

// per-character fields, addressed by Character
chp := stats.Character(ggst.Chipp)
level, _ := chp.Level()
rating, ok := chp.RankRating()
if ok {
    fmt.Printf("Chipp: level %d, rank rating %d\n", level, rating)
}
```

profile accessors include `NickName`, `PublicComment`, `UserID`, `OnlineID`,
`LobbyRank`, `MaxLobbyRank`, `VipStatus`, `WorldDollar`, `TotalRankMatch`,
`TotalPlayTime`, `WinChainMax`, `WinChainNow`, and `Rating3on3`

per-character accessors (on `stats.Character(c)`) include `Level`, `Exp`,
`NextLevelExp`, `RankRating`, `RankRatingPt`, `MasterRatingPt`,
`PlayerMatchWins`, `WinChainMax`, `WinChainNow`, `RankWinChainMax`, and
`RankWinChainNow`

### raw access

for any field without a typed accessor, read it directly with `Int(key)` / `Str(key)`, or inspect the whole document via `stats.Raw`

per-character keys use a three-letter prefix, e.g. `CHP_` for Chipp: `CHP_Lv`, `CHP_RankMatchRating`

character prefixes: `SOL KYK MAY AXL CHP POT FAU MLL ZAT RAM LEO NAG GIO ANJ INO GLD JKO COS BKN TST BGT SIN BED ASK JHN ELP ABA SLY DZY VEN UNI LUC USG ABS`

---

<p align="center">
  <img alt="ggst-go" title="ggst-go" src="https://www.guiltygear.com/ggst/jp/wordpress/wp-content/themes/ggst/img/fankit/01_chibi_sol.png" width="200">
  <img alt="ggst-go" title="ggst-go" src="https://www.guiltygear.com/ggst/jp/wordpress/wp-content/themes/ggst/img/fankit/02_chibi_ky.png" width="200">
</p>
