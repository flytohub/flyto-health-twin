package twin

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestAdapterRegistryIncludesImplementedCSVAndGatedRealSources(t *testing.T) {
	csv, ok := FindAdapterContract("manual_csv")
	if !ok {
		t.Fatal("manual_csv contract missing")
	}
	if !csv.PublicExportEligible || csv.Status != "implemented" {
		t.Fatalf("unexpected manual_csv contract: %#v", csv)
	}

	wearable, ok := FindAdapterContract("wearable_daily_aggregate")
	if !ok {
		t.Fatal("wearable_daily_aggregate contract missing")
	}
	if wearable.PublicExportEligible {
		t.Fatalf("real wearable contract should not be public-export eligible before gate: %#v", wearable)
	}
	if !wearable.RequiresConsent {
		t.Fatalf("real wearable contract should require consent: %#v", wearable)
	}
}

func TestSyntheticGeneratorIsDeterministicAndCSVLoadable(t *testing.T) {
	start := day("2026-06-01")
	first, err := GenerateSyntheticRecords("balanced", start, 12)
	if err != nil {
		t.Fatal(err)
	}
	second, err := GenerateSyntheticRecords("balanced", start, 12)
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 12 || len(second) != 12 {
		t.Fatalf("unexpected lengths %d %d", len(first), len(second))
	}
	if first[5] != second[5] {
		t.Fatalf("generator is not deterministic: %#v != %#v", first[5], second[5])
	}

	var buf bytes.Buffer
	if err := WriteDailyCSV(&buf, first); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(buf.String(), "gps") || strings.Contains(buf.String(), "access_token") {
		t.Fatalf("synthetic CSV should not include raw private fields: %s", buf.String())
	}
}

func TestBenchmarkReportSummarizesModelRegression(t *testing.T) {
	records, err := GenerateSyntheticRecords("balanced", day("2026-06-01"), 30)
	if err != nil {
		t.Fatal(err)
	}
	report, err := BuildBenchmarkReport("balanced", records)
	if err != nil {
		t.Fatal(err)
	}
	if report.ModelID != "baseline_strain" || report.EvaluationCount == 0 {
		t.Fatalf("unexpected report metadata: %#v", report)
	}
	if len(report.RecoveryStateCounts) == 0 {
		t.Fatalf("expected recovery state counts: %#v", report)
	}
	if report.ErrorSummary.HRVMeanAbsoluteError <= 0 {
		t.Fatalf("expected non-zero error summary: %#v", report.ErrorSummary)
	}
}

func TestEquipmentGateAllowsCSVAndBlocksRealAdapters(t *testing.T) {
	csv, err := CheckEquipmentGate("manual_csv")
	if err != nil {
		t.Fatal(err)
	}
	if csv.Status != "ready_for_synthetic_public_demo" {
		t.Fatalf("unexpected csv gate status: %#v", csv)
	}

	wearable, err := CheckEquipmentGate("wearable_daily_aggregate")
	if err != nil {
		t.Fatal(err)
	}
	if wearable.Status != "blocked_until_real_equipment_evidence" {
		t.Fatalf("unexpected wearable gate status: %#v", wearable)
	}
	for _, check := range wearable.Checks {
		if check.Required && !check.Passed {
			return
		}
	}
	t.Fatalf("expected at least one failed required check: %#v", wearable.Checks)
}

func TestWorkflowRecipesAreFlytoNative(t *testing.T) {
	recipes := WorkflowRecipes()
	if len(recipes) == 0 {
		t.Fatal("expected workflow recipes")
	}
	for _, recipe := range recipes {
		if !strings.Contains(recipe.PrimaryRuntime, "flyto-core") {
			t.Fatalf("expected flyto-core primary runtime: %#v", recipe)
		}
		hosted := strings.ToLower(recipe.HostedRuntime)
		if !strings.Contains(hosted, "flyto-cloud") && !strings.Contains(hosted, "github actions") {
			t.Fatalf("expected Flyto or CI hosted runtime: %#v", recipe)
		}
	}
}

func TestTelomereToyHasSafetyBoundary(t *testing.T) {
	result, err := RunTelomereToy(TelomereToyParams{
		InitialLengthKB: 10,
		Divisions:       20,
		StressIndex:     0.5,
		RepairBias:      0.1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.FinalLengthKB >= result.Params.InitialLengthKB {
		t.Fatalf("expected toy shortening, got %#v", result)
	}
	if !strings.Contains(result.Boundary, "educational toy model") {
		t.Fatalf("missing safety boundary: %#v", result)
	}
	if len(result.NotAllowedInterpretation) == 0 {
		t.Fatalf("expected explicit not-allowed interpretations: %#v", result)
	}
}

func TestModelAndDatasetRegistriesArePublicSafe(t *testing.T) {
	if len(ModelRegistry()) < 2 {
		t.Fatal("expected implemented and planned model cards")
	}
	if len(DatasetRegistry()) < 2 {
		t.Fatal("expected dataset registry entries")
	}
}

func TestGenerateSyntheticRejectsTooFewDays(t *testing.T) {
	_, err := GenerateSyntheticRecords("balanced", time.Time{}, 2)
	if err == nil {
		t.Fatal("expected error for too few days")
	}
}
