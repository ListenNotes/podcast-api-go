package podcast

import (
	"net/http"
	"strconv"
	"time"
)

type ResponseStatistics struct {
	FreeQuota       *int
	Usage           *int
	LatencySeconds  *float64
	NextBillingDate *time.Time
}

func parseStats(resp *http.Response) ResponseStatistics {
	// Do not allow this to fail, just return null values if we cannot parse the response.
	// We do this so that a POST that succeeds does not return an unrelated error.

	stats := ResponseStatistics{}

	if freeQuota, err := strconv.Atoi(resp.Header.Get(ResponseHeaderKeyFreeQuota)); err == nil {
		stats.FreeQuota = &freeQuota
	}
	if usage, err := strconv.Atoi(resp.Header.Get(ResponseHeaderKeyUsage)); err == nil {
		stats.Usage = &usage
	}
	if latency, err := strconv.ParseFloat(resp.Header.Get(ResponseHeaderKeyLatencySeconds), 64); err == nil {
		stats.LatencySeconds = &latency
	}
	if nextBill, err := time.Parse(resp.Header.Get(ResponseHeaderKeyNextBillingDate), ""); err == nil {
		stats.NextBillingDate = &nextBill
	}

	return stats
}
