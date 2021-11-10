test:
	go test -cover ./...

install: test
	go install 

run: install
	checkip 91.228.166.47
	checkip 209.141.33.65
	checkip 218.92.0.158

test-api:
	curl -s 'localhost:8000/api/v1/91.228.166.47' | jq -r '.ProbabilityMalicious'
	curl -s 'localhost:8000/api/v1/209.141.33.65' | jq -r '.ProbabilityMalicious'
	curl -s 'localhost:8000/api/v1/218.92.0.158' | jq -r '.ProbabilityMalicious'

PLATFORMS := linux/amd64 darwin/amd64 linux/arm windows/amd64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

release: $(PLATFORMS)

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -ldflags "-w" -o cmd/checkip-$(os)-$(arch) cmd/checkip.go
	tar -cf - cmd/checkip-$(os)-$(arch) | gzip -9c > cmd/checkip-$(os)-$(arch).tar.gz
	rm -f cmd/checkip-$(os)-$(arch)