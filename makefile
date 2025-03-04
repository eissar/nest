PWSH := pwsh.exe -NoProfile -Command


test:
	go test

.PHONY: build
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


