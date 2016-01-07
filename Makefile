build:
	go build ./...

gen:
	aeolus gen spec.json

serve:
	go build ./cmd/serve/serve.go && ./cmd/serve/serve

prepare:
	git add --all && git commit && git push

deploy:
	ssh -i ~/.ssh/elos.pem ubuntu@ec2-52-88-164-210.us-west-2.compute.amazonaws.com && ./update.sh
