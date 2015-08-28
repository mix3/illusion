GIT_VER := $(shell git describe --tags)
DATE    := $(shell date +%Y-%m-%dT%H:%M:%S%z)
OWNER   := "mix3"
REPO    := "illusion"

clean:
	rm -rf pkg/*

binary: clean
	env CGO_ENABLED=0 gox -osarch="linux/amd64 darwin/amd64" \
		-output "pkg/{{.Dir}}-${GIT_VER}-{{.OS}}-{{.Arch}}" \
		-ldflags "-X main.version ${GIT_VER} -X main.buildDate ${DATE}"

package: binary
	cd ./pkg && find . -name "*${GIT_VER}*" -type f \
		-exec mkdir -p illusion \; \
		-exec cp {} illusion/illusion \; \
		-exec cp -r ../config.toml.sample illusion/ \; \
		-exec zip -r {}.zip illusion \; \
		-exec rm -rf illusion \;
