SPEC_FILE 				:= docs/openapi.yaml
GEN_DIR 				:= ./api
OUTPUT_PATH             := bin/api
MAIN_MODULE             := cmd/main.go
IMAGE_VERSION           := latest

OAPI_GENERATOR_VERSION 	?= latest

##  
## Commands
##

.PHONY: run
run: compile
	./${OUTPUT_PATH}


.PHONY: test
test:
	go test -v -count=1 ./...


.PHONY: compile
compile:
	go build -o ${OUTPUT_PATH} ${MAIN_MODULE}


.PHONY: tools
tools:
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_GENERATOR_VERSION)


.PHONY: oapi-codegen
oapi-codegen: tools
	@mkdir -p $(GEN_DIR)
	oapi-codegen --package=api -generate="gin" -o $(GEN_DIR)/server.gen.go $(SPEC_FILE)