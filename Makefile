
PLATFORMS=darwin linux windows
ARCHITECTURES=amd64

pwgen:
	go build -o bin/pwgen

release:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -o bin/release/$(GOOS)/$(GOARCH)/pwgen)))
	tar -cvzf  bin/release.tgz bin/release
	

.PHONY = pwgen release
