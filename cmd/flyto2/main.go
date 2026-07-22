package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/flytohub/flyto-health-twin/internal/twin"
)

// main dispatches the local Flyto2 Health Twin command tree.
func main() {
	log.SetFlags(0)
	if len(os.Args) == 1 || strings.HasPrefix(os.Args[1], "-") {
		runEvaluate(os.Args[1:])
		return
	}

	switch os.Args[1] {
	case "demo":
		runEvaluate(append([]string{"-data", "examples/synthetic_daily.csv"}, os.Args[2:]...))
	case "evaluate":
		runEvaluate(os.Args[2:])
	case "predict":
		runPredict(os.Args[2:])
	case "export":
		runExport(os.Args[2:])
	case "generate":
		runGenerate(os.Args[2:])
	case "import":
		runImport(os.Args[2:])
	case "registry":
		runRegistry(os.Args[2:])
	case "report":
		runReport(os.Args[2:])
	case "benchmark":
		runBenchmark(os.Args[2:])
	case "equipment":
		runEquipment(os.Args[2:])
	case "simulate":
		runSimulate(os.Args[2:])
	case "workflow":
		runWorkflow(os.Args[2:])
	case "privacy":
		runPrivacy(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		log.Fatalf("unknown command %q\n\n%s", os.Args[1], usageText())
	}
}

// runGenerate writes deterministic synthetic CSV records.
func runGenerate(args []string) {
	if len(args) == 0 || args[0] != "synthetic" {
		log.Fatalf("expected: generate synthetic\n\n%s", usageText())
	}
	fs := flag.NewFlagSet("generate synthetic", flag.ExitOnError)
	profileID := fs.String("profile", "balanced", "synthetic profile id")
	days := fs.Int("days", 0, "number of generated days; 0 uses profile default")
	startRaw := fs.String("start", "2026-06-01", "start date in YYYY-MM-DD")
	outPath := fs.String("out", "-", "output CSV path; '-' writes stdout")
	_ = fs.Parse(args[1:])

	start := mustParseDate(*startRaw)
	records, err := twin.GenerateSyntheticRecords(*profileID, start, *days)
	if err != nil {
		log.Fatal(err)
	}
	if *outPath == "-" {
		if err := twin.WriteDailyCSV(os.Stdout, records); err != nil {
			log.Fatal(err)
		}
		return
	}
	f, err := os.Create(*outPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := twin.WriteDailyCSV(f, records); err != nil {
		_ = f.Close()
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

// runEvaluate prints rolling predictions and aggregate errors.
func runEvaluate(args []string) {
	fs := flag.NewFlagSet("evaluate", flag.ExitOnError)
	dataPath := fs.String("data", "examples/synthetic_daily.csv", "CSV file with daily aggregate records")
	limit := fs.Int("limit", 0, "maximum evaluation rows to print; 0 prints all")
	_ = fs.Parse(args)

	records := mustLoadRecords(*dataPath)
	evals, err := twin.Evaluate(records)
	if err != nil {
		log.Fatal(err)
	}

	printEvaluations(evals, *limit)
	hrv, rhr, fatigue, sleep := twin.MeanAbsoluteErrors(evals)
	fmt.Printf("\nmean_absolute_error hrv=%.2f rhr=%.2f fatigue=%.2f sleep=%.2f\n", hrv, rhr, fatigue, sleep)
}

// runPredict writes the latest next-day prediction as JSON.
func runPredict(args []string) {
	fs := flag.NewFlagSet("predict", flag.ExitOnError)
	dataPath := fs.String("data", "examples/synthetic_daily.csv", "CSV file with daily aggregate records")
	_ = fs.Parse(args)

	records := mustLoadRecords(*dataPath)
	prediction, err := twin.PredictNext(records)
	if err != nil {
		log.Fatal(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(map[string]any{
		"prediction": prediction,
	}); err != nil {
		log.Fatal(err)
	}
}

// runExport writes the explicit public-data JSON allowlist.
func runExport(args []string) {
	if len(args) == 0 || args[0] != "public" {
		log.Fatalf("expected: export public\n\n%s", usageText())
	}
	fs := flag.NewFlagSet("export public", flag.ExitOnError)
	dataPath := fs.String("data", "examples/synthetic_daily.csv", "CSV file with daily aggregate records")
	outPath := fs.String("out", "-", "output JSON path; '-' writes stdout")
	generatedAt := fs.String("generated-at", "", "optional RFC3339 timestamp for reproducible public exports")
	_ = fs.Parse(args[1:])

	records := mustLoadRecords(*dataPath)
	evals, err := twin.Evaluate(records)
	if err != nil {
		log.Fatal(err)
	}
	exportedAt := time.Now().UTC()
	if *generatedAt != "" {
		exportedAt, err = time.Parse(time.RFC3339, *generatedAt)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *outPath == "-" {
		if err := twin.WritePublicExportAt(os.Stdout, records, evals, exportedAt); err != nil {
			log.Fatal(err)
		}
		return
	}

	f, err := os.Create(*outPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := twin.WritePublicExportAt(f, records, evals, exportedAt); err != nil {
		_ = f.Close()
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

// runRegistry writes all or one metadata registry as JSON.
func runRegistry(args []string) {
	if len(args) == 0 {
		writeJSON(map[string]any{
			"adapters":   twin.AdapterContracts(),
			"models":     twin.ModelRegistry(),
			"datasets":   twin.DatasetRegistry(),
			"synthetic":  twin.SyntheticProfiles(),
			"benchmarks": twin.BenchmarkFixtures(),
			"workflows":  twin.WorkflowRecipes(),
		})
		return
	}
	switch args[0] {
	case "adapters":
		writeJSON(twin.AdapterContracts())
	case "models":
		writeJSON(twin.ModelRegistry())
	case "datasets":
		writeJSON(twin.DatasetRegistry())
	case "synthetic":
		writeJSON(twin.SyntheticProfiles())
	case "benchmarks":
		writeJSON(twin.BenchmarkFixtures())
	case "workflows":
		writeJSON(twin.WorkflowRecipes())
	default:
		log.Fatalf("unknown registry %q\n\n%s", args[0], usageText())
	}
}

// runReport writes one baseline model evidence report.
func runReport(args []string) {
	if len(args) == 0 || args[0] != "model" {
		log.Fatalf("expected: report model\n\n%s", usageText())
	}
	fs := flag.NewFlagSet("report model", flag.ExitOnError)
	dataPath := fs.String("data", "examples/synthetic_daily.csv", "CSV file with daily aggregate records")
	datasetID := fs.String("dataset", "synthetic_daily_v0", "dataset id for report provenance")
	outPath := fs.String("out", "-", "output JSON path; '-' writes stdout")
	_ = fs.Parse(args[1:])

	records := mustLoadRecords(*dataPath)
	report, err := twin.BuildModelEvaluationReport(twin.BaselineModel{}, records, *datasetID)
	if err != nil {
		log.Fatal(err)
	}
	writeJSONPath(*outPath, report)
}

// runBenchmark generates and evaluates a deterministic benchmark profile.
func runBenchmark(args []string) {
	if len(args) == 0 || args[0] != "run" {
		log.Fatalf("expected: benchmark run\n\n%s", usageText())
	}
	fs := flag.NewFlagSet("benchmark run", flag.ExitOnError)
	profileID := fs.String("profile", "balanced", "synthetic benchmark profile id")
	days := fs.Int("days", 30, "number of generated synthetic days")
	startRaw := fs.String("start", "2026-06-01", "start date in YYYY-MM-DD")
	outPath := fs.String("out", "-", "output JSON path; '-' writes stdout")
	_ = fs.Parse(args[1:])

	records, err := twin.GenerateSyntheticRecords(*profileID, mustParseDate(*startRaw), *days)
	if err != nil {
		log.Fatal(err)
	}
	report, err := twin.BuildBenchmarkReport(*profileID, records)
	if err != nil {
		log.Fatal(err)
	}
	writeJSONPath(*outPath, report)
}

// runEquipment writes one or all adapter readiness reports.
func runEquipment(args []string) {
	if len(args) == 0 || args[0] != "gate" {
		log.Fatalf("expected: equipment gate\n\n%s", usageText())
	}
	fs := flag.NewFlagSet("equipment gate", flag.ExitOnError)
	sourceID := fs.String("source", "", "adapter source id; empty checks all")
	_ = fs.Parse(args[1:])

	if *sourceID == "" {
		writeJSON(twin.CheckAllEquipmentGates())
		return
	}
	report, err := twin.CheckEquipmentGate(*sourceID)
	if err != nil {
		log.Fatal(err)
	}
	writeJSON(report)
}

// runSimulate writes the bounded educational telomere toy result.
func runSimulate(args []string) {
	if len(args) == 0 || args[0] != "telomere" {
		log.Fatalf("expected: simulate telomere\n\n%s", usageText())
	}
	fs := flag.NewFlagSet("simulate telomere", flag.ExitOnError)
	initialLength := fs.Float64("initial-kb", 10, "toy initial telomere length in kb")
	divisions := fs.Int("divisions", 24, "toy cell division count")
	stressIndex := fs.Float64("stress", 0.35, "toy stress index from 0 to 1")
	repairBias := fs.Float64("repair", 0.1, "toy repair bias from 0 to 1")
	outPath := fs.String("out", "-", "output JSON path; '-' writes stdout")
	_ = fs.Parse(args[1:])

	result, err := twin.RunTelomereToy(twin.TelomereToyParams{
		InitialLengthKB: *initialLength,
		Divisions:       *divisions,
		StressIndex:     *stressIndex,
		RepairBias:      *repairBias,
	})
	if err != nil {
		log.Fatal(err)
	}
	writeJSONPath(*outPath, result)
}

// runWorkflow writes documentation-only workflow recipes.
func runWorkflow(args []string) {
	if len(args) == 0 || args[0] != "recipes" {
		log.Fatalf("expected: workflow recipes\n\n%s", usageText())
	}
	writeJSON(twin.WorkflowRecipes())
}

// runImport validates CSV and prints normalized record provenance.
func runImport(args []string) {
	if len(args) == 0 || args[0] != "csv" {
		log.Fatalf("expected: import csv\n\n%s", usageText())
	}
	fs := flag.NewFlagSet("import csv", flag.ExitOnError)
	dataPath := fs.String("data", "examples/synthetic_daily.csv", "CSV file with daily aggregate records")
	_ = fs.Parse(args[1:])

	records := mustLoadRecords(*dataPath)
	if len(records) == 0 {
		fmt.Println("records=0")
		return
	}
	fmt.Printf("records=%d start=%s end=%s source=%s\n",
		len(records),
		records[0].Date.Format(twin.DateLayout),
		records[len(records)-1].Date.Format(twin.DateLayout),
		twin.CSVAdapter{}.Source().ID,
	)
}

// runPrivacy blocks forbidden headers and reviewable note content.
func runPrivacy(args []string) {
	if len(args) == 0 || args[0] != "check" {
		log.Fatalf("expected: privacy check\n\n%s", usageText())
	}
	fs := flag.NewFlagSet("privacy check", flag.ExitOnError)
	dataPath := fs.String("data", "examples/synthetic_daily.csv", "CSV file with daily aggregate records")
	_ = fs.Parse(args[1:])

	issues, err := twin.InspectCSVPrivacy(*dataPath)
	if err != nil {
		log.Fatal(err)
	}
	records := mustLoadRecords(*dataPath)
	for _, r := range records {
		for _, warning := range twin.PrivacyWarnings(r) {
			issues = append(issues, twin.PrivacyIssue{
				Field:   "notes",
				Reason:  warning,
				Level:   "review",
				Example: r.Date.Format(twin.DateLayout),
			})
		}
	}

	if len(issues) == 0 {
		fmt.Println("privacy_check=pass issues=0")
		return
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(map[string]any{"privacy_check": "fail", "issues": issues}); err != nil {
		log.Fatal(err)
	}
	os.Exit(1)
}

// mustLoadRecords loads a CSV path or terminates with a concise CLI error.
func mustLoadRecords(dataPath string) []twin.DailyRecord {
	records, err := twin.LoadCSV(dataPath)
	if err != nil {
		log.Fatal(err)
	}
	return records
}

// mustParseDate parses the canonical daily date or terminates the CLI.
func mustParseDate(raw string) time.Time {
	parsed, err := time.Parse(twin.DateLayout, raw)
	if err != nil {
		log.Fatal(err)
	}
	return parsed
}

// writeJSONPath writes indented JSON to stdout or a truncating local file.
func writeJSONPath(path string, value any) {
	if path == "-" {
		writeJSON(value)
		return
	}
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(value); err != nil {
		_ = f.Close()
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

// writeJSON writes indented JSON to stdout.
func writeJSON(value any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(value); err != nil {
		log.Fatal(err)
	}
}

// printEvaluations renders the human-readable rolling error table.
func printEvaluations(evals []twin.Evaluation, limit int) {
	fmt.Println("date        model             hrv_pred hrv_actual hrv_err rhr_pred rhr_actual fatigue_pred fatigue_actual sleep_pred sleep_actual recovery hints")

	count := len(evals)
	if limit > 0 && limit < count {
		count = limit
	}

	for _, e := range evals[:count] {
		p := e.Prediction
		fmt.Printf(
			"%s  %-16s %.1f     %.1f       %+.1f    %.1f     %.1f       %.1f          %.1f            %.1f      %.1f        %s %s\n",
			e.TargetDate.Format(twin.DateLayout),
			p.ModelID+"@"+p.ModelVersion,
			p.PredictedHRV,
			e.Actual.HRV,
			e.HRVError,
			p.PredictedRestingHeartRate,
			e.Actual.RestingHeartRate,
			p.PredictedFatigueScore,
			e.Actual.FatigueScore,
			p.PredictedSleepScore,
			e.Actual.SleepScore,
			p.RecoveryState,
			strings.Join(p.Hints, "; "),
		)
	}
}

// printUsage writes the complete command synopsis.
func printUsage() {
	fmt.Print(usageText())
}

// usageText returns the complete command synopsis.
func usageText() string {
	return `Flyto2

Usage:
  flyto2 demo [-limit N]
  flyto2 evaluate [-data path] [-limit N]
  flyto2 predict [-data path]
  flyto2 export public [-data path] [-out path|-] [-generated-at RFC3339]
  flyto2 generate synthetic [-profile balanced] [-days N] [-out path|-]
  flyto2 import csv [-data path]
  flyto2 registry [adapters|models|datasets|synthetic|benchmarks|workflows]
  flyto2 report model [-data path] [-dataset id] [-out path|-]
  flyto2 benchmark run [-profile balanced] [-days N] [-out path|-]
  flyto2 equipment gate [-source id]
  flyto2 simulate telomere [-divisions N] [-out path|-]
  flyto2 workflow recipes
  flyto2 privacy check [-data path]

Default:
  flyto2 -data examples/synthetic_daily.csv
`
}
