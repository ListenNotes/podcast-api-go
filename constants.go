package listennotes

// Base urls for access the available api endpoints
const (
	BaseURLProduction = "https://listen-api.listennotes.com/api/v2"
	BaseURLTest       = "https://listen-api-test.listennotes.com/api/v2"
)

// Request header keys
const (
	RequestHeaderKeyAPI = "X-ListenAPI-Key"
)

// Reponse header keys
const (
	ResponseHeaderKeyFreeQuota       = "X-ListenAPI-FreeQuota"
	ResponseHeaderKeyUsage           = "X-ListenAPI-Usage"
	ResponseHeaderKeyLatencySeconds  = "X-listenAPI-Latency-Seconds"
	ResponseHeaderKeyNextBillingDate = "X-Listenapi-NextBillingDate"
)

// TimeFormat is the string format of all response times
const TimeFormat = "2006-01-02T15:04:05.999999-07:00"
