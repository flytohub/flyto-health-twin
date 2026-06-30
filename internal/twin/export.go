package twin

import (
	"encoding/json"
	"io"
	"time"
)

type PublicExport struct {
	Project       string                `json:"project"`
	GeneratedAt   string                `json:"generated_at"`
	Records       []PublicRecordJSON    `json:"records"`
	Evaluations   []PublicEvaluation    `json:"evaluations"`
	ErrorSummary  ErrorSummary          `json:"error_summary"`
	RoadmapStatus PublicRoadmapStatus   `json:"roadmap_status"`
	Benchmark     ModelEvaluationReport `json:"benchmark"`
}

type PublicRoadmapStatus struct {
	AdapterContractCount  int                         `json:"adapter_contract_count"`
	ModelCardCount        int                         `json:"model_card_count"`
	DatasetCandidateCount int                         `json:"dataset_candidate_count"`
	WorkflowRecipeCount   int                         `json:"workflow_recipe_count"`
	EquipmentGates        []PublicEquipmentGateStatus `json:"equipment_gates"`
	SimulationBoundary    string                      `json:"simulation_boundary"`
}

type PublicEquipmentGateStatus struct {
	SourceID string `json:"source_id"`
	Status   string `json:"status"`
}

type PublicRecordJSON struct {
	Date             string  `json:"date"`
	SourceID         string  `json:"source_id,omitempty"`
	SleepScore       float64 `json:"sleep_score"`
	SleepHours       float64 `json:"sleep_hours"`
	HRV              float64 `json:"hrv"`
	RestingHeartRate float64 `json:"resting_heart_rate"`
	Steps            int     `json:"steps"`
	ExerciseMinutes  int     `json:"exercise_minutes"`
	StressScore      float64 `json:"stress_score"`
	FatigueScore     float64 `json:"fatigue_score"`
	CaffeineServings float64 `json:"caffeine_servings"`
	WaterLiters      float64 `json:"water_liters"`
}

type PublicPrediction struct {
	TargetDate                string   `json:"target_date"`
	ModelID                   string   `json:"model_id"`
	ModelVersion              string   `json:"model_version"`
	InputStartDate            string   `json:"input_start_date"`
	InputEndDate              string   `json:"input_end_date"`
	FeatureSet                []string `json:"feature_set"`
	Assumptions               []string `json:"assumptions"`
	PredictedHRV              float64  `json:"predicted_hrv"`
	PredictedRestingHeartRate float64  `json:"predicted_resting_heart_rate"`
	PredictedFatigueScore     float64  `json:"predicted_fatigue_score"`
	PredictedSleepScore       float64  `json:"predicted_sleep_score"`
	RecoveryState             string   `json:"recovery_state"`
	Hints                     []string `json:"hints"`
	MissingVariables          []string `json:"missing_variables"`
}

type PublicEvaluation struct {
	TargetDate   string           `json:"target_date"`
	Prediction   PublicPrediction `json:"prediction"`
	Actual       PublicRecordJSON `json:"actual"`
	HRVError     float64          `json:"hrv_error"`
	RHRError     float64          `json:"rhr_error"`
	FatigueError float64          `json:"fatigue_error"`
	SleepError   float64          `json:"sleep_error"`
}

func BuildPublicExport(records []DailyRecord, evals []Evaluation) PublicExport {
	return BuildPublicExportAt(records, evals, time.Now().UTC())
}

