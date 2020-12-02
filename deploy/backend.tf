terraform {
  backend "remote" {
    organization = "nickshine"

    workspaces {
      prefix = "boca-chica-bot-"
    }
  }
}
