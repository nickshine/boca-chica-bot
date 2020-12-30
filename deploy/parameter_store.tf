variable "twitter_consumer_key" {}    // provided in Terraform Cloud workspace
variable "twitter_consumer_secret" {} // provided in Terraform Cloud workspace
variable "twitter_access_token" {}    // provided in Terraform Cloud workspace
variable "twitter_access_secret" {}   // provided in Terraform Cloud workspace
variable "discord_bot_token" {}       // provided in Terraform Cloud workspace

locals {
  params = {
    "twitter_consumer_key"    = var.twitter_consumer_key
    "twitter_consumer_secret" = var.twitter_consumer_secret
    "twitter_access_token"    = var.twitter_access_token
    "twitter_access_secret"   = var.twitter_access_secret
    "discord_bot_token"       = var.discord_bot_token
  }
}

resource "aws_ssm_parameter" "twitter_creds" {
  for_each = local.params

  name  = "${var.param_store_path}/${each.key}"
  type  = "SecureString"
  value = each.value

  tags = local.tags
}
