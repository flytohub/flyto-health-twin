package twin

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestPredictNextRespondsToStrain verifies transparent strain classification.
func TestPredictNextRespondsToStrain(t *testing.T) {
	records := []DailyRecord{
		{Date: day("2026-06-01"), SleepScore: 82, SleepHours: 7.8, HRV: 55, RestingHeartRate: 60, FatigueScore: 2},
		{Date: day("2026-06-02"), SleepScore: 80, SleepHours: 7.6, HRV: 54, RestingHeartRate: 61, FatigueScore: 3},
		{Date: day("2026-06-03"), SleepScore: 58, SleepHours: 5.4, HRV: 46, RestingHeartRate: 66, ExerciseMinutes: 90, StressScore: 8, FatigueScore: 6, CaffeineServings: 4, WaterLiters: 1.6},
	}

	prediction, err := PredictNext(records)
	if err != nil {
		t.Fatal(err)
	}
	if prediction.RecoveryState != "strained" {
		t.Fatalf("expected strained recovery state, got %q", prediction.RecoveryState)
	}
	if len(prediction.Hints) < 3 {
		t.Fatalf("expected multiple strain hints, got %#v", prediction.Hints)
	}
	if prediction.ModelID == "" || prediction.ModelVersion == "" {
		t.Fatalf("expected model metadata, got %#v", prediction)
	}
	if len(prediction.FeatureSet) == 0 || len(prediction.Assumptions) == 0 {
		t.Fatalf("expected model feature and assumption metadata, got %#v", prediction)
	}
}

// TestEvaluateRequiresEnoughRecords verifies the rolling-history minimum.
func TestEvaluateRequiresEnoughRecords(t *testing.T) {
	_, err := Evaluate([]DailyRecord{
		{Date: day("2026-06-01"), HRV: 50},
		{Date: day("2026-06-02"), HRV: 51},
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

// TestPublicRecordDropsPrivateAndEquipmentGatedValues protects the public boundary.
func TestPublicRecordDropsPrivateAndEquipmentGatedValues(t *testing.T) {
	got := PublicRecord(DailyRecord{
		Notes: "home location", WeightKG: 72.1,
		BloodPressureSystolic: 120, BloodPressureDiastolic: 80,
		BloodGlucoseMGDL: 95, BodyTemperatureC: 36.8,
		IllnessScore: 3, TrainingLoad: 70,
	})
	if got.Notes != "" {
		t.Fatalf("expected notes to be redacted, got %q", got.Notes)
	}
	if got.WeightKG != 0 {
		t.Fatalf("expected weight to be redacted, got %v", got.WeightKG)
	}
	if got.BloodPressureSystolic != 0 || got.BloodPressureDiastolic != 0 ||
		got.BloodGlucoseMGDL != 0 || got.BodyTemperatureC != 0 ||
		got.IllnessScore != 0 || got.TrainingLoad != 0 {
		t.Fatalf("expected equipment-gated values to be redacted, got %#v", got)
	}
}

// TestPublicExportOmitsPrivateFields checks serialized field allowlisting.
func TestPublicExportOmitsPrivateFields(t *testing.T) {
	records := []DailyRecord{
		{Date: day("2026-06-01"), SleepScore: 80, SleepHours: 7.5, HRV: 55, RestingHeartRate: 60, FatigueScore: 2, Notes: "home", WeightKG: 72},
		{Date: day("2026-06-02"), SleepScore: 78, SleepHours: 7.2, HRV: 54, RestingHeartRate: 61, FatigueScore: 3},
		{Date: day("2026-06-03"), SleepScore: 76, SleepHours: 7.0, HRV: 53, RestingHeartRate: 62, FatigueScore: 3},
	}
	evals, err := Evaluate(records)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := WritePublicExport(&buf, records, evals); err != nil {
		t.Fatal(err)
	}
	raw := buf.String()
	if strings.Contains(raw, "home") || strings.Contains(raw, "weight_kg") {
		t.Fatalf("public export leaked private fields: %s", raw)
	}
	var decoded PublicExport
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Project != "flyto2" || len(decoded.Evaluations) == 0 {
		t.Fatalf("unexpected export: %#v", decoded)
	}
}

// TestPublicExportUsesProvidedGeneratedAt checks reproducible timestamps.
func TestPublicExportUsesProvidedGeneratedAt(t *testing.T) {
	records := []DailyRecord{
		{Date: day("2026-06-01"), SleepScore: 80, SleepHours: 7.5, HRV: 55, RestingHeartRate: 60, FatigueScore: 2},
		{Date: day("2026-06-02"), SleepScore: 78, SleepHours: 7.2, HRV: 54, RestingHeartRate: 61, FatigueScore: 3},
		{Date: day("2026-06-03"), SleepScore: 76, SleepHours: 7.0, HRV: 53, RestingHeartRate: 62, FatigueScore: 3},
	}
	evals, err := Evaluate(records)
	if err != nil {
		t.Fatal(err)
	}
	generatedAt := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)
	export := BuildPublicExportAt(records, evals, generatedAt)
	if export.GeneratedAt != "2026-06-30T00:00:00Z" {
		t.Fatalf("unexpected generated_at %q", export.GeneratedAt)
	}
}

// TestInspectCSVPrivacyBlocksRawFields checks known sensitive headers.
func TestInspectCSVPrivacyBlocksRawFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "raw.csv")
	err := os.WriteFile(path, []byte("date,hrv,gps,access_token\n2026-06-01,50,25.0/token,secret\n"), 0o600)
	if err != nil {
		t.Fatal(err)
	}
	issues, err := InspectCSVPrivacy(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 2 {
		t.Fatalf("expected two privacy issues, got %#v", issues)
	}
}

// TestLoadCSVRejectsNonFiniteNumbers prevents invalid downstream JSON values.
func TestLoadCSVRejectsNonFiniteNumbers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "non-finite.csv")
	if err := os.WriteFile(path, []byte("date,hrv\n2026-06-01,NaN\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadCSV(path); err == nil || !strings.Contains(err.Error(), "must be finite") {
		t.Fatalf("expected finite-number error, got %v", err)
	}
}

// TestLoadCSVRejectsDuplicateDates protects one-record-per-day semantics.
func TestLoadCSVRejectsDuplicateDates(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "duplicate.csv")
	content := "date,hrv\n2026-06-01,50\n2026-06-01,51\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadCSV(path); err == nil || !strings.Contains(err.Error(), "duplicate daily record") {
		t.Fatalf("expected duplicate-date error, got %v", err)
	}
}

// TestEvaluateRejectsUnorderedRecords protects rolling target alignment.
func TestEvaluateRejectsUnorderedRecords(t *testing.T) {
	_, err := Evaluate([]DailyRecord{
		{Date: day("2026-06-02")},
		{Date: day("2026-06-01")},
		{Date: day("2026-06-03")},
	})
	if err == nil || !strings.Contains(err.Error(), "strictly increasing") {
		t.Fatalf("expected chronology error, got %v", err)
	}
}

// day parses a fixture date and panics only for invalid test literals.
func day(s string) time.Time {
	t, err := time.Parse(DateLayout, s)
	if err != nil {
		panic(err)
	}
	return t
}
