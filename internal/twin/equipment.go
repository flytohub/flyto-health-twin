package twin

import "fmt"

type GateCheck struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Required bool   `json:"required"`
	Passed   bool   `json:"passed"`
	Evidence string `json:"evidence"`
}

type EquipmentGateReport struct {
	SourceID    string      `json:"source_id"`
	SourceName  string      `json:"source_name"`
	Status      string      `json:"status"`
	Checks      []GateCheck `json:"checks"`
	NextActions []string    `json:"next_actions"`
}

func CheckEquipmentGate(sourceID string) (EquipmentGateReport, error) {
	contract, ok := FindAdapterContract(sourceID)
	if !ok {
		return EquipmentGateReport{}, fmt.Errorf("unknown adapter source %q", sourceID)
	}

	if contract.Source.ID == "manual_csv" {
		return EquipmentGateReport{
			SourceID:   contract.Source.ID,
			SourceName: contract.Source.Name,
			Status:     "ready_for_synthetic_public_demo",
			Checks: []GateCheck{
				{ID: "contract", Label: "adapter contract exists", Required: true, Passed: true, Evidence: "CSVAdapter.Source"},
				{ID: "fixture", Label: "synthetic fixture exists", Required: true, Passed: true, Evidence: "examples/synthetic_daily.csv"},
				{ID: "privacy", Label: "privacy checker covers blocked headers", Required: true, Passed: true, Evidence: "InspectCSVPrivacy tests"},
				{ID: "export", Label: "public export omits private fields", Required: true, Passed: true, Evidence: "PublicRecord and export tests"},
			},
			NextActions: []string{"keep using manual_csv until a real device sample is reviewed"},
		}, nil
	}

	checks := []GateCheck{
		{ID: "sample", Label: "sample export or API documentation is available", Required: true, Passed: false, Evidence: "not provided"},
		{ID: "raw_fields", Label: "raw private fields are documented", Required: true, Passed: true, Evidence: "contract raw_fields"},
		{ID: "fixture", Label: "synthetic fixture exists", Required: true, Passed: false, Evidence: "device-specific fixture not committed"},
		{ID: "privacy", Label: "privacy checker covers blocked raw fields", Required: true, Passed: false, Evidence: "pending importer-specific tests"},
		{ID: "importer", Label: "importer tests pass", Required: true, Passed: false, Evidence: "real importer not implemented"},
		{ID: "export", Label: "public export redaction is proven", Required: true, Passed: false, Evidence: "pending fixture"},
	}
	return EquipmentGateReport{
		SourceID:   contract.Source.ID,
		SourceName: contract.Source.Name,
		Status:     "blocked_until_real_equipment_evidence",
		Checks:     checks,
		NextActions: []string{
			"obtain a sample export or API schema",
			"map private raw fields before writing an importer",
			"add synthetic and redacted fixtures",
			"prove public export redaction before enabling dashboard data",
		},
	}, nil
}

func CheckAllEquipmentGates() []EquipmentGateReport {
	var reports []EquipmentGateReport
	for _, contract := range AdapterContracts() {
		report, err := CheckEquipmentGate(contract.Source.ID)
		if err == nil {
			reports = append(reports, report)
		}
	}
	return reports
}
