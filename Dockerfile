FROM docker:dind

ENV PATH /go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	bash \
	ca-certificates \
	curl \
	make \
	gcc \
	go \
	git \
	libc-dev \
	libgcc \
	make \
	diffutils \
	jq \
	file

RUN go get github.com/golang/lint/golint \
	&& go get golang.org/x/tools/cmd/cover \
	&& go install cmd/vet cmd/cover

WORKDIR /go/src/github.com/dcos/dcos-checks
