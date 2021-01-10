# :rocket: boca-chica-bot

>I am a Twitter and Discord Bot that posts status updates on [beach and road closures related to SpaceX
Starship testing][cameron-county-spacex] in Boca Chica, TX.

[![twitter-badge]][@bocachicabot]  
[![discord-invite-badge]][bocachicabot-discord-invite] <sup>see [Discord Installation](#discord-installation) below</sup>

<p align="center">
<img width="500" src="assets/bocachicabot.jpg">
</p>

![go-version-badge]
[![go-report-card-badge]][go-report-card]
[![codecov-badge]][boca-chica-bot-codecov]
[![pkg-go-dev-badge]][pkg-go-dev]

## How I Work

I periodically pull the published road and beach closures from the [Cameron County SpaceX
page][cameron-county-spacex] to see if there are any changes or additions, then post status
updates via Twitter and Discord as [@BocaChicaBot].

Currently I will post a tweet/notification when:

- A closure is added or changed
- A closure has started or ended

I'm written in [Go] and run [serverless] in [AWS] using [AWS Lambda], [DynamoDB], and [EventBridge].

![arch diagram](./assets/boca-chica-bot.drawio.png)

---

## Discord Installation

I am not a typical Discord Bot that responds to commands on demand. **I only post notifications to
the channels in the Discord server you add me to.**

**I will post to *all* channels in the server by default**. This is likely not what you want, so
be sure to [disable the Send Messages
permission][discord-disable-send-messages] for the BocaChicaBot Role in the channels/categories
you'd like to disable me in.

[![discord-invite-badge]][bocachicabot-discord-invite]

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
- `discord_bot_token`

These are used to populate the Parameter Store with the required Twitter and Discord API creds.

## Local Development

Create a `.env` file with these env vars set:

```sh
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_REGION=us-east-1
DEBUG=true
DISABLE_PUBLISH=true
TWITTER_ENVIRONMENT=test
```

The `.env` file is leveraged in the `lambci/lambda:go1.x` Docker container used to run the lambda
handler locally:

```sh
make run
```

---

## Reference

- [AWS SDK for Go][aws-sdk-go]
- [AWS Lambda]
- [DynamoDB]
- [AWS Systems Manager Parameter Store][aws-param-store]
- [Twitter API Docs]
- [Twitter API authentication][twitter-api-auth]

[aws]:https://aws.amazon.com/
[aws lambda]:https://aws.amazon.com/lambda/
[aws-param-store]:https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html
[aws-sdk-go]:https://docs.aws.amazon.com/sdk-for-go/
[cameron-county-spacex]:https://www.cameroncounty.us/spacex/
[codecov-badge]:https://codecov.io/gh/nickshine/boca-chica-bot/branch/master/graph/badge.svg?token=171LQ10HAP
[discord-disable-send-messages]:https://support.discord.com/hc/en-us/articles/206029707-How-do-I-set-up-Permissions-
[discord-invite-badge]:https://img.shields.io/static/v1?label=Discord&logo=Discord&message=Invite%20@BocaChicaBot&colorB=7289DA
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
[bocachicabot-discord-invite]:https://discord.com/api/oauth2/authorize?client_id=782492119063199744&permissions=2048&scope=bot
[boca-chica-bot-codecov]:https://codecov.io/gh/nickshine/boca-chica-bot
