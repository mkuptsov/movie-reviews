OS ?= $(shell go env GOOS)
ARCH ?= $(shell  go env GOARCH)
VERSION ?= 0.1
PROVIDER_PATH = "/home/vagrant/.terraform.d/plugins/localhost/providers/moviereviews/$(VERSION)/$(OS)_$(ARCH)/terraform-provider-moviereviews"

install-provider:
	go build -o $(PROVIDER_PATH) .
