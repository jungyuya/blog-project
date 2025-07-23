# blog-backend/Makefile

SHELL := /bin/bash
.SHELLFLAGS := -e -c

SERVICE_NAME     = blog-api-dev
STAGE            = dev
REGION           = ap-northeast-2
LAMBDA_FUNCTIONS = getPosts getPost createPost updatePost deletePost
ARTIFACT_DIR     = artifacts

.PHONY: all build clean deploy

all: build

build: $(ARTIFACT_DIR)
	@echo "Building Go Lambdas into $(ARTIFACT_DIR)/ with CGO disabled"
	@for func in $(LAMBDA_FUNCTIONS); do \
		echo "  -> $$func"; \
		cd $$func && \
		  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		    go build -ldflags="-s -w" -o bootstrap ./main.go && \
		  chmod +x bootstrap && \
		  zip -r ../$(ARTIFACT_DIR)/$$func.zip bootstrap && \
		  rm bootstrap; \
		cd ..; \
	done
	@echo "✔ Artifacts ready in $(ARTIFACT_DIR)/"

clean:
	@echo "Cleaning up artifacts..."
	@rm -rf $(ARTIFACT_DIR)
	@find . -type f -name "bootstrap" -delete
	@echo "✔ Clean complete."

deploy: clean build
	@echo "Deploying service to AWS..."
	sls deploy --stage $(STAGE) --region $(REGION)
	@echo "✔ Deployment complete."

$(ARTIFACT_DIR):
	mkdir -p $(ARTIFACT_DIR)
