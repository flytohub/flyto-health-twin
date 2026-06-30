.PHONY: verify fmt-check lint test build vet privacy-check public-export-smoke demo

verify: fmt-check lint test build privacy-check public-export-smoke

fmt-check:
	@test -z "$$(gofmt -l cmd internal)"

lint: vet

test:
	go test ./...

build:
	go build -o bin/healthtwin ./cmd/healthtwin

vet:
	go vet ./...

privacy-check:
	go run ./cmd/healthtwin privacy check -data examples/synthetic_daily.csv

public-export-smoke:
	go run ./cmd/healthtwin export public -data examples/synthetic_daily.csv -out /tmp/flyto-health-public.json
	@test -s /tmp/flyto-health-public.json

demo:
	go run ./cmd/healthtwin demo -limit 5
