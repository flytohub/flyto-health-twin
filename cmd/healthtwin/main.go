package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/flytohub/flyto-health-twin/internal/twin"
)

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
	case "import":
		runImport(os.Args[2:])
	case "privacy":
		runPrivacy(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		log.Fatalf("unknown command %q\n\n%s", os.Args[1], usageText())
	}
}

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

func runExport(args []string) {
	if len(args) == 0 || args[0] != "public" {
		log.Fatalf("expected: export public\n\n%s", usageText())
	}
	fs := flag.NewFlagSet("export public", flag.ExitOnError)
	dataPath := fs.String("data", "examples/synthetic_daily.csv", "CSV file with daily aggregate records")
	outPath := fs.String("out", "-", "output JSON path; '-' writes stdout")
	_ = fs.Parse(args[1:])

	records := mustLoadRecords(*dataPath)
	evals, err := twin.Evaluate(records)
	if err != nil {
		log.Fatal(err)
	}

	if *outPath == "-" {
		if err := twin.WritePublicExport(os.Stdout, records, evals); err != nil {
			log.Fatal(err)
		}
		return
	}

	f, err := os.Create(*outPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := twin.WritePublicExport(f, records, evals); err != nil {
		log.Fatal(err)
	}
}

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

func mustLoadRecords(dataPath string) []twin.DailyRecord {
	records, err := twin.LoadCSV(dataPath)
	if err != nil {
		log.Fatal(err)
	}
	return records
}

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

func printUsage() {
	fmt.Print(usageText())
}

func usageText() string {
	return `Flyto Health Twin

Usage:
  healthtwin demo [-limit N]
  healthtwin evaluate [-data path] [-limit N]
  healthtwin predict [-data path]
  healthtwin export public [-data path] [-out path|-]
  healthtwin import csv [-data path]
  healthtwin privacy check [-data path]

Default:
  healthtwin -data examples/synthetic_daily.csv
`
}
