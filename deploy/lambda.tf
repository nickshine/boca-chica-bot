data "aws_iam_policy_document" "lambda-exec" {
  statement {
    actions = ["sts:AssumeRole"]
    effect  = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy" "AmazonDynamoDBFullAccess" {
  arn = "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"
}

data "aws_iam_policy" "AWSLambdaBasicExecutionRole" {
  arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

data "aws_iam_policy" "AmazonSSMReadOnlyAccess" {
  arn = "arn:aws:iam::aws:policy/AmazonSSMReadOnlyAccess"
}

resource "aws_iam_role" "lambda-exec" {
  name               = "${var.app}-exec-role-${var.env}"
  assume_role_policy = data.aws_iam_policy_document.lambda-exec.json

  tags = local.tags
}

resource "aws_iam_role_policy_attachment" "AmazonDynamoDBFullAccess" {
  role       = aws_iam_role.lambda-exec.name
  policy_arn = data.aws_iam_policy.AmazonDynamoDBFullAccess.arn
}

resource "aws_iam_role_policy_attachment" "AWSLambdaBasicExecutionRole" {
  role       = aws_iam_role.lambda-exec.name
  policy_arn = data.aws_iam_policy.AWSLambdaBasicExecutionRole.arn
}

resource "aws_iam_role_policy_attachment" "AmazonSSMReadOnlyAccess" {
  role       = aws_iam_role.lambda-exec.name
  policy_arn = data.aws_iam_policy.AmazonSSMReadOnlyAccess.arn
}

############################ scraper lambda ####################################

resource "aws_lambda_function" "scraper_lambda" {
  filename      = "${path.root}/lambda/scraper-lambda.zip"
  function_name = "${var.app}-scraper-${var.env}"
  role          = aws_iam_role.lambda-exec.arn
  handler       = "${var.app}-scraper"
  publish       = true

  source_code_hash = filebase64sha256("${path.root}/lambda/scraper-lambda.zip")
  runtime          = "go1.x"

  environment {
    variables = {
      DEBUG           = var.debug
      AWS_ENVIRONMENT = var.env
    }
  }

  tags = local.tags
}

resource "aws_lambda_alias" "scraper_lambda" {
  name             = var.env
  description      = "environment"
  function_name    = aws_lambda_function.scraper_lambda.arn
  function_version = aws_lambda_function.scraper_lambda.version
}

resource "aws_lambda_permission" "allow_eventbridge" {
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.scraper_lambda.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.cron.arn
  qualifier     = aws_lambda_alias.scraper_lambda.name
}

############################ publisher lambda ####################################

resource "aws_lambda_function" "publisher_lambda" {
  filename      = "${path.root}/lambda/publisher-lambda.zip"
  function_name = "${var.app}-publisher-${var.env}"
  role          = aws_iam_role.lambda-exec.arn
  handler       = "${var.app}-publisher"
  publish       = true
  timeout       = 60

  source_code_hash = filebase64sha256("${path.root}/lambda/publisher-lambda.zip")
  runtime          = "go1.x"

  environment {
    variables = {
      DEBUG           = var.debug
      DISABLE_PUBLISH = var.disable_publish
      AWS_ENVIRONMENT = var.env
    }
  }

  tags = local.tags
}

resource "aws_lambda_alias" "publisher_lambda" {
  name             = var.env
  description      = "environment"
  function_name    = aws_lambda_function.publisher_lambda.arn
  function_version = aws_lambda_function.publisher_lambda.version
}


resource "aws_lambda_event_source_mapping" "dynamodb" {
  event_source_arn       = aws_dynamodb_table.table.stream_arn
  function_name          = aws_lambda_function.publisher_lambda.arn
  batch_size             = 1
  starting_position      = "TRIM_HORIZON"
  maximum_retry_attempts = 1
}
