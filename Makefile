gen: lint test
	buf generate proto
	for service in "inventory"; do \
	cd golang/$$service/v1 && \
		(go mod init github.com/ebisaan/proto/golang/$$service || true) && \
		go mod tidy; \
	done

test:
	buf breaking proto --against "https://github.com/ebisaan/proto.git#branch=main,subdir=proto"
lint:
	buf lint proto