func BuildPublicExportAt(records []DailyRecord, evals []Evaluation, generatedAt time.Time) PublicExport {
	publicRecords := make([]PublicRecordJSON, 0, len(records))
	for _, r := range records {
		publicRecords = append(publicRecords, toPublicRecordJSON(PublicRecord(r)))
	}

	publicEvals := make([]PublicEvaluation, 0, len(evals))
	for _, e := range evals {
		publicEvals = append(publicEvals, PublicEvaluation{
			TargetDate:   e.TargetDate.Format(DateLayout),
			Prediction:   toPublicPrediction(e.Prediction),
			Actual:       toPublicRecordJSON(PublicRecord(e.Actual)),
			HRVError:     e.HRVError,
			RHRError:     e.RHRError,
			FatigueError: e.FatigueError,
			SleepError:   e.SleepError,
		})
	}

	hrv, rhr, fatigue, sleep := MeanAbsoluteErrors(evals)
	errorSummary := ErrorSummary{
		HRVMeanAbsoluteError:     hrv,
		RHRMeanAbsoluteError:     rhr,
		FatigueMeanAbsoluteError: fatigue,
		SleepMeanAbsoluteError:   sleep,
	}
	baseline := BaselineModel{}
	return PublicExport{
		Project:      "flyto2",
		GeneratedAt:  generatedAt.UTC().Format(time.RFC3339),
		Records:      publicRecords,
		Evaluations:  publicEvals,
		ErrorSummary: errorSummary,
		RoadmapStatus: PublicRoadmapStatus{
			AdapterContractCount:  len(AdapterContracts()),
			ModelCardCount:        len(ModelRegistry()),
			DatasetCandidateCount: len(DatasetRegistry()),
			WorkflowRecipeCount:   len(WorkflowRecipes()),
			EquipmentGates:        publicEquipmentGateStatuses(),
			SimulationBoundary:    "safe educational toy simulations only; not personal biology or medical guidance",
		},
		Benchmark: BuildModelEvaluationReportFromEvaluations(
			"synthetic_daily_v0",
			baseline.ID(),
			baseline.Version(),
			len(records),
			evals,
		),
	}
}

func publicEquipmentGateStatuses() []PublicEquipmentGateStatus {
	reports := CheckAllEquipmentGates()
	statuses := make([]PublicEquipmentGateStatus, 0, len(reports))
	for _, report := range reports {
		statuses = append(statuses, PublicEquipmentGateStatus{
			SourceID: report.SourceID,
			Status:   report.Status,
		})
	}
	return statuses
}

func WritePublicExport(w io.Writer, records []DailyRecord, evals []Evaluation) error {
	return WritePublicExportAt(w, records, evals, time.Now().UTC())
}

func WritePublicExportAt(w io.Writer, records []DailyRecord, evals []Evaluation, generatedAt time.Time) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(BuildPublicExportAt(records, evals, generatedAt))
}

func toPublicRecordJSON(r DailyRecord) PublicRecordJSON {
	return PublicRecordJSON{
		Date:             r.Date.Format(DateLayout),
		SourceID:         r.SourceID,
		SleepScore:       r.SleepScore,
		SleepHours:       r.SleepHours,
		HRV:              r.HRV,
		RestingHeartRate: r.RestingHeartRate,
		Steps:            r.Steps,
		ExerciseMinutes:  r.ExerciseMinutes,
		StressScore:      r.StressScore,
		FatigueScore:     r.FatigueScore,
		CaffeineServings: r.CaffeineServings,
		WaterLiters:      r.WaterLiters,
	}
}

func toPublicPrediction(p Prediction) PublicPrediction {
	return PublicPrediction{
		TargetDate:                p.TargetDate.Format(DateLayout),
		ModelID:                   p.ModelID,
		ModelVersion:              p.ModelVersion,
		InputStartDate:            p.InputStartDate.Format(DateLayout),
		InputEndDate:              p.InputEndDate.Format(DateLayout),
		FeatureSet:                p.FeatureSet,
		Assumptions:               p.Assumptions,
		PredictedHRV:              p.PredictedHRV,
		PredictedRestingHeartRate: p.PredictedRestingHeartRate,
		PredictedFatigueScore:     p.PredictedFatigueScore,
		PredictedSleepScore:       p.PredictedSleepScore,
		RecoveryState:             p.RecoveryState,
		Hints:                     p.Hints,
		MissingVariables:          p.MissingVariables,
	}
}
