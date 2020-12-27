output "table_arn" {
  value = aws_dynamodb_table.table.arn
}

output "table_id" {
  value = aws_dynamodb_table.table.id
}

output "scraper_function_arn" {
  value = aws_lambda_function.scraper_lambda.qualified_arn
}

output "event_rule_arn" {
  value = aws_cloudwatch_event_rule.cron.arn
}

output "publisher_function_arn" {
  value = aws_lambda_function.publisher_lambda.qualified_arn
}
