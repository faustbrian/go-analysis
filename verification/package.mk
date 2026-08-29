.PHONY: build compatibility docs package-specific

build:
	mkdir -p .build
	go build -trimpath -o .build/golib-analysis ./cmd/golib-analysis

compatibility:
	./scripts/compatibility.sh check

docs:
	cd policy && go test . -run='^TestDocumentation$$'

package-specific: build
	./scripts/toolchain.sh
	./scripts/toolchain_test.sh
	./scripts/corpus_test.sh
	./scripts/owned_corpus_test.sh
	./scripts/performance_test.sh
	./scripts/corpus.sh check corpus/manifest.tsv
	./scripts/performance.sh corpus/performance.tsv
	./scripts/reproducible-build.sh
	./scripts/verify-release.sh 1.0.0
	go vet -vettool=.build/golib-analysis ./analysis ./policy ./internal/driver
