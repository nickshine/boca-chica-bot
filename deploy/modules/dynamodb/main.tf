resource "aws_dynamodb_table" "table" {
  name         = var.tablename
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "Date"
  range_key    = "Time"

  ttl {
    enabled        = true
    attribute_name = "Expires"
  }

  attribute {
    name = "Date"
    type = "S"
  }

  attribute {
    name = "Time"
    type = "S"
  }

  tags = var.tags
}
