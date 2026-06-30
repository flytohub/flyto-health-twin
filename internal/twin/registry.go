package twin

type AdapterContract struct {
	Source               DeviceSource       `json:"source"`
	Status               string             `json:"status"`
	RawFields            []string           `json:"raw_fields"`
	AggregateFields      []string           `json:"aggregate_fields"`
	PublicFields         []string           `json:"public_fields"`
	PublicExportEligible bool               `json:"public_export_eligible"`
	RequiresConsent      bool               `json:"requires_consent"`
	GateChecks           []string           `json:"gate_checks"`
	Notes                []string           `json:"notes"`
	Capabilities         []DeviceCapability `json:"capabilities"`
}

type ModelCard struct {
	ID             string   `json:"id"`
	Version        string   `json:"version"`
	Status         string   `json:"status"`
	Purpose        string   `json:"purpose"`
	Inputs         []string `json:"inputs"`
	Outputs        []string `json:"outputs"`
	Assumptions    []string `json:"assumptions"`
	Limitations    []string `json:"limitations"`
	PublicEligible bool     `json:"public_eligible"`
}

type DatasetCandidate struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	License         string   `json:"license"`
	DataType        string   `json:"data_type"`
	IntendedUse     string   `json:"intended_use"`
	PrivacyConcerns []string `json:"privacy_concerns"`
	Limitations     []string `json:"limitations"`
	Status          string   `json:"status"`
}

func AdapterContracts() []AdapterContract {
	return []AdapterContract{
		{
			Source:               CSVAdapter{}.Source(),
			Status:               "implemented",
			RawFields:            dailyCSVHeader(),
			AggregateFields:      dailyCSVHeader(),
			PublicFields:         publicDailyFields(),
			PublicExportEligible: true,
			RequiresConsent:      false,
			GateChecks:           []string{"privacy header inspection", "public export redaction test", "synthetic fixture coverage"},
			Notes:                []string{"manual daily aggregate CSV is the default open-source adapter"},
			Capabilities:         CSVAdapter{}.Source().Capabilities,
		},
		{
			Source: DeviceSource{
				ID:          "wearable_daily_aggregate",
				Name:        "Wearable daily aggregate import",
				Kind:        "file_or_api",
				SyncMode:    "manual_first",
				PrivacyRisk: "medium_until_vendor_fields_are_mapped",
				Capabilities: []DeviceCapability{
					CapabilitySleep,
					CapabilityHRV,
					CapabilityRestingHeartRate,
					CapabilitySteps,
					CapabilityExercise,
				},
			},
			Status:               "contract_ready",
			RawFields:            []string{"device_user_id", "sleep_timeline", "heart_rate_timeseries", "activity_sessions"},
			AggregateFields:      []string{"sleep_score", "sleep_hours", "hrv", "resting_heart_rate", "steps", "exercise_minutes", "training_load"},
			PublicFields:         []string{"sleep_score", "sleep_hours", "hrv", "resting_heart_rate", "steps", "exercise_minutes"},
			PublicExportEligible: false,
			RequiresConsent:      true,
			GateChecks:           realAdapterGateChecks(),
			Notes:                []string{"raw timelines must aggregate before public export", "vendor credentials stay outside the repository"},
			Capabilities: []DeviceCapability{
				CapabilitySleep,
				CapabilityHRV,
				CapabilityRestingHeartRate,
				CapabilitySteps,
				CapabilityExercise,
			},
		},
		{
			Source: DeviceSource{
				ID:          "blood_pressure_monitor",
				Name:        "Blood pressure monitor summary import",
				Kind:        "file_or_api",
				SyncMode:    "manual_first",
				PrivacyRisk: "medium_private_by_default",
				Capabilities: []DeviceCapability{
					CapabilityBloodPressure,
				},
			},
			Status:               "contract_ready",
			RawFields:            []string{"session_timestamp", "systolic", "diastolic", "pulse"},
			AggregateFields:      []string{"blood_pressure_systolic", "blood_pressure_diastolic"},
			PublicFields:         []string{},
			PublicExportEligible: false,
			RequiresConsent:      true,
			GateChecks:           realAdapterGateChecks(),
			Notes:                []string{"session-level readings remain private unless explicitly aggregated and approved"},
			Capabilities: []DeviceCapability{
				CapabilityBloodPressure,
			},
		},
		{
			Source: DeviceSource{
				ID:          "lab_snapshot",
				Name:        "Periodic lab snapshot import",
				Kind:        "file",
				SyncMode:    "manual_review",
				PrivacyRisk: "high_private_by_default",
				Capabilities: []DeviceCapability{
					CapabilityLabSnapshot,
					CapabilityBloodGlucose,
				},
			},
			Status:               "gated_research_track",
			RawFields:            []string{"lab_panel_name", "collection_date", "result_value", "reference_range", "clinic_context"},
			AggregateFields:      []string{"blood_glucose_mgdl", "illness_score"},
			PublicFields:         []string{},
			PublicExportEligible: false,
			RequiresConsent:      true,
			GateChecks:           realAdapterGateChecks(),
			Notes:                []string{"full reports and clinical context must never enter public fixtures by default"},
			Capabilities: []DeviceCapability{
				CapabilityLabSnapshot,
				CapabilityBloodGlucose,
			},
		},
	}
}

