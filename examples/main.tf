locals {
  location = "westeurope"
  
}

resource "azurerm_resource_group" "example" {
  name     = "rg-golangexample-weu-dev-01"
  location = local.location
}

// -new-subnets '[{"fill-gap-01":25}, {"fill-gap-02": 26}, {"fill-gap-03": 26}]'
resource "azurerm_virtual_network" "example" {
  name                = "vnet-golangexample-weu-dev-01"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["192.168.0.0/23", "10.90.90.0/24"]
  dns_servers         = ["10.0.0.4", "10.0.0.5"]

  subnet {
    name           = "pre-existing-subnet01"
    address_prefix = "192.168.0.0/25"
  }

  subnet {
    name           = "pre-existing-subnet-02"
    address_prefix = "10.90.90.0/26"
  }

  subnet {
    name           = "pre-existing-subnet-04"
    address_prefix = "10.90.90.128/26"
  }

  tags = var.tags
}

