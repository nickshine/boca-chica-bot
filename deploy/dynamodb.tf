resource "aws_dynamodb_table" "table" {
  name             = var.tablename
  billing_mode     = "PAY_PER_REQUEST"
  hash_key         = "Date"
  range_key        = "RawTimeRange"
  stream_enabled   = true
  stream_view_type = "NEW_AND_OLD_IMAGES"

  ttl {
    enabled        = true
    attribute_name = "Expires"
  }

  attribute {
    name = "Date"
    type = "S"
  }

  attribute {
    name = "RawTimeRange"
    type = "S"
  }

  tags = local.tags
}
