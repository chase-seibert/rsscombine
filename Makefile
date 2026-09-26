.PHONY: install run docker-build docker-run

IMAGE ?= rsscombine:local
DOCKER ?= $(shell if command -v docker >/dev/null 2>&1; then command -v docker; else echo /Users/cseibert/.docker/bin/docker; fi)

install:
	go mod download

run:
	set -a && source .env && source $(file) && set +a && go run rsscombine.go

docker-build:
	$(DOCKER) build -t $(IMAGE) .

docker-run:
	@test -n "$(file)" || (echo "Usage: make docker-run file=profile.env" >&2; exit 1)
	@test -f "$(CURDIR)/$(file)" || (echo "Configuration file not found: $(CURDIR)/$(file)" >&2; exit 1)
	$(DOCKER) run --rm \
		--env-file $(CURDIR)/.env \
		--env-file $(CURDIR)/$(file) \
		-v $(CURDIR)/rsscombine.yml:/app/rsscombine.yml:ro \
		$(IMAGE)
