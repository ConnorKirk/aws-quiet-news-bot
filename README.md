# Quiet AWS News

## What is this?

Staying on top of AWS Announcements is tricky. There are a lot.

This lambda polls the AWS Whats new RSS feed every 24 hours. For anything published in the previous 24 hours, it publishes it in a chime channel.

## How to Use

1. [Create a new chime room](https://docs.aws.amazon.com/chime/latest/ug/chime-chat-room.html)
2. [Add a chime bot. Make note of the Webhook url](https://docs.aws.amazon.com/chime/latest/ug/webhooks.html)
3. Populate the `environment.dummy.json` file with your webhook url. Rename the file to `environment.json`
3. Run `make package`, then `make deploy`


## Environment Variables
`NEWS_OUTPUT_WEBHOOK_URL` - Chime Webhook URL
`NEWS_INPUT_FEED_URL` - Optional - Alternative AWS Whats New RSS Feed URL
`NEWS_TIME_WINDOW_DAYS` - Optional - Window size to filter RSS posts

## Local Development Workflow

Create a `environment.json` file to contain your local environment variables. **Don't check this in** 

* Build with `make build`  
* Start local testing with `make local`  
* Package with `make package`  
* Deploy with `make deploy`

## Deploying

Remember to set an environment variable `NEWS_OUTPUT_WEBHOOK_URL` in the lambda


### What is `deploy.fake`?

`make` works by building a target file if any of the files dependencies have changed. If the target file doesn't exist, then `make` will always run the process. `deploy.fake` acts as the target file for make to build. When `make deploy` or `make deploy.fake` is run, it will only perform the deployment process if any of it's dependencies have changed.


