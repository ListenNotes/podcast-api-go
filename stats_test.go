package listennotes

import (
	"net/http"
	"testing"
)

func TestParseStatsGoodValues(t *testing.T) {
	headers := http.Header{}
	headers.Set(ResponseHeaderKeyFreeQuota, "10")
	headers.Set(ResponseHeaderKeyUsage, "11")
	headers.Set(ResponseHeaderKeyLatencySeconds, "12.3")
	headers.Set(ResponseHeaderKeyNextBillingDate, "2020-09-26T17:27:33.110641+00:00")

	resp := &http.Response{
		Header: headers,
	}

	stats := parseStats(resp)

	if stats.FreeQuota != 10 {
		t.Errorf("FreeQuota did not parse correctly: %v", stats.FreeQuota)
	}

	if stats.Usage != 11 {
		t.Errorf("Usage did not parse correctly: %v", stats.Usage)
	}

	if stats.LatencySeconds != 12.3 {
		t.Errorf("LatencySeconds did not parse correctly: %v", stats.LatencySeconds)
	}

	if stats.NextBillingDate.Month() != 9 {
		t.Errorf("NextBillingDate did not parse correctly: %v", stats.NextBillingDate)
	}
}

func TestParseStatsBadValues(t *testing.T) {
	headers := http.Header{}
	headers.Set(ResponseHeaderKeyFreeQuota, "a")
	headers.Set(ResponseHeaderKeyUsage, "b")
	headers.Set(ResponseHeaderKeyLatencySeconds, "c")
	headers.Set(ResponseHeaderKeyNextBillingDate, "d")

	resp := &http.Response{
		Header: headers,
	}

	stats := parseStats(resp)

	if stats.FreeQuota != 0 {
		t.Errorf("FreeQuota did not parse correctly: %v", stats.FreeQuota)
	}

	if stats.Usage != 0 {
		t.Errorf("Usage did not parse correctly: %v", stats.Usage)
	}

	if stats.LatencySeconds != 0 {
		t.Errorf("LatencySeconds did not parse correctly: %v", stats.LatencySeconds)
	}

	if !stats.NextBillingDate.IsZero() {
		t.Errorf("NextBillingDate did not parse correctly: %v", stats.NextBillingDate)
	}
}
