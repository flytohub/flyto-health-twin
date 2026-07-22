.PHONY: verify fmt-check lint docs-check test build vet privacy-check public-export-smoke roadmap-smoke fixtures web-install web-data web-check web-build web-dev web-preview demo

verify: fmt-check lint test build privacy-check public-export-smoke roadmap-smoke web-data web-check web-build

fmt-check:
	@test -z "$$(gofmt -l cmd internal scripts)"

lint: vet docs-check

docs-check:
	go run ./scripts/generate-reference.go
	npm --prefix web run docs:check

test:
	go test ./...

build:
	go build -o bin/flyto2 ./cmd/flyto2

vet:
	go vet ./...

privacy-check:
	go run ./cmd/flyto2 privacy check -data examples/synthetic_daily.csv
	go run ./cmd/flyto2 privacy check -data examples/benchmark_balanced.csv

public-export-smoke:
	go run ./cmd/flyto2 export public -data examples/synthetic_daily.csv -out /tmp/flyto-health-public.json
	@test -s /tmp/flyto-health-public.json

roadmap-smoke:
	go run ./cmd/flyto2 registry adapters >/tmp/flyto2-adapters.json
	go run ./cmd/flyto2 registry models >/tmp/flyto2-models.json
	go run ./cmd/flyto2 registry datasets >/tmp/flyto2-datasets.json
	go run ./cmd/flyto2 registry workflows >/tmp/flyto2-workflows.json
	go run ./cmd/flyto2 report model -data examples/synthetic_daily.csv >/tmp/flyto2-model-report.json
	go run ./cmd/flyto2 benchmark run -profile balanced -days 30 >/tmp/flyto2-benchmark-report.json
	go run ./cmd/flyto2 equipment gate >/tmp/flyto2-equipment-gate.json
	go run ./cmd/flyto2 simulate telomere -divisions 24 >/tmp/flyto2-telomere-toy.json
	@test -s /tmp/flyto2-adapters.json
	@test -s /tmp/flyto2-benchmark-report.json
	@test -s /tmp/flyto2-equipment-gate.json
	@test -s /tmp/flyto2-telomere-toy.json

fixtures:
	go run ./cmd/flyto2 generate synthetic -profile balanced -days 30 -start 2026-06-01 -out examples/benchmark_balanced.csv

web-install:
	npm install --prefix web

web-data:
	go run ./cmd/flyto2 export public -data examples/synthetic_daily.csv -out web/public/public-data.json -generated-at 2026-06-30T00:00:00Z

web-check:
	npm --prefix web run check

web-build:
	npm --prefix web run build

web-dev: web-data
	npm --prefix web run dev

web-preview: web-data web-build
	npm --prefix web run preview

demo:
	go run ./cmd/flyto2 demo -limit 5
