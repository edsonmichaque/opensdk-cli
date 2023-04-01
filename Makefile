COPYRIGHT_HOLDER := "Edson Michaque"
COPYRIGHT_YEARS := "2023"

.PHONY: update-deps
update-deps:
	go get -u ./...

.PHONY: build
build:
	go build -o ./bin/template cmd/template/main.go

.PHONY: test
test:
	go test -race -v ./...

.PHONY: dep
dep:
	go mod download

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: release
release:
	goreleaser release --clean

.PHONY: addlicense
addlicense:
	go install github.com/google/addlicense@latest

.PHONY: copyright
copyright: addlicense
	addlicense -c ${COPYRIGHT_HOLDER} -y ${COPYRIGHT_YEARS} -l apache -s  .

.PHONY: check-license
check-license: addlicense
	addlicense -check .

.PHONY: e2e
e2e: build bats
	./e2e/bats/bin/bats e2e
	go test -race -tags=e2e -v ./...


.PHONY: bats
bats:
	bash hack/bats.sh
