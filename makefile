PWSH := pwsh.exe -NoProfile -Command

.PHONY: test build stop open docs dev windows darwin

test:
	go test

build:
	go build

stop:
	./nest.exe -stop

open:
	wt.exe -w 0 nt -d . $(PWSH) ./nest.exe -start


docs:
	swag init -d .\core\,.\handlers\,.\plugins\nest\ -g .\core.go --parseInternal --parseFuncBody --parseDependency


dev: stop build open


# TODO: rewrite build-pub in this;
# BUILD-PUB

GOARCHList := amd64 arm64

# HACK: ew
windows:
	@echo "Building for Windows..."
	$(PWSH) "'amd64', 'arm64' | ForEach-Object { $$arch = $$_; $$env:GOOS='windows'; $$env:GOARCH=$$arch; go build -ldflags -H=windowsgui -o build/nest-windows-$$arch.exe; Write-Host \"Built build/nest-windows-$$arch.exe\"; }"

darwin:
	@echo "Building for Darwin..."
	@for arch in $(GOARCHList); do \
		GOOS=darwin GOARCH=$$arch go build -o build/nest-darwin-$$arch; \
		echo "Built build/nest-darwin-$$arch"; \
		done

clean:
	rm -rf build/*



