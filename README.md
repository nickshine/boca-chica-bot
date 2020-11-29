# boca-chica-bot :rocket:

<p align="center">
<img width="800" src="assets/boca-chica-bot.jpg">
</p>

>I am a Twitter Bot that tweets status updates to [beach and road closures related to SpaceX
Starship testing][cameron-county-spacex] in Boca Chica, TX.

[![Twitter](https://img.shields.io/twitter/follow/BocaChicaBot?style=social)][@bocachicabot]

## How I Work

I periodically pull the published road and beach closures from the [Cameron County SpaceX
page][cameron-county-spacex] to see if there are any changes or additions, then tweet them out as
[@BocaChicaBot].

I'm written in [Go] and run [serverless] in [AWS] using [AWS Lambda], [DynamoDB], and [EventBridge].

---

## Development

### Environment

Required env vars for the bot to execute:

```sh
TWITTER_CONSUMER_KEY
TWITTER_CONSUMER_SECRET
TWITTER_ACCESS_TOKEN
TWITTER_ACCESS_SECRET

AWS_REGION=us-east-1
```

Set this env var for verbose logging during development:

```sh
DEBUG=true
```

## Reference

* [AWS SDK for Go][aws-sdk-go]
* [AWS Lambda]
* [DynamoDB]
* [AWS Systems Manager Parameter Store][aws-param-store]
* [Twitter API Docs]
* [Twitter API authentication][twitter-api-auth]

[aws]:https://aws.amazon.com/
[aws lambda]:https://aws.amazon.com/lambda/
[aws-param-store]:https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html
[aws-sdk-go]:https://docs.aws.amazon.com/sdk-for-go/
[cameron-county-spacex]:https://www.cameroncounty.us/spacex/
[dynamodb]:https://aws.amazon.com/dynamodb/
[EventBridge]:https://aws.amazon.com/eventbridge/
[go]:https://golang.org/
[serverless]:https://aws.amazon.com/serverless/
[twitter api docs]:https://developer.twitter.com/en/docs/twitter-api
[twitter-api-auth]:https://developer.twitter.com/en/docs/authentication/overview
[@BocaChicaBot]:https://twitter.com/bocachicabot
