# GNU Make
# See also https://blog.gopheracademy.com/advent-2017/make/

GO11MODULES = on
APP := server
MAINMODULE := main.go

# To generate a production-ready executable
# simply invoke "make NODEBUG=1"
ifdef NODEBUG
	LDFLAGS := -ldflags="-s -w"
	BINARYEXT := -nodebug
else
	BINARYEXT := -debug
endif

.PHONY: compile
compile: dep
	go mod download
	go build ${LDFLAGS} -o ${APP}${BINARYEXT} ${MAINMODULE}

.PHONY: build
build: clean compile

.PHONY: depmod
depmod:
	$(MAKE) -C models dep

.PHONY: depcont
depcont:
	$(MAKE) -C controller dep

.PHONY: dep
dep: depmod depcont

.PHONY: run
run: compile
	./server -conf=config/config.json

.PHONY: clean
clean:
	@echo "Cleaning"
	$(MAKE) -C models clean
	$(MAKE) -C controller clean
	-rm ${APP}-debug ${APP}-nodebug

.PHONY: test
test: dep
	go test -v -count=1 -race ./...

#.PHONY: setup
### setup: setup go modules
#setup:
#	@go mod init
#	@go mod tidy
#	@go mod vendor
