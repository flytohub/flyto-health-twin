package twin

import (
	"errors"
	"math"
)

type TelomereToyParams struct {
	InitialLengthKB float64 `json:"initial_length_kb"`
	Divisions       int     `json:"divisions"`
	StressIndex     float64 `json:"stress_index"`
	RepairBias      float64 `json:"repair_bias"`
}

type TelomereToyPoint struct {
	Division int     `json:"division"`
	LengthKB float64 `json:"length_kb"`
	State    string  `json:"state"`
}

type TelomereToyResult struct {
	SimulationID             string             `json:"simulation_id"`
	Boundary                 string             `json:"boundary"`
	Params                   TelomereToyParams  `json:"params"`
	ShorteningPerDivisionBP  float64            `json:"shortening_per_division_bp"`
	FinalLengthKB            float64            `json:"final_length_kb"`
	Points                   []TelomereToyPoint `json:"points"`
	Caveats                  []string           `json:"caveats"`
	ResearchQuestions        []string           `json:"research_questions"`
	NotAllowedInterpretation []string           `json:"not_allowed_interpretation"`
}

func RunTelomereToy(params TelomereToyParams) (TelomereToyResult, error) {
	if params.InitialLengthKB == 0 {
		params.InitialLengthKB = 10
	}
	if params.Divisions == 0 {
		params.Divisions = 24
	}
	if params.InitialLengthKB < 2 || params.InitialLengthKB > 20 {
		return TelomereToyResult{}, errors.New("initial telomere length must be between 2 and 20 kb")
	}
	if params.Divisions < 1 || params.Divisions > 500 {
		return TelomereToyResult{}, errors.New("divisions must be between 1 and 500")
	}
	params.StressIndex = clamp(params.StressIndex, 0, 1)
	params.RepairBias = clamp(params.RepairBias, 0, 1)

	shorteningBP := 45 + params.StressIndex*25 - params.RepairBias*20
	if shorteningBP < 10 {
		shorteningBP = 10
	}
	length := params.InitialLengthKB
	points := []TelomereToyPoint{{Division: 0, LengthKB: round2(length), State: toyCellState(length)}}
	for i := 1; i <= params.Divisions; i++ {
		length -= shorteningBP / 1000
		if length < 0 {
			length = 0
		}
		if i == params.Divisions || i%5 == 0 {
			points = append(points, TelomereToyPoint{Division: i, LengthKB: round2(length), State: toyCellState(length)})
		}
	}

	return TelomereToyResult{
		SimulationID:            "telomere_toy_v0",
		Boundary:                "educational toy model only; not personal biology, diagnosis, treatment, or anti-aging guidance",
		Params:                  params,
		ShorteningPerDivisionBP: round1(shorteningBP),
		FinalLengthKB:           round2(length),
		Points:                  points,
		Caveats: []string{
			"telomere dynamics vary by cell type, measurement method, and biology not represented here",
			"telomerase concepts can involve cancer-risk tradeoffs and are not recommendations",
			"this simulation must remain separate from personal wearable predictions",
		},
		ResearchQuestions: []string{
			"what public papers are appropriate for concept-level parameters",
			"how should biology explainers be separated from the personal response model",
		},
		NotAllowedInterpretation: []string{
			"do not infer a user's telomere length",
			"do not claim rejuvenation, repair, treatment, or lifespan extension",
		},
	}, nil
}

func toyCellState(lengthKB float64) string {
	switch {
	case lengthKB <= 3:
		return "toy-critical"
	case lengthKB <= 5:
		return "toy-short"
	default:
		return "toy-normal"
	}
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
