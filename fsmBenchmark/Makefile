
.PHONY: all
all: bench
	mkdir -p ./testDataTmp
	mkdir -p ./logDatas

outBin:
	go build -o $@ .

.PHONY: test
test:
	@echo "Running test..."
	mkdir -p ./logDatas
	go test .
#	go test -race $(PKGS)

.PHONY: cover
cover:
	@echo "Running coverage testing..."
	mkdir -p ./logDatas
	./scripts/cover.sh .

.PHONY: bench
BENCH ?= .
bench:
#	@$(foreach pkg,$(PKGS),go test -bench=$(BENCH) -run="^$$" $(BENCH_FLAGS) $(pkg);)
	@echo "Running benchmark..."
	mkdir -p ./logDatas
	go test -bench=. -run=none -benchtime=1s -benchmem -cpuprofile cpu.prof -memprofile mem.prof
	go tool pprof -pdf mem.prof
	go tool pprof -pdf cpu.prof

.PHONY: updatereadme
updatereadme:
	@echo "README.md updating..."
#	rm -f README.md
#	cat .readme.tmpl | go run internal/readme/readme.go > README.md
#
#observer
.PHONY:docker
docker: outBin
#	docker rmi -f myapp:falcon
	docker build --rm -t myapp:falcon \
				--build-arg "group=${USER}" \
				--build-arg "gid=$$(id -u)" \
				--build-arg "user=${USER}" \
				--build-arg	"uid=$$(id -g)" \
				-f docker/Dockerfile.img .

.PHONY: clean
clean:
	rm -rf outBin cover.out cover.html cpu.prof mem.prof profile*.pdf
	rm -rf ./testDataTmp
	rm -rf logDatas
	rm -rf *.test
	go clean -cache .
#	docker rmi -f myapp:falcon


