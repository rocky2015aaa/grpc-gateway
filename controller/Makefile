# GNU Make

PBCONTFILES := $(wildcard *.proto)
PBCONTTARGETS := $(PBCONTFILES:.proto=.pb.go)
PBGWCONTTARGETS := $(PBCONTFILES:.proto=.pb.gw.go)
PBCONTGEN := protoc -I${GOPATH}/googleapis/ -I. -I=${GOPATH}/src -I=${GOPATH}/protobuf --gogo_out=plugins=grpc,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types:. --grpc-gateway_out=allow_patch_feature=false,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api:.
UNAME_S := $(shell uname -s)

.PHONY: build
build: clean dep

.PHONY: dep
dep: ${PBCONTTARGETS}

%.pb.go: %.proto
	${PBCONTGEN} $<
	protoc-go-inject-tag -input=$@
	ls *.pb.go | xargs -n1 -IX bash -c 'sed s/json=.*,proto3/proto3/ X > X.tmp && mv X{.tmp,}'
ifeq ($(UNAME_S),Linux)
	sed "" s/json=.*,proto3/proto3/ $@
endif

.PHONY: clean
clean:
	@-rm -f ${PBCONTTARGETS}
	@-rm -f ${PBGWCONTTARGETS}
