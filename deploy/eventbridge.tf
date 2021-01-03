resource "aws_cloudwatch_event_rule" "cron" {
  name                = "boca-chica-bot-cron-${var.env}"
  description         = "Execute the boca-chica-bot lambda every 10 mins from 8am to 8pm CST every day."
  schedule_expression = "cron(0/6 14-23,0-2 * * ? *)"
  is_enabled          = var.cron_schedule_enabled

  tags = local.tags
}

resource "aws_cloudwatch_event_target" "lambda" {
  rule = aws_cloudwatch_event_rule.cron.name
  arn  = aws_lambda_alias.scraper_lambda.arn
}
