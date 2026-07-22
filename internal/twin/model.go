package twin

import (
	"errors"
	"math"
)

// Model defines a versioned next-day prediction implementation.
type Model interface {
	ID() string
	Version() string
	Predict(history []DailyRecord) (Prediction, error)
}

// BaselineModel is the deterministic transparent strain heuristic.
type BaselineModel struct{}

// ID returns the stable baseline model identifier.
func (BaselineModel) ID() string {
	return "baseline_strain"
}

// Version returns the baseline model contract version.
func (BaselineModel) Version() string {
	return "0.1.0"
}

// PredictNext generates a transparent next-day baseline prediction from recent
// daily records.
func PredictNext(records []DailyRecord) (Prediction, error) {
	return BaselineModel{}.Predict(records)
}

// Predict derives one next-day result from strictly ordered daily history.
func (m BaselineModel) Predict(records []DailyRecord) (Prediction, error) {
	if len(records) < 2 {
		return Prediction{}, errors.New("at least two records are required")
	}
	if err := validateChronology(records); err != nil {
		return Prediction{}, err
	}

	recent := tail(records, 3)
	last := records[len(records)-1]

	avgHRV := avg(recent, func(r DailyRecord) float64 { return r.HRV })
	avgRHR := avg(recent, func(r DailyRecord) float64 { return r.RestingHeartRate })
	avgFatigue := avg(recent, func(r DailyRecord) float64 { return r.FatigueScore })
	avgSleep := avg(recent, func(r DailyRecord) float64 { return r.SleepScore })

	strain := 0.0
	if last.SleepHours < 6.5 {
		strain += 1.0
	}
	if last.StressScore >= 7 {
		strain += 1.0
	}
	if last.ExerciseMinutes >= 70 {
		strain += 0.8
	}
	if last.CaffeineServings >= 3 {
		strain += 0.5
	}
	if last.WaterLiters > 0 && last.WaterLiters < 2 {
		strain += 0.4
	}
	if last.IllnessScore >= 5 {
		strain += 1.2
	}
	if last.TrainingLoad >= 75 {
		strain += 0.8
	}

	p := Prediction{
		ModelID:                   m.ID(),
		ModelVersion:              m.Version(),
		TargetDate:                last.Date.AddDate(0, 0, 1),
		InputStartDate:            records[0].Date,
		InputEndDate:              last.Date,
		FeatureSet:                baselineFeatureSet(),
		Assumptions:               baselineAssumptions(),
		PredictedHRV:              clamp(avgHRV-(strain*2.2), 20, 120),
		PredictedRestingHeartRate: clamp(avgRHR+(strain*1.4), 35, 120),
		PredictedFatigueScore:     clamp(avgFatigue+(strain*0.7), 1, 10),
		PredictedSleepScore:       clamp(avgSleep-(strain*4.0), 1, 100),
		Hints:                     explainStrain(last, strain),
		MissingVariables:          missingVariableHints(last),
	}
	p.RecoveryState = recoveryState(p)
	return p, nil
}

// Evaluate replays predictions through the time series and compares each
// prediction with the next actual daily record.
func Evaluate(records []DailyRecord) ([]Evaluation, error) {
	return EvaluateWithModel(BaselineModel{}, records)
}

// EvaluateWithModel runs a supplied model over strictly ordered rolling history.
func EvaluateWithModel(model Model, records []DailyRecord) ([]Evaluation, error) {
	if len(records) < 3 {
		return nil, errors.New("at least three records are required")
	}
	if err := validateChronology(records); err != nil {
		return nil, err
	}

	var evaluations []Evaluation
	for i := 2; i < len(records); i++ {
		prediction, err := model.Predict(records[:i])
		if err != nil {
			return nil, err
		}
		actual := records[i]
		evaluations = append(evaluations, Evaluation{
			TargetDate:   actual.Date,
			Prediction:   prediction,
			Actual:       actual,
			HRVError:     actual.HRV - prediction.PredictedHRV,
			RHRError:     actual.RestingHeartRate - prediction.PredictedRestingHeartRate,
			FatigueError: actual.FatigueScore - prediction.PredictedFatigueScore,
			SleepError:   actual.SleepScore - prediction.PredictedSleepScore,
		})
	}
	return evaluations, nil
}

