package twin

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// LoadCSV reads daily aggregate records from a CSV file and returns them sorted
// by date.
func LoadCSV(path string) (records []DailyRecord, err error) {
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

	index := map[string]int{}
	for i, name := range header {
		normalized := strings.ToLower(strings.TrimSpace(name))
		if normalized == "" {
			return nil, fmt.Errorf("header column %d is empty", i+1)
		}
		if _, exists := index[normalized]; exists {
			return nil, fmt.Errorf("duplicate header %q", normalized)
		}
		index[normalized] = i
	}
	if _, ok := index["date"]; !ok {
		return nil, fmt.Errorf("required header %q is missing", "date")
	}

	for rowNum := 2; ; rowNum++ {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read row %d: %w", rowNum, err)
		}
		rec, err := parseRow(index, row)
		if err != nil {
			return nil, fmt.Errorf("parse row %d: %w", rowNum, err)
		}
		records = append(records, rec)
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Date.Before(records[j].Date)
	})
	for i := 1; i < len(records); i++ {
		if records[i].Date.Equal(records[i-1].Date) {
			return nil, fmt.Errorf("duplicate daily record for %s", records[i].Date.Format(DateLayout))
		}
	}
	return records, nil
}

// parseRow converts one indexed CSV row into the normalized daily contract.
func parseRow(index map[string]int, row []string) (DailyRecord, error) {
	get := func(name string) string {
		i, ok := index[name]
		if !ok || i >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[i])
	}
	parseFloat := func(name string) (float64, error) {
		raw := get(name)
		if raw == "" {
			return 0, nil
		}
		value, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return 0, err
		}
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return 0, fmt.Errorf("must be finite")
		}
		return value, nil
	}
	parseInt := func(name string) (int, error) {
		raw := get(name)
		if raw == "" {
			return 0, nil
		}
		return strconv.Atoi(raw)
	}

	date, err := time.Parse(DateLayout, get("date"))
	if err != nil {
		return DailyRecord{}, fmt.Errorf("date: %w", err)
	}
	sleepScore, err := parseFloat("sleep_score")
	if err != nil {
		return DailyRecord{}, fmt.Errorf("sleep_score: %w", err)
	}
	sleepHours, err := parseFloat("sleep_hours")
	if err != nil {
		return DailyRecord{}, fmt.Errorf("sleep_hours: %w", err)
	}
	hrv, err := parseFloat("hrv")
	if err != nil {
		return DailyRecord{}, fmt.Errorf("hrv: %w", err)
	}
	rhr, err := parseFloat("resting_heart_rate")
	if err != nil {
		return DailyRecord{}, fmt.Errorf("resting_heart_rate: %w", err)
	}
	steps, err := parseInt("steps")
	if err != nil {
		return DailyRecord{}, fmt.Errorf("steps: %w", err)
	}
	exercise, err := parseInt("exercise_minutes")
	if err != nil {
		return DailyRecord{}, fmt.Errorf("exercise_minutes: %w", err)
	}
	stress, err := parseFloat("stress_score")
	if err != nil {
		return DailyRecord{}, fmt.Errorf("stress_score: %w", err)
	}
	fatigue, err := parseFloat("fatigue_score")
	if err != nil {
		return DailyRecord{}, fmt.Errorf("fatigue_score: %w", err)
	}
	caffeine, err := parseFloat("caffeine_servings")
	if err != nil {
		return DailyRecord{}, fmt.Errorf("caffeine_servings: %w", err)
	}
	water, err := parseFloat("water_liters")
	if err != nil {
		return DailyRecord{}, fmt.Errorf("water_liters: %w", err)
	}
	weight, err := parseFloat("weight_kg")
	if err != nil {
		return DailyRecord{}, fmt.Errorf("weight_kg: %w", err)
	}
	bpSys, err := parseInt("blood_pressure_systolic")
	if err != nil {
		return DailyRecord{}, fmt.Errorf("blood_pressure_systolic: %w", err)
	}
	bpDia, err := parseInt("blood_pressure_diastolic")
	if err != nil {
		return DailyRecord{}, fmt.Errorf("blood_pressure_diastolic: %w", err)
	}
	glucose, err := parseFloat("blood_glucose_mgdl")
	if err != nil {
		return DailyRecord{}, fmt.Errorf("blood_glucose_mgdl: %w", err)
	}
	temp, err := parseFloat("body_temperature_c")
	if err != nil {
		return DailyRecord{}, fmt.Errorf("body_temperature_c: %w", err)
	}
	illness, err := parseFloat("illness_score")
	if err != nil {
		return DailyRecord{}, fmt.Errorf("illness_score: %w", err)
	}
	trainingLoad, err := parseFloat("training_load")
	if err != nil {
		return DailyRecord{}, fmt.Errorf("training_load: %w", err)
	}

	return DailyRecord{
		Date:                   date,
		SourceID:               get("source_id"),
		SleepScore:             sleepScore,
		SleepHours:             sleepHours,
		HRV:                    hrv,
		RestingHeartRate:       rhr,
		Steps:                  steps,
		ExerciseMinutes:        exercise,
		StressScore:            stress,
		FatigueScore:           fatigue,
		CaffeineServings:       caffeine,
		WaterLiters:            water,
		WeightKG:               weight,
		BloodPressureSystolic:  bpSys,
		BloodPressureDiastolic: bpDia,
		BloodGlucoseMGDL:       glucose,
		BodyTemperatureC:       temp,
		IllnessScore:           illness,
		TrainingLoad:           trainingLoad,
		Notes:                  get("notes"),
	}, nil
}