func FindAdapterContract(id string) (AdapterContract, bool) {
	for _, contract := range AdapterContracts() {
		if contract.Source.ID == id {
			return contract, true
		}
	}
	return AdapterContract{}, false
}

func ModelRegistry() []ModelCard {
	model := BaselineModel{}
	return []ModelCard{
		{
			ID:      model.ID(),
			Version: model.Version(),
			Status:  "implemented",
			Purpose: "transparent next-day response baseline for error analysis",
			Inputs:  baselineFeatureSet(),
			Outputs: []string{
				"predicted_hrv",
				"predicted_resting_heart_rate",
				"predicted_fatigue_score",
				"predicted_sleep_score",
				"recovery_state",
				"missing_variables",
			},
			Assumptions:    baselineAssumptions(),
			Limitations:    []string{"not clinical", "not calibrated to a real device yet", "uses short recent windows"},
			PublicEligible: true,
		},
		{
			ID:             "weighted_trend",
			Version:        "planned",
			Status:         "planned",
			Purpose:        "compare recency-weighted trend features against the transparent baseline",
			Inputs:         baselineFeatureSet(),
			Outputs:        []string{"same response metrics as baseline_strain"},
			Assumptions:    []string{"recent days carry higher weight than older daily aggregates"},
			Limitations:    []string{"not implemented", "requires benchmark comparison before dashboard exposure"},
			PublicEligible: false,
		},
		{
			ID:             "equipment_calibrated",
			Version:        "blocked_until_device_gate",
			Status:         "gated",
			Purpose:        "future adapter-calibrated model after real equipment samples exist",
			Inputs:         []string{"daily aggregates", "validated device adapter outputs", "approved lab summary features"},
			Outputs:        []string{"personal response predictions with provenance"},
			Assumptions:    []string{"device schemas and privacy filters are validated first"},
			Limitations:    []string{"must not be enabled without the equipment integration gate"},
			PublicEligible: false,
		},
	}
}

func DatasetRegistry() []DatasetCandidate {
	return []DatasetCandidate{
		{
			ID:          "synthetic_daily_v0",
			Name:        "Flyto2 deterministic daily aggregate fixture",
			License:     "Apache-2.0 repository fixture",
			DataType:    "synthetic daily wearable and lifestyle aggregates",
			IntendedUse: "local development, model regression checks, public dashboard demos",
			PrivacyConcerns: []string{
				"contains no real person",
				"must not be treated as cohort evidence",
			},
			Limitations: []string{"small deterministic fixture", "not population representative"},
			Status:      "implemented",
		},
		{
			ID:          "public_wearable_research",
			Name:        "Public wearable datasets in scientific repositories",
			License:     "dataset-specific",
			DataType:    "wearable time series or daily summaries",
			IntendedUse: "future importer and benchmark inspiration",
			PrivacyConcerns: []string{
				"public datasets can still contain re-identification risk",
				"licenses must be checked before redistribution",
			},
			Limitations: []string{"different devices", "different populations", "not personal calibration"},
			Status:      "tracked_not_imported",
		},
		{
			ID:          "biology_literature_context",
			Name:        "Biology literature and public omics portals",
			License:     "source-specific",
			DataType:    "population or cellular research context",
			IntendedUse: "safe toy simulation and research notes only",
			PrivacyConcerns: []string{
				"must not be mixed into personal predictions as if personal evidence",
			},
			Limitations: []string{"not a full-body digital twin", "not intervention evidence"},
			Status:      "docs_only",
		},
	}
}

func realAdapterGateChecks() []string {
	return []string{
		"sample export or API documentation is available",
		"raw private fields are documented",
		"synthetic fixture exists",
		"privacy checker covers blocked raw fields",
		"importer tests pass",
		"public export redaction is proven",
	}
}

func publicDailyFields() []string {
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
	}
}
