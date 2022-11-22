resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_network_security_group" "example" {
  name                = "example-security-group"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["192.168.0.0/23", "10.90.90.0/24"]
  dns_servers         = ["10.0.0.4", "10.0.0.5"]

  subnet {
    name           = "pre-existing-subnet01"
    address_prefix = "10.90.90.0/25"
  }

  subnet {
    name           = "pre-existing-subnet-02"
    address_prefix = "10.90.90.128/25"
  }

  tags = {
    environment = "Production"
  }
}