package twin

type DeviceCapability string

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

type DeviceSource struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Kind         string             `json:"kind"`
	SyncMode     string             `json:"sync_mode"`
	PrivacyRisk  string             `json:"privacy_risk"`
	Capabilities []DeviceCapability `json:"capabilities"`
}

type DeviceAdapter interface {
	Source() DeviceSource
	Import(path string) ([]DailyRecord, error)
}

type CSVAdapter struct{}

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

func (CSVAdapter) Import(path string) ([]DailyRecord, error) {
	return LoadCSV(path)
}
