gen: lint generate
	buf generate proto
	for service in "inventory"; do \
	cd golang/$$service && \
		(go mod init github.com/ebisaan/proto/golang/$$service || true) && \
		go mod tidy; \
	done

break:
	buf breaking proto --against "https://github.com/ebisaan/proto.git#branch=main,subdir=proto"
lint:
	buf lint proto
