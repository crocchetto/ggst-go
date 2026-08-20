package ggst

// --- profile-level fields ---

func (s PlayerStats) NickName() (string, bool)      { return s.Str("NickName") }
func (s PlayerStats) PublicComment() (string, bool) { return s.Str("PublicComment") }
func (s PlayerStats) UserID() (int, bool)           { return s.Int("UserID") }
func (s PlayerStats) OnlineID() (string, bool)      { return s.Str("OnlineID") }
func (s PlayerStats) LobbyRank() (int, bool)        { return s.Int("LobbyRank") }
func (s PlayerStats) MaxLobbyRank() (int, bool)     { return s.Int("MaxLobbyRank") }
func (s PlayerStats) VipStatus() (int, bool)        { return s.Int("VipStatus") }
func (s PlayerStats) WorldDollar() (int, bool)      { return s.Int("WorldDollar") }
func (s PlayerStats) TotalRankMatch() (int, bool)   { return s.Int("TotalRankMatch") }
func (s PlayerStats) TotalPlayTime() (int, bool)    { return s.Int("TotalPlayTime") }
func (s PlayerStats) WinChainMax() (int, bool)      { return s.Int("PlayerWinChainMax") }
func (s PlayerStats) WinChainNow() (int, bool)      { return s.Int("PlayerWinChainNow") }
func (s PlayerStats) Rating3on3() (int, bool)       { return s.Int("3on3Rating") }

// --- per-character fields ---

type CharacterStats struct {
	stats  PlayerStats
	prefix string
}

func (s PlayerStats) Character(c Character) CharacterStats {
	return CharacterStats{stats: s, prefix: c.statPrefix()}
}

func (cs CharacterStats) field(suffix string) (int, bool) {
	if cs.prefix == "" {
		return 0, false
	}
	return cs.stats.Int(cs.prefix + "_" + suffix)
}

func (cs CharacterStats) Level() (int, bool)           { return cs.field("Lv") }
func (cs CharacterStats) Exp() (int, bool)             { return cs.field("Exp") }
func (cs CharacterStats) NextLevelExp() (int, bool)    { return cs.field("NextLvExp") }
func (cs CharacterStats) RankRating() (int, bool)      { return cs.field("RankMatchRating") }
func (cs CharacterStats) RankRatingPt() (int, bool)    { return cs.field("RankMatchRatingPt") }
func (cs CharacterStats) MasterRatingPt() (int, bool)  { return cs.field("MasterRatingPt") }
func (cs CharacterStats) PlayerMatchWins() (int, bool) { return cs.field("PM_Wins") }
func (cs CharacterStats) WinChainMax() (int, bool)     { return cs.field("WinChainMax") }
func (cs CharacterStats) WinChainNow() (int, bool)     { return cs.field("WinChainNow") }
func (cs CharacterStats) RankWinChainMax() (int, bool) { return cs.field("RankMatchWinChainMax") }
func (cs CharacterStats) RankWinChainNow() (int, bool) { return cs.field("RankMatchWinChainNow") }
