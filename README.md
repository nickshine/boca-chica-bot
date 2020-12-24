# :rocket: boca-chica-bot

>I am a Twitter Bot that tweets status updates to [beach and road closures related to SpaceX
Starship testing][cameron-county-spacex] in Boca Chica, TX.

[![twitter-badge]][@bocachicabot]

<p align="center">
<img width="500" src="assets/boca-chica-bot.jpg">
</p>

![go-version-badge]
[![go-report-card-badge]][go-report-card]
[![pkg-go-dev-badge]][pkg-go-dev]

## How I Work

I periodically pull the published road and beach closures from the [Cameron County SpaceX
page][cameron-county-spacex] to see if there are any changes or additions, then tweet them out as
[@BocaChicaBot].

I'm written in [Go] and run [serverless] in [AWS] using [AWS Lambda], [DynamoDB], and [EventBridge].

![arch diagram](./assets/boca-chica-bot.drawio.svg)

---

## Deployment

[Terraform Cloud] is utilized for deploying my infrastructure. I have two workspaces (access required to see these):

- [boca-chica-bot-test][terraform-cloud-workspace-test]
- [boca-chica-bot-prod][terraform-cloud-workspace-prod]

### Terraform Cloud Workspace Setup

Required Environment Variables for these workspaces in Terraform Cloud:

- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `TF_CLI_ARGS_plan=-var-file=workspaces/test.tfvars` (for test)
- `TF_CLI_ARGS_plan=-var-file=workspaces/prod.tfvars` (for prod)

Required Terraform Variables in Terraform Cloud:

- `twitter_consumer_key`
- `twitter_consumer_secret`
- `twitter_access_secret`
- `twitter_access_token`

These are used to populate the Parameter Store with the required Twitter API creds.

## Local Development

Create a `.env` file with these env vars set:

```sh
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_REGION=us-east-1
DEBUG=true
DISABLE_TWEETS=true
TWITTER_ENVIRONMENT=test
```

The `.env` file is leveraged in the `lambci/lambda:go1.x` Docker container used to run the lambda
handler locally:

```sh
make run
```

---

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
[go-report-card]:https://goreportcard.com/report/github.com/nickshine/boca-chica-bot
[go-report-card-badge]:https://goreportcard.com/badge/github.com/nickshine/boca-chica-bot
[go-version-badge]:https://img.shields.io/github/go-mod/go-version/nickshine/boca-chica-bot
[pkg-go-dev]:https://pkg.go.dev/github.com/nickshine/boca-chica-bot
[pkg-go-dev-badge]:https://pkg.go.dev/badge/github.com/nickshine/boca-chica-bot
[serverless]:https://aws.amazon.com/serverless/
[terraform cloud]:https://www.hashicorp.com/products/terraform
[terraform-cloud-workspace-test]:https://app.terraform.io/app/nickshine/workspaces/boca-chica-bot-test
[terraform-cloud-workspace-prod]:https://app.terraform.io/app/nickshine/workspaces/boca-chica-bot-prod
[twitter api docs]:https://developer.twitter.com/en/docs/twitter-api
[twitter-api-auth]:https://developer.twitter.com/en/docs/authentication/overview
[twitter-badge]:https://img.shields.io/twitter/follow/BocaChicaBot?style=social
[@BocaChicaBot]:https://twitter.com/bocachicabot
