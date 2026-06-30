package twin

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// PublicRecord removes free-text notes and keeps only daily aggregates intended
// for examples or public dashboards.
func PublicRecord(r DailyRecord) DailyRecord {
	r.Notes = ""
	r.WeightKG = 0
	return r
}

type PrivacyIssue struct {
	Field   string `json:"field"`
	Reason  string `json:"reason"`
	Level   string `json:"level"`
	Example string `json:"example,omitempty"`
}

var forbiddenPublicHeaders = map[string]string{
	"gps":                   "raw location is not public-safe",
	"latitude":              "raw location is not public-safe",
	"longitude":             "raw location is not public-safe",
	"route":                 "raw route history is not public-safe",
	"address":               "home/work/clinic address is not public-safe",
	"access_token":          "device or account token must never be committed",
	"refresh_token":         "device or account token must never be committed",
	"password":              "credential material must never be committed",
	"medical_report":        "full medical reports stay private by default",
	"diagnosis":             "diagnosis history stays private by default",
	"medication":            "medication history stays private by default",
	"sleep_timeline":        "exact sleep timeline should be aggregated before publishing",
	"heart_rate_timeseries": "raw heart-rate time series should be aggregated before publishing",
}

func InspectCSVPrivacy(path string) ([]PrivacyIssue, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true
	header, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}

	var issues []PrivacyIssue
	for _, raw := range header {
		name := strings.ToLower(strings.TrimSpace(raw))
		if reason, ok := forbiddenPublicHeaders[name]; ok {
			issues = append(issues, PrivacyIssue{Field: raw, Reason: reason, Level: "block"})
		}
	}
	return issues, nil
}

// PrivacyWarnings flags obvious free-text markers that should be reviewed
// before publishing a record.
func PrivacyWarnings(r DailyRecord) []string {
	var warnings []string
	note := strings.ToLower(r.Notes)
	for _, marker := range []string{"home", "work", "clinic", "hospital", "gps", "address"} {
		if strings.Contains(note, marker) {
			warnings = append(warnings, "notes may contain location or sensitive context")
			break
		}
	}
	return warnings
}
