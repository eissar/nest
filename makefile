#COMSPEC := pwsh.exe -noprofile
all:
	$(info SHELL is "$(SHELL)")
#Write-Host 'test'

test:
	write-host 'test'


stop:
	$$ErrorActionPreference = "silentlyContinue"; iwr "http://localhost:1323/api/ping" -ConnectionTimeoutSeconds 1 | select-object StatusCode, Content &&\ iwr "http://localhost:1323/api/server/close" -ConnectionTimeoutSeconds 10 |\ Select StatusCode, Content

## build: build the application
build:
	go build

dev:
	go build && wt.exe -w 0 nt -d . nest.exe -serve


