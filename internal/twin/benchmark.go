package twin

import "fmt"

type BenchmarkFixture struct {
	ID              string       `json:"id"`
	ProfileID       string       `json:"profile_id"`
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	RecordCount     int          `json:"record_count"`
	PrivacyProfile  string       `json:"privacy_profile"`
	PassThresholds  ErrorSummary `json:"pass_thresholds"`
	ExpectedSignals []string     `json:"expected_signals"`
}

type ModelEvaluationReport struct {
	DatasetID             string         `json:"dataset_id"`
	ModelID               string         `json:"model_id"`
	ModelVersion          string         `json:"model_version"`
	RecordCount           int            `json:"record_count"`
	EvaluationCount       int            `json:"evaluation_count"`
	ErrorSummary          ErrorSummary   `json:"error_summary"`
	PassThresholds        ErrorSummary   `json:"pass_thresholds"`
	Pass                  bool           `json:"pass"`
	RecoveryStateCounts   map[string]int `json:"recovery_state_counts"`
	MissingVariableCounts map[string]int `json:"missing_variable_counts"`
	Notes                 []string       `json:"caveats"`
}

func BenchmarkFixtures() []BenchmarkFixture {
	return []BenchmarkFixture{
		{
			ID:             "synthetic_balanced_v0",
			ProfileID:      "balanced",
			Name:           "Balanced synthetic benchmark",
			Description:    "Deterministic daily aggregate fixture with recovery, strain, and rebound periods.",
			RecordCount:    30,
			PrivacyProfile: "synthetic_daily_aggregate",
			PassThresholds: benchmarkThresholds(),
			ExpectedSignals: []string{
				"strained state after high load or stress pulses",
				"recovered state after sleep and low-load rebound",
				"missing-variable hints for optional lab-style fields",
			},
		},
		{
			ID:             "synthetic_high_strain_v0",
			ProfileID:      "high-strain",
			Name:           "High-strain synthetic benchmark",
			Description:    "Deterministic fixture for stress-testing strained predictions and error reports.",
			RecordCount:    30,
			PrivacyProfile: "synthetic_daily_aggregate",
			PassThresholds: benchmarkThresholds(),
			ExpectedSignals: []string{
				"frequent strained states",
				"higher resting heart rate predictions",
				"sleep and HRV degradation under load",
			},
		},
	}
}

func BuildBenchmarkReport(profileID string, records []DailyRecord) (ModelEvaluationReport, error) {
	fixture, ok := benchmarkFixtureForProfile(profileID)
	if !ok {
		return ModelEvaluationReport{}, fmt.Errorf("unknown benchmark profile %q", profileID)
	}
	return BuildModelEvaluationReport(BaselineModel{}, records, fixture.ID)
}

func BuildModelEvaluationReport(model Model, records []DailyRecord, datasetID string) (ModelEvaluationReport, error) {
	evals, err := EvaluateWithModel(model, records)
	if err != nil {
		return ModelEvaluationReport{}, err
	}
	return BuildModelEvaluationReportFromEvaluations(datasetID, model.ID(), model.Version(), len(records), evals), nil
}

func BuildModelEvaluationReportFromEvaluations(datasetID, modelID, modelVersion string, recordCount int, evals []Evaluation) ModelEvaluationReport {
	hrv, rhr, fatigue, sleep := MeanAbsoluteErrors(evals)
	report := ModelEvaluationReport{
		DatasetID:       datasetID,
		ModelID:         modelID,
		ModelVersion:    modelVersion,
		RecordCount:     recordCount,
		EvaluationCount: len(evals),
		ErrorSummary: ErrorSummary{
			HRVMeanAbsoluteError:     hrv,
			RHRMeanAbsoluteError:     rhr,
			FatigueMeanAbsoluteError: fatigue,
			SleepMeanAbsoluteError:   sleep,
		},
		PassThresholds:        benchmarkThresholds(),
		RecoveryStateCounts:   map[string]int{},
		MissingVariableCounts: map[string]int{},
		Notes: []string{
			"benchmark uses synthetic daily aggregates only",
			"pass/fail is for regression guardrails, not clinical accuracy",
		},
	}
	for _, eval := range evals {
		report.RecoveryStateCounts[eval.Prediction.RecoveryState]++
		for _, missing := range eval.Prediction.MissingVariables {
			report.MissingVariableCounts[missing]++
		}
	}
	report.Pass = report.ErrorSummary.HRVMeanAbsoluteError <= report.PassThresholds.HRVMeanAbsoluteError &&
		report.ErrorSummary.RHRMeanAbsoluteError <= report.PassThresholds.RHRMeanAbsoluteError &&
		report.ErrorSummary.FatigueMeanAbsoluteError <= report.PassThresholds.FatigueMeanAbsoluteError &&
		report.ErrorSummary.SleepMeanAbsoluteError <= report.PassThresholds.SleepMeanAbsoluteError
	return report
}

func benchmarkFixtureForProfile(profileID string) (BenchmarkFixture, bool) {
	for _, fixture := range BenchmarkFixtures() {
		if fixture.ProfileID == profileID {
			return fixture, true
		}
	}
	return BenchmarkFixture{}, false
}

func benchmarkThresholds() ErrorSummary {
	return ErrorSummary{
		HRVMeanAbsoluteError:     9,
		RHRMeanAbsoluteError:     6,
		FatigueMeanAbsoluteError: 2.5,
		SleepMeanAbsoluteError:   18,
	}
}
