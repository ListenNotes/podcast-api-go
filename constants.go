package listennotes

// Version is the SDK version sent in the User-Agent header.
const Version = "3.0.0"

// Base urls for access the available api endpoints
const (
	BaseURLProduction = "https://listen-api.listennotes.com/api/v2"
	BaseURLTest       = "https://listen-api-test.listennotes.com/api/v2"
)

// Request header keys
const (
	RequestHeaderKeyAPI = "X-ListenAPI-Key"
)

// Response header keys
const (
	ResponseHeaderKeyFreeQuota       = "X-ListenAPI-FreeQuota"
	ResponseHeaderKeyUsage           = "X-ListenAPI-Usage"
	ResponseHeaderKeyLatencySeconds  = "X-listenAPI-Latency-Seconds"
	ResponseHeaderKeyNextBillingDate = "X-Listenapi-NextBillingDate"
)

// TimeFormat is the string format of all response times
const TimeFormat = "2006-01-02T15:04:05.999999-07:00"
