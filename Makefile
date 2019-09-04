#########################
###      DEFS         ###
#########################

# Don't ask, for to understand it is to look 
# into the void and know the void is not only 
#looking back but also reading your emails.
SHELL=/bin/bash -e -o pipefail

#########################
###      TARGETS      ###
#########################

.PHONY: test cover build

test: ## Runs basic go test
	go test -v ./... --cover --coverprofile=wurlwind.out

cover: ## Generate coverage report
	go tool cover --html=wurlwind.out

build: ## Build plugin binary
	go build -o terraform-provider-striketracker