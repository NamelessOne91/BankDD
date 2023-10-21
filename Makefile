test:
	@ginkgo -r -p --race --randomize-all --randomize-suites --fail-on-pending --keep-going

coverage:
	@ginkgo -r -p  --race --randomize-all --randomize-suites --fail-on-pending --keep-going --cover --coverprofile=cover.profile

ci-test:
	@go run github.com/onsi/ginkgo/v2/ginkgo -r --procs=2 --compilers=2 --randomize-all --randomize-suites --fail-on-pending --keep-going --cover --coverprofile=cover.profile --race --trace --json-report=report.json --poll-progress-after=120s --poll-progress-interval=30s

coverage-report: coverage
	@go tool cover -html=cover.profile