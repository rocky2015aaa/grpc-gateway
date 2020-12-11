# GNU Make

OS := $(shell uname)
PBMODFILES := $(wildcard *.proto)
PBMODTARGETS := $(PBMODFILES:.proto=.pb.go)
PBMODGEN :=  protoc -I${GOPATH}/third_party/googleapis/ -I. -I=${GOPATH}/src -I=${GOPATH}/protobuf --gogo_out=plugins=grpc,Mgoogle/protobuf/timestamp.proto=github.com/types:.
UNAME_S := $(shell uname -s)

.PHONY: build
build: clean dep

.PHONY: dep
dep: ${PBMODTARGETS}

%.pb.go: %.proto
	${PBMODGEN} $<
	protoc-go-inject-tag -input=$@
ifeq ($(UNAME_S),Linux)
	sed "" s/json=.*,proto3/proto3/ $@
else
	sed s/json=.*,proto3/proto3/ $@ > mk_sed_tmp.tmp
	mv mk_sed_tmp.tmp $@
endif

.PHONY: clean
clean:
	@-rm -f ${PBMODTARGETS}
