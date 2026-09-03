
.PHONY: deps all test test-lmdb09 test-lmdb10 full-test bin bench check check-c coexist-test

deps:
	go mod download

bin:
	mkdir -p bin
	GOBIN=${PWD}/bin go install ./exp/cmd/...
	GOBIN=${PWD}/bin go install ./cmd/...

all: deps check check-c full-test bin

# The full test suite runs once per bundled LMDB engine: the process-wide
# default version decides which engine drives the databases the tests create.
test: test-lmdb09 test-lmdb10

test-lmdb09:
	go test -cover ./...

test-lmdb10:
	LMDBGO_DEFAULT_FORMAT=10 go test -count=1 -cover ./...

full-test: test
	go test -race ./...
	LMDBGO_DEFAULT_FORMAT=10 go test -race -count=1 ./...

bench:
	go test -run=NONE -bench=. -benchmem -count=6 ./lmdb

# C-level checks: shared-define parity between the two vendored headers, and
# the ASan reproducer for the patched LMDB 1.0 cursor-close use-after-free.
check-c:
	scripts/check-defines.sh
	scripts/test-cursor-close-asan.sh

# Verifies that lmdb-go v1 and v2 link and work together in one binary
# (requires network access to download v1).
coexist-test:
	cd tests/coexist && go test -count=1 ./...

check:
	which goimports > /dev/null
	find . -name '*.go' | grep -v '^\./\.claude/' | xargs goimports -d | tee /dev/stderr | wc -l | xargs test 0 -eq
	which golint > /dev/null
	golint ./... | tee /dev/stderr | wc -l | xargs test 0 -eq
