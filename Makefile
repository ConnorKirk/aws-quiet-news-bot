.PHONY:  clean build


clean: 
	rm -rf ./news/news
	
build:
	GOOS=linux GOARCH=amd64 go build -o news/news ./news