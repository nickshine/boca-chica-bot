variable "twitter_consumer_key" {}    // provided in Terraform Cloud workspace
variable "twitter_consumer_secret" {} // provided in Terraform Cloud workspace
variable "twitter_access_token" {}    // provided in Terraform Cloud workspace
variable "twitter_access_secret" {}   // provided in Terraform Cloud workspace

locals {
  params = {
    "consumer_key"    = var.twitter_consumer_key
    "consumer_secret" = var.twitter_consumer_secret
    "access_token"    = var.twitter_access_token
    "access_secret"   = var.twitter_access_secret
  }
}

resource "aws_ssm_parameter" "twitter_creds" {
  for_each = local.params

  name  = "${var.param_store_path}/twitter_${each.key}"
  type  = "SecureString"
  value = each.value

  tags = local.tags
}
