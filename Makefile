news := bin/news
package := packaged.yaml

local:
	sam local start-api --env-vars environment.json

.PHONY: build
build : bin/news
SOURCES := $(wildcard news/*.go)
bin/news: $(SOURCES)
	GOOS=linux GOARCH=amd64 go build -o bin/news ./news

.PHONY:  package
package: packaged.yaml
packaged.yaml: $(news) template.yaml
	sam package --template-file template.yaml \
	--s3-bucket quiet-aws-news-lambda \
	--output-template-file packaged.yaml

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
