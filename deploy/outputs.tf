output "table_arn" {
  value = module.dynamodb.arn
}

output "table_id" {
  value = module.dynamodb.id
}

output "function_arn" {
  value = aws_lambda_function.lambda.qualified_arn
}

output "event_rule_arn" {
  value = aws_cloudwatch_event_rule.cron.arn
}
