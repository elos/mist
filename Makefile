build:
	go build ./...

gen:
	aeolus gen spec.json

serve:
	make build && ./cmd/serve/serve