// MeanAbsoluteErrors summarizes prediction error across evaluation rows.
func MeanAbsoluteErrors(evals []Evaluation) (hrv, rhr, fatigue, sleep float64) {
	if len(evals) == 0 {
		return 0, 0, 0, 0
	}
	for _, e := range evals {
		hrv += math.Abs(e.HRVError)
		rhr += math.Abs(e.RHRError)
		fatigue += math.Abs(e.FatigueError)
		sleep += math.Abs(e.SleepError)
	}
	n := float64(len(evals))
	return hrv / n, rhr / n, fatigue / n, sleep / n
}

// tail returns at most the final n records without copying them.
func tail(records []DailyRecord, n int) []DailyRecord {
	if len(records) <= n {
		return records
	}
	return records[len(records)-n:]
}

// avg computes the mean of one projected record value.
func avg(records []DailyRecord, value func(DailyRecord) float64) float64 {
	if len(records) == 0 {
		return 0
	}
	total := 0.0
	for _, r := range records {
		total += value(r)
	}
	return total / float64(len(records))
}

// clamp bounds a floating-point value to an inclusive interval.
func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// recoveryState classifies a prediction as strained, recovered, or normal.
func recoveryState(p Prediction) string {
	if p.PredictedFatigueScore >= 6 || p.PredictedHRV < 45 || p.PredictedRestingHeartRate > 65 {
		return "strained"
	}
	if p.PredictedFatigueScore <= 3 && p.PredictedHRV >= 52 {
		return "recovered"
	}
	return "normal"
}

// explainStrain turns active heuristic conditions into human-readable hints.
func explainStrain(last DailyRecord, strain float64) []string {
	if strain == 0 {
		return []string{"baseline trend; no major strain marker in current inputs"}
	}

	var hints []string
	if last.SleepHours < 6.5 {
		hints = append(hints, "short sleep may reduce next-day HRV")
	}
	if last.StressScore >= 7 {
		hints = append(hints, "high stress may raise resting heart rate and fatigue")
	}
	if last.ExerciseMinutes >= 70 {
		hints = append(hints, "high exercise load may require more recovery time")
	}
	if last.CaffeineServings >= 3 {
		hints = append(hints, "high caffeine may affect sleep quality")
	}
	if last.WaterLiters > 0 && last.WaterLiters < 2 {
		hints = append(hints, "low water intake may be a missing recovery variable")
	}
	if last.IllnessScore >= 5 {
		hints = append(hints, "illness marker may dominate wearable recovery signals")
	}
	if last.TrainingLoad >= 75 {
		hints = append(hints, "high training load may require more recovery time")
	}
	return hints
}

// missingVariableHints lists optional signals absent from the latest record.
func missingVariableHints(last DailyRecord) []string {
	var hints []string
	if last.StressScore == 0 {
		hints = append(hints, "stress_score")
	}
	if last.FatigueScore == 0 {
		hints = append(hints, "fatigue_score")
	}
	if last.IllnessScore == 0 {
		hints = append(hints, "illness_score")
	}
	if last.TrainingLoad == 0 && last.ExerciseMinutes > 0 {
		hints = append(hints, "training_load")
	}
	if last.BloodGlucoseMGDL == 0 {
		hints = append(hints, "blood_glucose_mgdl")
	}
	if last.BloodPressureSystolic == 0 || last.BloodPressureDiastolic == 0 {
		hints = append(hints, "blood_pressure")
	}
	if len(hints) == 0 {
		return []string{"no obvious missing variable from current schema"}
	}
	return hints
}

// baselineFeatureSet returns the model-card input field names.
func baselineFeatureSet() []string {
	return []string{
		"sleep_score",
		"sleep_hours",
		"hrv",
		"resting_heart_rate",
		"exercise_minutes",
		"stress_score",
		"fatigue_score",
		"caffeine_servings",
		"water_liters",
		"illness_score",
		"training_load",
	}
}

// baselineAssumptions returns the non-clinical heuristic assumptions.
func baselineAssumptions() []string {
	return []string{
		"recent daily aggregates approximate short-term recovery trend",
		"short sleep, high stress, high load, illness, and low hydration increase strain",
		"model is a transparent baseline for error analysis, not clinical prediction",
	}
}

// validateChronology requires unique records in strictly increasing date order.
func validateChronology(records []DailyRecord) error {
	for i := 1; i < len(records); i++ {
		if !records[i].Date.After(records[i-1].Date) {
			return errors.New("records must be in strictly increasing date order")
		}
	}
	return nil
}
