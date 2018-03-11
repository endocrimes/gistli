BUILDDIR := ${PWD}/cross
NAME := gistli
GOOSARCHES = darwin/amd64 darwin/386 freebsd/amd64 freebsd/386 linux/arm linux/arm64 linux/amd64 linux/386 solaris/amd64 windows/amd64 windows/386

all: clean fmt lint vet cross

.PHONY: fmt
fmt:
	@echo "+ $@"
	@gofmt -s -l . | grep -v vendor | tee /dev/stderr

.PHONY: lint
lint:
	@echo "+ $@"
	@golint ./... | grep -v vendor | tee /dev/stderr

.PHONY: vet
vet:
	@echo "+ $@"
	@go vet $(shell go list ./... | grep -v vendor) | tee /dev/stderr

.PHONY: install
install:
	@echo "+ $@"
	@go install .

define build
mkdir -p $(BUILDDIR)/$(1)/$(2);
GOOS=$(1) GOARCH=$(2) CGO_ENABLED=0 go build \
	 -o $(BUILDDIR)/$(1)/$(2)/$(NAME);
md5sum $(BUILDDIR)/$(1)/$(2)/$(NAME) > $(BUILDDIR)/$(1)/$(2)/$(NAME).md5;
sha256sum $(BUILDDIR)/$(1)/$(2)/$(NAME) > $(BUILDDIR)/$(1)/$(2)/$(NAME).sha256;
endef

.PHONY: cross
cross: *.go
	@echo "+ $@"
	$(foreach GOOSARCH,$(GOOSARCHES), $(call build,$(subst /,,$(dir $(GOOSARCH))),$(notdir $(GOOSARCH))))

.PHONY: clean
clean:
	@echo "+ $@"
	$(RM) $(NAME)
	$(RM) -r $(BUILDDIR)
