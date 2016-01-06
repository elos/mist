build:
	go build ./...

gen:
	aeolus gen spec.json

serve:
	make ./cmd/serve/serve.go && ./cmd/serve/serve

