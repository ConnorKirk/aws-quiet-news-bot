news := bin/news
package := packaged.yaml

# Start the local development api with local environment
# variables
local:
	sam local start-api --env-vars environment.json

# Build the go binary
.PHONY: build
build : bin/news
SOURCES := $(wildcard news/*.go)
bin/news: $(SOURCES)
	GOOS=linux GOARCH=amd64 go build -o bin/news ./news

# Package the binary and CF template
.PHONY:  package
package: packaged.yaml
packaged.yaml: $(news) template.yaml
	sam package --template-file template.yaml \
	--s3-bucket quiet-aws-news-lambda \
	--output-template-file packaged.yaml

# Deploy to lambda
.PHONY: deploy
deploy: deploy.fake
deploy.fake: packaged.yaml
	aws cloudformation deploy \
		--template-file /Users/ckp/go/src/github.com/ConnorKirk/aws-quiet-news-bot/packaged.yaml \
		--stack-name quiet-aws-news \
		--capabilities CAPABILITY_IAM
	touch deploy.fake

.PHONY: clean
clean: 
	rm -rf ./bin/news
	rm $(package)
	rm deploy.fake
