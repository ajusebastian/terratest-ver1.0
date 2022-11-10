terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "3.8.0"
    }
  }
}
provider "azurerm" {
  subscription_id = "a3d772dc-93d3-4b00-a52e-e8c86dfb1feb"
  tenant_id       = "5ab73113-e42a-44d9-b750-ab52b3bbb0b7"
  client_id       = "0cc2f215-7288-4f93-91e9-db21835113d2"
  client_secret   = "6D68Q~xmpP4~17txCkdWHJqzIi~gwsCdnnw21ccd"
  features {

  }
}
terraform {
  cloud {
    organization = "rohithpadmanabhan"

    workspaces {
      name = "testworkspace"
    }
  }
}