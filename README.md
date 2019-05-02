# Quiet AWS News

## What is this?

Staying on top of AWS Announcements is tricky. There are a lot.

This lambda polls the AWS Whats new RSS feed every 24 hours. For anything published in the previous 24 hours, it publishes it in a chime channel.

The fancy part is that you can (attempt) to specify only the regions you are interested in.

## Environment Variables
webhook: Webhook URL to post to

## Local Workflow

Create a local `environment.json` file to contain your local environment variables. Don't check this in

Build with `make build`
Start local testing with `make local`
Package with `make package`
Deploy with `make deploy`

## Deploying

Remember to set an environment variable `webhook`


