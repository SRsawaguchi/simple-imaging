VERSION := 1.0.0
NAME := srimage
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

ifeq ($(GOOS), windows)
dist/$(NAME)-$(VERSION)_$(GOOS)_$(GOARCH).zip:
	go build -ldflags="X main.version=$(VERSION)" -o dist/cli/$(NAME)-$(VERSION)/$(NAME).exe cmd/cli/main.go
	zip -j dist/cli/$(NAME)-$(VERSION)_$(GOOS)_$(GOARCH).zip dist/cli/$(NAME)-$(VERSION)/$(NAME).exe
else
dist/$(NAME)-$(VERSION)_$(GOOS)_$(GOARCH).tar.gz:
	go build -ldflags="-X main.version=$(VERSION)" -o dist/cli/$(NAME)-$(VERSION)/$(NAME) cmd/cli/main.go
	tar cfz dist/cli/$(NAME)-$(VERSION)_$(GOOS)_$(GOARCH).tar.gz -C dist/cli/$(NAME)-$(VERSION) $(NAME)
endif

.PHONY: clean
clean:
	rm -rf dist

.PHONY: test
test:
	go test github.com/SRsawaguchi/srimage/imaging