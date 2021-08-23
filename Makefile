
PLATFORMS=darwin linux windows
ARCHITECTURES=amd64

SRC := pwgen.go

default: release

define GEN_RULE
bin/release/$(goos)/$(goarch)/%: $(SRC)
	GOOS=$(goos) GOARCH=$(goarch) go build -o $$@ .
endef

$(foreach goos,$(PLATFORMS), \
  $(foreach goarch,$(ARCHITECTURES), \
    $(eval $(GEN_RULE)) \
  ) \
)

bin/release.tgz: \
  bin/release/linux/amd64/pwgen \
  bin/release/darwin/amd64/pwgen \
  bin/release/windows/amd64/pwgen.exe
	tar -czf  bin/release.tgz bin/release

release: bin/release.tgz
	
clean:
	rm -rf bin

.PHONY = clean release
