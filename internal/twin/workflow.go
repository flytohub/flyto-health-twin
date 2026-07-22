package twin

// WorkflowStep describes one non-executing recipe command and expected output.
type WorkflowStep struct {
	ID      string `json:"id"`
	Runtime string `json:"runtime"`
	Command string `json:"command"`
	Output  string `json:"output"`
}

// WorkflowRecipe documents orchestration metadata and its privacy boundary.
type WorkflowRecipe struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	PrimaryRuntime  string         `json:"primary_runtime"`
	HostedRuntime   string         `json:"hosted_runtime"`
	PrivacyBoundary string         `json:"privacy_boundary"`
	Steps           []WorkflowStep `json:"steps"`
}

// WorkflowRecipes returns documentation-only local and hosted recipes.
func WorkflowRecipes() []WorkflowRecipe {
	return []WorkflowRecipe{
		{
			ID:              "daily_public_export",
			Name:            "Daily public export refresh",
			Description:     "Regenerate predictions and redacted public JSON from daily aggregate input.",
			PrimaryRuntime:  "flyto-core local module workflow",
			HostedRuntime:   "flyto-cloud scheduled job after private source mounting is configured",
			PrivacyBoundary: "only web/public/public-data.json can be published",
			Steps: []WorkflowStep{
				{ID: "privacy_check", Runtime: "flyto-core", Command: "go run ./cmd/flyto2 privacy check -data examples/synthetic_daily.csv", Output: "privacy_check=pass"},
				{ID: "evaluate", Runtime: "flyto-core", Command: "go run ./cmd/flyto2 evaluate -data examples/synthetic_daily.csv", Output: "model errors"},
				{ID: "export", Runtime: "flyto-core", Command: "make web-data", Output: "web/public/public-data.json"},
				{ID: "build", Runtime: "GitHub Actions or flyto-cloud", Command: "make web-build", Output: "web/dist"},
			},
		},
		{
			ID:              "benchmark_regression",
			Name:            "Benchmark regression report",
			Description:     "Generate synthetic benchmark data and verify the baseline model stays within guardrails.",
			PrimaryRuntime:  "flyto-core local module workflow",
			HostedRuntime:   "flyto-cloud controlled scheduled report",
			PrivacyBoundary: "synthetic fixtures only",
			Steps: []WorkflowStep{
				{ID: "generate", Runtime: "flyto-core", Command: "go run ./cmd/flyto2 generate synthetic -profile balanced -days 30 -out /tmp/flyto2-benchmark.csv", Output: "synthetic CSV"},
				{ID: "benchmark", Runtime: "flyto-core", Command: "go run ./cmd/flyto2 benchmark run -profile balanced -days 30", Output: "JSON benchmark report"},
			},
		},
		{
			ID:              "equipment_gate_review",
			Name:            "Equipment adapter gate review",
			Description:     "Review every adapter contract before real device data is allowed into the model loop.",
			PrimaryRuntime:  "flyto-core local module workflow",
			HostedRuntime:   "flyto-cloud human-reviewed workflow",
			PrivacyBoundary: "raw samples stay private and never publish automatically",
			Steps: []WorkflowStep{
				{ID: "registry", Runtime: "flyto-core", Command: "go run ./cmd/flyto2 registry adapters", Output: "adapter contracts"},
				{ID: "gate", Runtime: "flyto-core", Command: "go run ./cmd/flyto2 equipment gate", Output: "gate reports"},
			},
		},
	}
}
