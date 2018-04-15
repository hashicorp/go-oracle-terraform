TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

test: fmtcheck errcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=60m -parallel=4

testacc: 
	ORACLE_ACC=1 go test -v $(TEST) $(TESTARGS) -timeout 120m

testrace: fmtcheck
	ORACLE_ACC= go test -race $(TEST) $(TESTARGS)

cover:
	@go tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u golang.org/x/tools/cmd/cover; \
	fi
	go test $(TEST) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	echo "'$(CURDIR)/scripts/gofmtcheck.sh'"
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

.PHONY: tools build test testacc testrace cover vet fmt fmtcheck errcheck

test-analyze:

	@echo "Installing test plugins"
	go get -u github.com/jstemmer/go-junit-report
	go get -u github.com/axw/gocov/...
	go get -u github.com/AlekSi/gocov-xml
	go get -u gopkg.in/alecthomas/gometalinter.v1
	../bin/gometalinter.v1 --install

	@rm -rf reports
	@mkdir -p ../reports/coverage
	@mkdir -p ../reports/unit-tests
	@mkdir -p ../reports/golint

	@echo "... Running golint"
	../bin/gometalinter.v1 --checkstyle ./... > ../reports/golint/report.xml

	@echo "... Running Go Unit Tests and coverage"
	@go test "./..." -coverprofile=../reports/coverage/cover.out -v 2>&1 | "../bin/go-junit-report" > ../reports/unit-tests/unit-test-report.xml

	@echo "... Generating coverage reports"
	@go tool cover -html=../reports/coverage/cover.out -o ../reports/coverage/coverage.html
	@../bin/gocov convert ../reports/coverage/cover.out | ../bin/gocov-xml > ../reports/coverage/coverage.xml

