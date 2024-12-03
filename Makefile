TEST_COVERAGE_OUTPUT_ROOT := $(CURDIR)/.coverage/unit
INTEGRATION_COVERAGE_OUTPUT_ROOT := $(CURDIR)/.coverage/integration
COVERAGE_PROFILE := $(CURDIR)/.coverage/profile
COVERAGE_REPORT := $(CURDIR)/.coverage/report.html

test: unit-test integration-test

unit-test:
	go test ./app/{billing,order,shipment}

integration-test:
	go test -tags=integration ./app/test

$(TEST_COVERAGE_OUTPUT_ROOT):
	mkdir -p $(TEST_COVERAGE_OUTPUT_ROOT)

$(INTEGRATION_COVERAGE_OUTPUT_ROOT):
	mkdir -p $(INTEGRATION_COVERAGE_OUTPUT_ROOT)

$(SUMMARY_COVERAGE_OUTPUT_ROOT):
	mkdir -p $(SUMMARY_COVERAGE_OUTPUT_ROOT)

unit-test-coverage: $(TEST_COVERAGE_OUTPUT_ROOT)
	@echo Unit test coverage
	go test -cover ./app/billing -args -test.gocoverdir=$(TEST_COVERAGE_OUTPUT_ROOT) 
	go test -cover ./app/order -args -test.gocoverdir=$(TEST_COVERAGE_OUTPUT_ROOT) 
	go test -cover ./app/shipment -args -test.gocoverdir=$(TEST_COVERAGE_OUTPUT_ROOT) 

integration-test-coverage: $(INTEGRATION_COVERAGE_OUTPUT_ROOT)
	@echo Integration test coverage
	go test -tags integration -cover ./app/test -coverpkg ./... -args -test.gocoverdir=$(INTEGRATION_COVERAGE_OUTPUT_ROOT)

coverage-report: $(TEST_COVERAGE_OUTPUT_ROOT) $(INTEGRATION_COVERAGE_OUTPUT_ROOT) $(SUMMARY_COVERAGE_OUTPUT_ROOT)
	@echo Summary coverage report
	go tool covdata textfmt -i $(TEST_COVERAGE_OUTPUT_ROOT),$(INTEGRATION_COVERAGE_OUTPUT_ROOT) -o $(COVERAGE_PROFILE)
	go tool covdata percent -i $(TEST_COVERAGE_OUTPUT_ROOT),$(INTEGRATION_COVERAGE_OUTPUT_ROOT)
	go tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_REPORT)
