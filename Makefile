default: release

bin/release/linux/amd64/pwgen: pwgen.go
	GOOS=linux GOARCH=amd64 go build -o $@ $<

bin/release/linux/amd64/pwgen_hard: pwgen_hard.go
	GOOS=linux GOARCH=amd64 go build -o $@ $<

bin/release/darwin/amd64/pwgen: pwgen.go
	GOOS=darwin GOARCH=amd64 go build -o $@ $<

bin/release/darwin/amd64/pwgen_hard: pwgen_hard.go
	GOOS=darwin GOARCH=amd64 go build -o $@ $<

bin/release/windows/amd64/pwgen.exe: pwgen.go
	GOOS=windows GOARCH=amd64 go build -o $@ $<

bin/release/windows/amd64/pwgen_hard.exe: pwgen_hard.go
	GOOS=windows GOARCH=amd64 go build -o $@ $<

bin/release.tgz: \
  bin/release/linux/amd64/pwgen \
  bin/release/darwin/amd64/pwgen \
  bin/release/windows/amd64/pwgen.exe \
  bin/release/linux/amd64/pwgen_hard \
  bin/release/darwin/amd64/pwgen_hard \
  bin/release/windows/amd64/pwgen_hard.exe
	tar -czf  bin/release.tgz bin/release

release: bin/release.tgz
	
clean:
	rm -rf bin

.PHONY = clean release
