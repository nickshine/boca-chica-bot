terraform {
  backend "remote" {
    organization = "nickshine"

    workspaces {
      prefix = "boca-chica-bot-"
    }
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.19"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

locals {
  tags = {
    app = var.app
    env = var.env
  }
}
