SHELL				?= /bin/bash

# Update these for your dev env
ONEPASSWORD_VAULT	?= dev-local
ONEPASSWORD_SECRET  ?= jumpcloud_techjavelin_oss
ONEPASSWORD_FIELD	?= credential
ONEPASSWORD_SIGNIN	?=

# These can be overridden at any time on the env
VERIFY_LOG_PRIORITY ?= INFO
VERIFY_LOG_FILE		?= $(VERIFY_PATH)/terraform-verify.log
VERIFY_RESOURCE     ?=

# Only modify these if you know wtf you're doing
TERRAFORM_BIN		?= $(shell which terraform)
VERIFY_PATH			?= verify/$(VERIFY_RESOURCE)

provider_registry	:= registry.terraform.io
provider_group 		:= techjavelin
provider_name  		:= jumpcloud
provider_version	:= 0.0.1
provider_executable := terraform-provider-$(provider_name)

# Don't modify these
build_os 		   	:= $(shell go env GOOS)
build_arch 		   	:= $(shell go env GOARCH)
plugin_path			:=  ~/.terraform.d/plugins

# Debugging Stuff
PROVIDER_NETWORK	?=
PROVIDER_PID		?=
debug_flags			:= -gcflags="all=-N -l"
debug_attach		:= 	'{"registry.terraform.io/techjavelin/jumpcloud":{"Protocol":"grpc","ProtocolVersion":6,"Pid":$(PROVIDER_PID),"Test":true,"Addr":{"Network":"unix","String":"/tmp/$(PROVIDER_NETWORK)"}}}'


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

verify: install verify/cleanup
	rm -f $(VERIFY_LOG_FILE)
	make verify/create
	make verify/import
	make verify/update
	make verify/destroy
	make verify/cleanup

verify/nodestroy: install verify/cleanup
	make verify/add
	make verify/import

verify/init:
	@echo "-=-=-=-=-=-=] init [=-=-=-=-=-=-" 
	@$(TERRAFORM_BIN) -chdir=$(VERIFY_PATH) init 

verify/prepare:
	@echo "-=-=-=-=-=-=] prepare [=-=-=-=-=-=-" 
	rm -f $(VERIFY_PATH)/terraform.tfplan
	@export TF_VAR_jumpcloud_api_key=$(jumpcloud_api_key) && \
	export TF_LOG=$(VERIFY_LOG_PRIORITY) && \
	export TF_LOG_PATH=$(VERIFY_LOG_FILE) && \
	op run -- $(TERRAFORM_BIN) -chdir=$(VERIFY_PATH) plan -out=terraform.tfplan 

verify/create: 
	@echo "-=-=-=-=-=-=] create [=-=-=-=-=-=-" && \
	cp $(VERIFY_PATH)/main.create $(VERIFY_PATH)/main.tf
	@make verify/init verify/prepare verify/apply

verify/import: 
	@export TF_VAR_jumpcloud_api_key=$(jumpcloud_api_key) && \
	export TF_LOG=$(VERIFY_LOG_PRIORITY) && \
	export TF_LOG_PATH=$(VERIFY_LOG_FILE) && \
	echo "-=-=-=-=-=-=] import [=-=-=-=-=-=-" && \
	export RESOURCE_ID=`$(TERRAFORM_BIN) -chdir=$(VERIFY_PATH) state show "$(VERIFY_RESOURCE).test" | grep id | head -n1 | grep id | cut -d'=' -f2 | cut -d'"' -f2` && \
	op run -- $(TERRAFORM_BIN) -chdir=$(VERIFY_PATH) state rm "$(VERIFY_RESOURCE).test" && \
	echo "Resource ID: $$RESOURCE_ID" && \
	op run -- $(TERRAFORM_BIN) -chdir=$(VERIFY_PATH) import "$(VERIFY_RESOURCE).test" $$RESOURCE_ID 

verify/update:
	@echo "-=-=-=-=-=-=] update [=-=-=-=-=-=-" 
	cp $(VERIFY_PATH)/main.update $(VERIFY_PATH)/main.tf
	@make verify/init verify/prepare verify/apply

verify/apply:
	@echo "-=-=-=-=-=-=] apply [=-=-=-=-=-=-"
	@export TF_VAR_jumpcloud_api_key=$(jumpcloud_api_key) && \
	export TF_LOG=$(VERIFY_LOG_PRIORITY) && \
	export TF_LOG_PATH=$(VERIFY_LOG_FILE) && \
	op run -- $(TERRAFORM_BIN) -chdir=$(VERIFY_PATH) apply --auto-approve terraform.tfplan

verify/destroy:
	@echo "-=-=-=-=-=-=] destroy [=-=-=-=-=-=-"
	@export TF_VAR_jumpcloud_api_key=$(jumpcloud_api_key) && \
	export TF_LOG=$(VERIFY_LOG_PRIORITY) && \
	export TF_LOG_PATH=$(VERIFY_LOG_FILE) && \
	op run -- $(TERRAFORM_BIN) -chdir=$(VERIFY_PATH) destroy --auto-approve 
	@make verify/cleanup

verify/cleanup:
	@echo "-=-=-=-=-=-=] cleanup [=-=-=-=-=-=-"
	rm -rf $(VERIFY_PATH)/.terraform $(VERIFY_PATH)/.terraform.lock.hcl $(VERIFY_PATH)/terraform.tfstate $(VERIFY_PATH)/terraform.tfstate* $(VERIFY_PATH)/main.tf $(VERIFY_LOG_FILE) $(VERIFY_PATH)terraform.tfplan

build/debug:
	go build $(debug_flags)

