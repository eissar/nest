PWSH := pwsh.exe -NoProfile -Command


test:
	write-host 'test'

build:
	go build

stop:
	./nest.exe -stop

open:
	wt.exe -w 0 nt -d . $(PWSH) ./nest.exe -serve


docs:
	swag init -d .\core\,.\handlers\,.\plugins\nest\ -g .\core.go --parseInternal --parseFuncBody --parseDependency


dev: stop open


# TODO: rewrite build-pub in this;


