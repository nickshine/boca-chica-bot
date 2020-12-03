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

variable "tablename" {
  description = "DynamoDB table name"
}
