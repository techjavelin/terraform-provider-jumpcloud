# Update these for your dev env
ONEPASSWORD_VAULT	?= techjavelin-automation
ONEPASSWORD_SECRET  ?= jumpcloud_cschmidt.techjavelin.com
ONEPASSWORD_FIELD	?= credential

# Only modify these if you know wtf you're doing
TERRAFORM_BIN		?= $(shell which terraform)
VERIFY_PATH			?= examples/provider-install-verification

provider_registry	:= registry.terraform.io
provider_group 		:= techjavelin
provider_name  		:= jumpcloud
provider_version	:= 0.0.1
provider_executable := terraform-provider-$(provider_name)

# Don't modify these
build_os 		   	:= $(shell go env GOOS)
build_arch 		   	:= $(shell go env GOARCH)
plugin_path			:=  ~/.terraform.d/plugins

# ifeq($(build_os),windows)
# 	SHELL 		   		:= pwsh -NoProfile
# 	plugin_path			:= ${Env:APPDATA}/terraform.d/plugins
# 	provider_executable := $(provider_executable).exe
# endif

plugin_install_path ?= $(plugin_path)/$(provider_registry)/$(provider_group)/$(provider_name)/$(provider_version)/$(build_os)_$(build_arch)
jumpcloud_api_key   ?= op://$(ONEPASSWORD_VAULT)/$(ONEPASSWORD_SECRET)/$(ONEPASSWORD_FIELD)

default: testacc

# Run acceptance tests
.PHONY: testacc install format dependencies build verify_cleanup

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

format:
	go fmt

dependencies:
	go mod tidy

build: format dependencies
	go build -o $(provider_executable)

install: build
	mkdir -p $(plugin_install_path)
	cp $(provider_executable) $(plugin_install_path)

verify_prepare:
	export TF_VAR_jumpcloud_api_key=$(jumpcloud_api_key) && \
	export TF_LOG=INFO && \
	op run -- $(TERRAFORM_BIN) -chdir=$(VERIFY_PATH) state refresh 
verify: install verify_cleanup
	export TF_VAR_jumpcloud_api_key=$(jumpcloud_api_key) && \
	export TF_LOG=INFO && \
	op run -- $(TERRAFORM_BIN) -chdir=$(VERIFY_PATH) plan -out=/tmp/tfplan && \
	op run -- $(TERRAFORM_BIN) -chdir=$(VERIFY_PATH) apply --auto-approve /tmp/tfplan && \
	op run -- $(TERRAFORM_BIN) -chdir=$(VERIFY_PATH) destroy --auto-approve

verify_cleanup:
	rm -rf $(VERIFY_PATH)/.terraform 
	rm -f $(VERIFY_PATH)/.terraform.lock.hcl $(VERIFY_PATH)/terraform.tfstate $(VERIFY_PATH)/terraform.tfstate.backup