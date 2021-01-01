variable "app" {
  default = "boca-chica-bot"
}

variable "debug" {
  description = "Environment variable used to enable/disable debug logs in Lambda."
  default     = false
}

variable "disable_publish" {
  description = "Environment variable used in lambda function for disabling publishing to Twitter/Discord."
  type        = bool
  default     = false
}

variable "cron_schedule_enabled" {
  description = "EventBridge rule enabled or disabled."
  type        = bool
  default     = true
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
