build:
	go build ./...

gen:
	aeolus gen spec.json

serve:
	go build ./cmd/serve/serve.go && ./cmd/serve/serve

