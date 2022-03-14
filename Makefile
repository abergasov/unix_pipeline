check: ## check that tests are passing
	go test -v -race

run: ## run chain of command calculating hash for fibonacci sequence
	go run .