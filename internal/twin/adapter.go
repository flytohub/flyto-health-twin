package twin

// DeviceCapability identifies one normalized signal an adapter can provide.
type DeviceCapability string

// Supported device capabilities used by adapter contracts and gates.
const (
	CapabilitySleep            DeviceCapability = "sleep"
	CapabilityHRV              DeviceCapability = "hrv"
	CapabilityRestingHeartRate DeviceCapability = "resting_heart_rate"
	CapabilitySteps            DeviceCapability = "steps"
	CapabilityExercise         DeviceCapability = "exercise"
	CapabilityBloodPressure    DeviceCapability = "blood_pressure"
	CapabilityBloodGlucose     DeviceCapability = "blood_glucose"
	CapabilityBodyWeight       DeviceCapability = "body_weight"
	CapabilityManualSurvey     DeviceCapability = "manual_survey"
	CapabilityLabSnapshot      DeviceCapability = "lab_snapshot"
)

// DeviceSource describes the provenance, transport, risk, and signals of input.
type DeviceSource struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Kind         string             `json:"kind"`
	SyncMode     string             `json:"sync_mode"`
	PrivacyRisk  string             `json:"privacy_risk"`
	Capabilities []DeviceCapability `json:"capabilities"`
}

// DeviceAdapter imports one source into normalized daily records.
type DeviceAdapter interface {
	Source() DeviceSource
	Import(path string) ([]DailyRecord, error)
}

// CSVAdapter implements the local manual daily-aggregate CSV source.
type CSVAdapter struct{}

// Source returns the manual CSV source contract.
func (CSVAdapter) Source() DeviceSource {
	return DeviceSource{
		ID:          "manual_csv",
		Name:        "Manual daily aggregate CSV",
		Kind:        "file",
		SyncMode:    "manual",
		PrivacyRisk: "low_when_daily_aggregate",
		Capabilities: []DeviceCapability{
			CapabilitySleep,
			CapabilityHRV,
			CapabilityRestingHeartRate,
			CapabilitySteps,
			CapabilityExercise,
			CapabilityBloodPressure,
			CapabilityBloodGlucose,
			CapabilityBodyWeight,
			CapabilityManualSurvey,
		},
	}
}

// Import loads daily records from the supplied CSV path.
func (CSVAdapter) Import(path string) ([]DailyRecord, error) {
	return LoadCSV(path)
}
