GO    := go
BIN   := boast
COV   := cov.out
RM_RF := rm -rf

all: $(BIN)

.PHONY: $(BIN)
$(BIN):
	@echo "> Building $(BIN)..."
	$(GO) build -o $(BIN) cmd/boast/*.go
	@echo "> Done."

.PHONY: test
test:
	@echo "> Testing..."
	$(GO) test ./...
	@echo "> Done."

test_verbose:
	@echo "> Testing..."
	$(GO) test -v ./...
	@echo "> Done."

.PHONY: cover
cover:
	@echo "> Generating coverage profile..."
	$(GO) test -covermode=atomic -coverprofile=$(COV) ./...
	@echo "> Done."

.PHONY: cover_html
cover_html: cover
	@echo "> Generating HTML output for coverage profile..."
	$(GO) tool cover -html=cov.out 
	@echo "> Done."

.PHONY: cover_clean
cover_clean:
	$(RM_RF) $(COV)

.PHONY: clean
clean:
	$(RM_RF) $(BIN)
