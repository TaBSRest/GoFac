runTestEnvironment:
	docker compose run --rm --build default bash

test:
	python3 ./buildscript.py $(ARGS)

checkCoverage:
	@echo "Checking the Total Line Coverage..."
	@go test ./... -race -coverprofile=./coverage/cover.txt > ./coverage/dump.txt
	@go tool cover -func=./coverage/cover.txt | fgrep total | awk '{print $$3}'
	@echo "All Done!"

updateCoverage:
	@echo "Updating the Total Line Coverage..."
	@go test ./... -race -coverprofile=./coverage/cover.txt > ./coverage/dump.txt
	@go tool cover -func=./coverage/cover.txt | fgrep total | awk '{print $$3}'
	@go tool cover -func=./coverage/cover.txt | fgrep total | awk '{print $$3}' > ./coverage/LineCoverage.txt
	@echo "All Done!"

adjustCoverage:
	@echo "Adjusting the Total Line Coverage Manually..."
	@python3 ./adjustcoverage.py $(COVERAGE)
	@echo "All Done!"

generateCoverageReport:
	@echo "Generating Coverage Report..."
	@go test ./... -race -coverprofile=./coverage/cover.txt > ./coverage/dump.txt
	@go tool cover -html=./coverage/cover.txt -o ./coverage/Results.html
	@echo "All Done!"
