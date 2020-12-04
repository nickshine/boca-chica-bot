variable "twitter_consumer_key" {}    // provided in Terraform Cloud workspace
variable "twitter_consumer_secret" {} // provided in Terraform Cloud workspace
variable "twitter_access_token" {}    // provided in Terraform Cloud workspace
variable "twitter_access_secret" {}   // provided in Terraform Cloud workspace

resource "aws_ssm_parameter" "twitter_consumer_key" {
  name  = "${var.param_store_path}/twitter_consumer_key"
  type  = "SecureString"
  value = var.twitter_consumer_key

  tags = local.tags
}

resource "aws_ssm_parameter" "twitter_consumer_secret" {
  name  = "${var.param_store_path}/twitter_consumer_secret"
  type  = "SecureString"
  value = var.twitter_consumer_secret

  tags = local.tags
}

resource "aws_ssm_parameter" "twitter_access_token" {
  name  = "${var.param_store_path}/twitter_access_token"
  type  = "SecureString"
  value = var.twitter_access_token

  tags = local.tags
}

resource "aws_ssm_parameter" "twitter_access_secret" {
  name  = "${var.param_store_path}/twitter_access_secret"
  type  = "SecureString"
  value = var.twitter_access_secret

  tags = local.tags
}
