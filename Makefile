.PHONY: verify fmt-check lint test build vet privacy-check public-export-smoke web-install web-data web-check web-build web-dev web-preview demo

verify: fmt-check lint test build privacy-check public-export-smoke web-data web-check web-build

fmt-check:
	@test -z "$$(gofmt -l cmd internal)"

lint: vet

test:
	go test ./...

build:
	go build -o bin/flyto2 ./cmd/flyto2

vet:
	go vet ./...

privacy-check:
	go run ./cmd/flyto2 privacy check -data examples/synthetic_daily.csv

public-export-smoke:
	go run ./cmd/flyto2 export public -data examples/synthetic_daily.csv -out /tmp/flyto-health-public.json
	@test -s /tmp/flyto-health-public.json

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
