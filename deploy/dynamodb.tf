resource "aws_dynamodb_table" "table" {
  name             = var.tablename
  billing_mode     = "PAY_PER_REQUEST"
  hash_key         = "Date"
  range_key        = "Time"
  stream_enabled   = true
  stream_view_type = "NEW_AND_OLD_IMAGES"

  ttl {
    enabled        = true
    attribute_name = "Time"
  }

  attribute {
    name = "Date"
    type = "S"
  }

  attribute {
    name = "Time"
    type = "N"
  }

  tags = local.tags
}
