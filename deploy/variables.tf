variable "app" {
  default = "boca-chica-bot"
}

variable "debug" {
  description = "Environment variable used to enable/disable debug logs in Lambda."
  default     = false
}

variable "disable_tweets" {
  description = "Environment variable used in lambda function for disabling Tweets."
  type        = bool
  default     = false
}

variable "env" {
  description = "environment"
}

variable "param_store_path" {
  description = "Prefix for SSM Parameter Store SecureString secret names."
}

variable "tablename" {
  description = "DynamoDB table name"
}
