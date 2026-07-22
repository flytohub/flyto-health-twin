package twin

import "time"

// DateLayout is the canonical local-day representation used by CSV and JSON.
const DateLayout = "2006-01-02"

// DailyRecord is the normalized internal aggregate for one local calendar day.
// It may contain private optional values and must not be serialized publicly.
type DailyRecord struct {
	Date                   time.Time
	SourceID               string
	SleepScore             float64
	SleepHours             float64
	HRV                    float64
	RestingHeartRate       float64
	Steps                  int
	ExerciseMinutes        int
	StressScore            float64
	FatigueScore           float64
	CaffeineServings       float64
	WaterLiters            float64
	WeightKG               float64
	BloodPressureSystolic  int
	BloodPressureDiastolic int
	BloodGlucoseMGDL       float64
	BodyTemperatureC       float64
	IllnessScore           float64
	TrainingLoad           float64
	Notes                  string
}

// Prediction records one transparent next-day model result and its provenance.
type Prediction struct {
	TargetDate                time.Time
	ModelID                   string
	ModelVersion              string
	InputStartDate            time.Time
	InputEndDate              time.Time
	FeatureSet                []string
	Assumptions               []string
	PredictedHRV              float64
	PredictedRestingHeartRate float64
	PredictedFatigueScore     float64
	PredictedSleepScore       float64
	RecoveryState             string
	Hints                     []string
	MissingVariables          []string
}

// Evaluation compares one prediction with the actual target-day record.
type Evaluation struct {
	TargetDate   time.Time
	Prediction   Prediction
	Actual       DailyRecord
	HRVError     float64
	RHRError     float64
	FatigueError float64
	SleepError   float64
}

// ErrorSummary contains mean absolute error for each modeled response metric.
type ErrorSummary struct {
	HRVMeanAbsoluteError     float64 `json:"hrv_mean_absolute_error"`
	RHRMeanAbsoluteError     float64 `json:"rhr_mean_absolute_error"`
	FatigueMeanAbsoluteError float64 `json:"fatigue_mean_absolute_error"`
	SleepMeanAbsoluteError   float64 `json:"sleep_mean_absolute_error"`
}
