package twin

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"time"
)

// SyntheticProfile describes a deterministic privacy-safe fixture scenario.
type SyntheticProfile struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	PrivacyProfile string `json:"privacy_profile"`
	DefaultDays    int    `json:"default_days"`
	SourceID       string `json:"source_id"`
}

// SyntheticProfiles returns all supported fixture scenarios.
func SyntheticProfiles() []SyntheticProfile {
	return []SyntheticProfile{
		{
			ID:             "balanced",
			Name:           "Balanced recovery loop",
			Description:    "Alternates normal training, strain, and recovery days for benchmark demos.",
			PrivacyProfile: "synthetic_daily_aggregate",
			DefaultDays:    30,
			SourceID:       "synthetic_balanced",
		},
		{
			ID:             "high-strain",
			Name:           "High strain period",
			Description:    "Stress, short sleep, and training load stay elevated to test strained predictions.",
			PrivacyProfile: "synthetic_daily_aggregate",
			DefaultDays:    30,
			SourceID:       "synthetic_high_strain",
		},
		{
			ID:             "missing-sensors",
			Name:           "Missing optional sensors",
			Description:    "Keeps core wearable fields but omits optional blood pressure, glucose, illness, and training load.",
			PrivacyProfile: "synthetic_daily_aggregate",
			DefaultDays:    30,
			SourceID:       "synthetic_missing_sensors",
		},
	}
}

// GenerateSyntheticRecords creates deterministic records for one profile.
func GenerateSyntheticRecords(profileID string, start time.Time, days int) ([]DailyRecord, error) {
	profile, ok := findSyntheticProfile(profileID)
	if !ok {
		return nil, fmt.Errorf("unknown synthetic profile %q", profileID)
	}
	if days <= 0 {
		days = profile.DefaultDays
	}
	if days < 3 {
		return nil, errors.New("at least three synthetic days are required")
	}

	records := make([]DailyRecord, 0, days)
	for i := 0; i < days; i++ {
		cycle := math.Sin(float64(i) * math.Pi / 3.5)
		loadPulse := 0.0
		if i%7 == 2 || i%7 == 3 {
			loadPulse = 1.0
		}
		stressPulse := 0.0
		if i%10 == 6 || i%10 == 7 {
			stressPulse = 1.0
		}

		sleepHours := 7.35 + (cycle * 0.35) - (loadPulse * 0.6) - (stressPulse * 0.55)
		stress := 4.0 + (stressPulse * 3.2) + (loadPulse * 1.0) - (cycle * 0.4)
		exercise := 28 + int(loadPulse*48) + int(math.Max(0, cycle)*18)
		trainingLoad := 32 + (loadPulse * 42) + (math.Max(0, cycle) * 18)
		fatigue := 3.0 + (loadPulse * 1.6) + (stressPulse * 1.2) - (cycle * 0.4)
		sleepScore := 78 + (cycle * 5) - (loadPulse * 9) - (stressPulse * 10)
		hrv := 53 + (cycle * 3.2) - (loadPulse * 4.2) - (stressPulse * 3.5)
		rhr := 60 - (cycle * 1.5) + (loadPulse * 2.5) + (stressPulse * 2.0)
		water := 2.35 - (stressPulse * 0.45) - (loadPulse * 0.2)
		caffeine := 1.4 + (stressPulse * 1.4)
		illness := 1.0

		switch profile.ID {
		case "high-strain":
			sleepHours -= 0.65
			stress += 1.5
			exercise += 12
			trainingLoad += 18
			fatigue += 1.2
			sleepScore -= 11
			hrv -= 5
			rhr += 3
			caffeine += 0.8
			water -= 0.2
			if i%9 == 4 {
				illness = 4
			}
		case "missing-sensors":
			trainingLoad = 0
			illness = 0
		}

		rec := DailyRecord{
			Date:                   start.AddDate(0, 0, i),
			SourceID:               profile.SourceID,
			SleepScore:             round1(clamp(sleepScore, 35, 96)),
			SleepHours:             round1(clamp(sleepHours, 4.5, 9.2)),
			HRV:                    round1(clamp(hrv, 28, 82)),
			RestingHeartRate:       round1(clamp(rhr, 48, 82)),
			Steps:                  6500 + (i%5)*650 + int(loadPulse*2800),
			ExerciseMinutes:        exercise,
			StressScore:            round1(clamp(stress, 1, 10)),
			FatigueScore:           round1(clamp(fatigue, 1, 10)),
			CaffeineServings:       round1(clamp(caffeine, 0, 5)),
			WaterLiters:            round1(clamp(water, 1.2, 3.4)),
			WeightKG:               round1(72.0 - float64(i)*0.02),
			BloodPressureSystolic:  int(math.Round(118 + loadPulse*4 + stressPulse*5)),
			BloodPressureDiastolic: int(math.Round(76 + loadPulse*3 + stressPulse*4)),
			BloodGlucoseMGDL:       round1(92 + loadPulse*8 + stressPulse*5),
			BodyTemperatureC:       round1(36.5 + illness*0.05),
			IllnessScore:           round1(illness),
			TrainingLoad:           round1(clamp(trainingLoad, 0, 95)),
			Notes:                  "synthetic fixture",
		}
		if profile.ID == "missing-sensors" {
			rec.BloodPressureSystolic = 0
			rec.BloodPressureDiastolic = 0
			rec.BloodGlucoseMGDL = 0
			rec.IllnessScore = 0
			rec.TrainingLoad = 0
		}
		records = append(records, rec)
	}
	return records, nil
}

