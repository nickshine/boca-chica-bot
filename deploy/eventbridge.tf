resource "aws_cloudwatch_event_rule" "cron" {
  name                = "boca-chica-bot-cron-${var.env}"
  description         = "Execute the boca-chica-bot lambda every 10 mins from 8amto 6pm CST Mon-Fri."
  schedule_expression = "cron(0/10 14-23 ? * MON-FRI *)"
  is_enabled          = true

  tags = local.tags
}

resource "aws_cloudwatch_event_target" "lambda" {
  rule = aws_cloudwatch_event_rule.cron.name
  arn  = "arn:aws:lambda:us-east-1:401054096140:function:boca-chica-bot-test:test"
}
