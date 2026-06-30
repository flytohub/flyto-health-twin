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

func TestEvaluateRequiresEnoughRecords(t *testing.T) {
	_, err := Evaluate([]DailyRecord{
		{Date: day("2026-06-01"), HRV: 50},
		{Date: day("2026-06-02"), HRV: 51},
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPublicRecordDropsNotes(t *testing.T) {
	got := PublicRecord(DailyRecord{Notes: "home location", WeightKG: 72.1})
	if got.Notes != "" {
		t.Fatalf("expected notes to be redacted, got %q", got.Notes)
	}
	if got.WeightKG != 0 {
		t.Fatalf("expected weight to be redacted, got %v", got.WeightKG)
	}
}

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

func day(s string) time.Time {
	t, err := time.Parse(DateLayout, s)
	if err != nil {
		panic(err)
	}
	return t
}
