module "dynamodb" {
  source    = "./modules/dynamodb"
  tablename = var.tablename
  tags      = local.tags
}

locals {
  tags = {
    app = var.app
    env = var.env
  }
}
