## help: Print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## gen: Generate sdk, gateway,... for protobuffer
.PHONY: gen
gen: # lint # break
	buf generate proto
	for service in "inventory"; do \
	cd golang/$$service && \
		(go mod init github.com/ebisaan/proto/golang/$$service || true) && \
		go mod tidy; \
	done

## break: Check if proto files break backward-compatibility
break:
	buf breaking proto --against "https://github.com/ebisaan/proto.git#branch=main,subdir=proto"

## lint: Lint proto files
lint:
	buf lint proto


## run/gateway: Run gRPC-Gateway params=$1
.PHONY: run/gateway
run/gateway:
	cd gateway; go run ./cmd/gateway ${param}
