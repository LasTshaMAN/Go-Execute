test:
	go test ./... -v -race -count=1

goimports:
	goimports -local "github.com/LasTshaMAN/" -w ./
