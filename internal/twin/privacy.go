package twin

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// PublicRecord removes private and equipment-gated values from a daily record.
// Public JSON should still use PublicRecordJSON so private fields are omitted,
// rather than serialized as zero values.
func PublicRecord(r DailyRecord) DailyRecord {
	r.Notes = ""
	r.WeightKG = 0
	r.BloodPressureSystolic = 0
	r.BloodPressureDiastolic = 0
	r.BloodGlucoseMGDL = 0
	r.BodyTemperatureC = 0
	r.IllnessScore = 0
	r.TrainingLoad = 0
	return r
}

// PrivacyIssue describes one blocked or manually reviewed source-data finding.
type PrivacyIssue struct {
	Field   string `json:"field"`
	Reason  string `json:"reason"`
	Level   string `json:"level"`
	Example string `json:"example,omitempty"`
}

// forbiddenPublicHeaders maps raw sensitive fields to operator-facing reasons.
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

// InspectCSVPrivacy blocks known raw or credential-bearing CSV headers.
func InspectCSVPrivacy(path string) (issues []PrivacyIssue, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close CSV %q: %w", path, closeErr)
		}
	}()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true
	header, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}

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
