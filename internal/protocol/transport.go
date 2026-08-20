package protocol

// HTTP transport constants
const (
	BaseURL = "https://ggst-game.guiltygear.com"

	// endpoint paths
	PathLogin      = "/api/user/login"
	PathReplays    = "/api/catalog/get_replay"
	PathStatistics = "/api/statistics/get"

	// HTTP header values
	userAgent     = "GGST/Steam"
	contentType   = "application/x-www-form-urlencoded"
	cacheControl  = "no-store"
	clientVersion = "1"
)