// WriteDailyCSV serializes normalized daily records with the canonical header.
func WriteDailyCSV(w io.Writer, records []DailyRecord) error {
	cw := csv.NewWriter(w)
	if err := cw.Write(dailyCSVHeader()); err != nil {
		return err
	}
	for _, rec := range records {
		if err := cw.Write([]string{
			rec.Date.Format(DateLayout),
			rec.SourceID,
			formatFloat(rec.SleepScore),
			formatFloat(rec.SleepHours),
			formatFloat(rec.HRV),
			formatFloat(rec.RestingHeartRate),
			strconv.Itoa(rec.Steps),
			strconv.Itoa(rec.ExerciseMinutes),
			formatFloat(rec.StressScore),
			formatFloat(rec.FatigueScore),
			formatFloat(rec.CaffeineServings),
			formatFloat(rec.WaterLiters),
			formatFloat(rec.WeightKG),
			strconv.Itoa(rec.BloodPressureSystolic),
			strconv.Itoa(rec.BloodPressureDiastolic),
			formatFloat(rec.BloodGlucoseMGDL),
			formatFloat(rec.BodyTemperatureC),
			formatFloat(rec.IllnessScore),
			formatFloat(rec.TrainingLoad),
			rec.Notes,
		}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

// dailyCSVHeader returns the ordered internal import/export column contract.
func dailyCSVHeader() []string {
	return []string{
		"date",
		"source_id",
		"sleep_score",
		"sleep_hours",
		"hrv",
		"resting_heart_rate",
		"steps",
		"exercise_minutes",
		"stress_score",
		"fatigue_score",
		"caffeine_servings",
		"water_liters",
		"weight_kg",
		"blood_pressure_systolic",
		"blood_pressure_diastolic",
		"blood_glucose_mgdl",
		"body_temperature_c",
		"illness_score",
		"training_load",
		"notes",
	}
}

// findSyntheticProfile performs exact profile lookup.
func findSyntheticProfile(id string) (SyntheticProfile, bool) {
	for _, profile := range SyntheticProfiles() {
		if profile.ID == id {
			return profile, true
		}
	}
	return SyntheticProfile{}, false
}

// round1 rounds a value to one decimal place for stable fixtures.
func round1(v float64) float64 {
	return math.Round(v*10) / 10
}

// formatFloat emits one-decimal deterministic CSV values.
func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', 1, 64)
}
