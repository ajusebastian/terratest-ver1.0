locals {
  resource_group_name = "aju-resourse-grp"
  location            = "West Europe"
  virtual_network = {
    name          = "app-network"
    address_space = "10.0.0.0/16"
  }
  subnets = [
    {
      name           = "subnetA"
      address_prefix = "10.0.0.0/24"
    },
    {
      name           = "subnetB"
      address_prefix = "10.0.1.0/24"
    }
  ]
  networksecuritygroup_rules = [
    {
      priority               = 200
      destination_port_range = "22"
    },
    {
      priority               = 300
      destination_port_range = "80"
    }
  ]

}
resource "azurerm_resource_group" "aju-resourse-grp" {
  name     = local.resource_group_name
  location = local.location

}
/*resource "azurerm_virtual_network" "appnetwork" {
  name                = local.virtual_network.name
  location            = local.location
  resource_group_name = local.resource_group_name
  address_space       = [local.virtual_network.address_space]
  depends_on = [
    azurerm_resource_group.aju-resourse-grp
  ]
}
resource "azurerm_subnet" "subnetA" {
  name                 = local.subnets[0].name
  resource_group_name  = local.resource_group_name
  virtual_network_name = azurerm_virtual_network.appnetwork.name
  address_prefixes     = [local.subnets[0].address_prefix]
  depends_on = [
    azurerm_virtual_network.appnetwork
  ]
}

resource "azurerm_subnet" "subnetB" {
  name                 = local.subnets[1].name
  resource_group_name  = local.resource_group_name
  virtual_network_name = azurerm_virtual_network.appnetwork.name
  address_prefixes     = [local.subnets[1].address_prefix]
  depends_on = [
    azurerm_virtual_network.appnetwork
  ]
}

resource "azurerm_network_interface" "appinterface" {
  name                = "appinterface"
  location            = local.location
  resource_group_name = local.resource_group_name
  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.subnetA.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.appip.id
  }
  depends_on = [
    azurerm_subnet.subnetA
  ]
}

resource "azurerm_public_ip" "appip" {
  name                = "app-ip"
  resource_group_name = local.resource_group_name
  location            = local.location
  allocation_method   = "Static"
  depends_on = [
    azurerm_resource_group.aju-resourse-grp
  ]
}

resource "azurerm_network_security_group" "appnsg" {
  name                = "app-nsg"
  location            = local.location
  resource_group_name = local.resource_group_name
  dynamic "security_rule" {
    for_each = local.networksecuritygroup_rules
    content {
      name                       = "Allow-${security_rule.value.destination_port_range}"
      priority                   = security_rule.value.priority
      direction                  = "Inbound"
      access                     = "Allow"
      protocol                   = "Tcp"
      source_port_range          = "*"
      destination_port_range     = security_rule.value.destination_port_range
      source_address_prefix      = "*"
      destination_address_prefix = "*"
    }
  }

  depends_on = [
    azurerm_resource_group.aju-resourse-grp
  ]
}

resource "azurerm_subnet_network_security_group_association" "appnsglink" {
  subnet_id                 = azurerm_subnet.subnetA.id
  network_security_group_id = azurerm_network_security_group.appnsg.id
}



resource "tls_private_key" "linuxkey" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "local_file" "linuxpemkey" {
  filename = "linuxkey.pem"
  content  = tls_private_key.linuxkey.private_key_pem
  depends_on = [
    tls_private_key.linuxkey
  ]
}

resource "azurerm_linux_virtual_machine" "linuxvm" {
  name                = "linuxvm"
  resource_group_name = local.resource_group_name
  location            = local.location
  size                = "Standard_D2s_v3"
  admin_username      = "ajupulickal"
  network_interface_ids = [
    azurerm_network_interface.appinterface.id
  ]

  admin_ssh_key {
    username   = "ajupulickal"
    public_key = tls_private_key.linuxkey.public_key_openssh
  }
  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-focal"
    sku       = "20_04-lts"
    version   = "latest"
  }

  depends_on = [
    azurerm_network_interface.appinterface,
    azurerm_resource_group.aju-resourse-grp,
    tls_private_key.linuxkey
  ]

}*/

resource "azurerm_storage_account" "aju-storageaccount" {
  #count                    = 2
  name                      = var.storageaccountname
  resource_group_name       = local.resource_group_name
  location                  = local.location
  account_tier              = "Standard"
  account_replication_type  = "LRS"
  account_kind              = "StorageV2"
  enable_https_traffic_only = var.enable_https_traffic_only
  depends_on = [
    azurerm_resource_group.aju-resourse-grp
  ]
  tags = {
    department = "TEST",
    environemt = "prod"
  }

}

