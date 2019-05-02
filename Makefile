.PHONY:  clean build

local:
	sam local start-api --env-vars environment.json
clean: 
	rm -rf ./news/news
	
build:
	GOOS=linux GOARCH=amd64 go build -o bin/news ./news

package:
	sam package --template-file template.yaml \
	 --s3-bucket quiet-aws-news-lambda \
	 --output-template-file packaged.yaml

deploy:
	aws cloudformation deploy \
		--template-file /Users/ckp/go/src/github.com/ConnorKirk/aws-quiet-news-bot/packaged.yaml \
		--stack-name quiet-aws-news \
		--capabilities CAPABILITY_IAM \
		--profile $(AWS_PROFILE) \
		--env-vars environment.json