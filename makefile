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

1:
	go build && nest add "../R.png"

dev: stop build open


# TODO: rewrite build-pub in this;
# BUILD-PUB

GOARCHList := amd64 arm64

#@$(PWSH) "foreach (\$$arch in @('amd64','arm64')) { \$$env:GOOS = 'windows'; \$$env:GOARCH = \$$arch; go build -ldflags '-H=windowsgui' -o build/nest-windows-\$$arch.exe; Write-Host \"Built build/nest-windows-\$$arch.exe\" }"

windows:
	@echo "Building for Windows..."

	# HACK: ew
	@$(PWSH) "foreach (\$$arch in @('amd64','arm64')) { \
		\$$env:GOOS = 'windows'; \
		\$$env:GOARCH = \$$arch; \
		go build -o build/nest-windows-\$$arch.exe; \
		Write-Host \"Built build/nest-windows-\$$arch.exe\" \
	}"

# TODO: fix darwin and add more tests...
darwin:
	@echo "Building for Darwin..."
	@for arch in $(GOARCHList); do \
		GOOS=darwin GOARCH=$$arch go build -o build/nest-darwin-$$arch; \
		echo "Built build/nest-darwin-$$arch"; \
		done

clean:
	rm -rf build/*


## all commits since last tag
# git log "$(git describe --tags --abbrev=0)..HEAD" --oneline
#
# gh release create v1.0.0 yourfile.zip -t "Release 1.0.0" -n "Release notes"
#  git tag -a v0.0.7 -m 'improve error handling for nest switch' 4db8bc0 && git push --tags

# tag last commit (local?)
# git tag -a <VERSION>
#
# tag last commit (local?)
# git tag -a <VERSION> -m <message> -m <description>
#
# git tag -a <VERSION> -m <message> -m (git log "$(git describe --tags --abbrev=0)..HEAD" --oneline)
#
