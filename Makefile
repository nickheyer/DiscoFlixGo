# Install Ent code-generation module
.PHONY: ent-install
ent-install:
	go get -d entgo.io/ent/cmd/ent

# Generate Ent code
.PHONY: ent-gen
ent-gen:
	go generate ./ent

# Create a new Ent entity
.PHONY: ent-new
ent-new:
	go run entgo.io/ent/cmd/ent new $(name)

# Run the application
.PHONY: run
run:
	go run cmd/web/main.go

# Run all tests
.PHONY: test
test:
	go test -count=1 -p 1 ./...

# Check for direct dependency updates
.PHONY: check-updates
check-updates:
	go list -u -m -f '{{if not .Indirect}}{{.}}{{end}}' all | grep "\["

# Clean the build cache and module cache
.PHONY: clean
clean:
	go clean -cache -modcache -i -r

# Clean and rebuild the project
.PHONY: rebuild
rebuild: clean
	go build cmd/web/main.go

# Tidy go.mod by removing unused dependencies
.PHONY: tidy
tidy:
	go mod tidy

# Install all dependencies
.PHONY: install-deps
install-deps:
	go mod download

# Generate all necessary code and build the project
.PHONY: all
all: ent-gen tidy install-deps rebuild

