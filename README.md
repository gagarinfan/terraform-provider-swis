# Terraform provider (plugin) for Solar Winds

https://www.solarwinds.com/

## Authors
- Michał Gębka <michal.j.gebka@gmail.com>

## Issues and Contributing
If you find an issue with this provider, please report it. Contributions are welcome.


## Available resources:

[swis_ipam](https://www.solarwinds.com/ip-address-manager):
- allows to create, delete, manage IP addresses
- validates if subnet has DHCP Scope

## Argument Reference

swis_ipam:
- **vlan_address** (string): (Required) vlan address to find/search IP Address
- **vlan_name** (string): (Optional) the vlan name of IP Address
- **comment** (string): (Required) comment to assign to IP Address. Must start with `Reserved by`
- **ip_address** (string): (Optional) if IP Address must be static then write it here. Default nil
- **status_code** (int): (Optional) status code for IP Address instance: 1- assigned, 2- free. Default 1
- **avoid_dhcp_scope** (bool): (Optional) set to true if want to avoid subnets with DHCP Scope. Default true

## Attributes Reference

swis_ipam:
- **vlan_address** : the vlan address of IP Address
- **vlan_name** : the vlan name of IP Address
- **comment** : the comment of IP Address
- **ip_address** : the IP Address
- **status_code** : the value of IP Address's status code: 1- assigned, 2- free


## Example terraform file:

```tf
variable "swis_host" {
    default = <swis host>
}

variable "swis_username" {
    default = <swis username>
}

variable "swis_password" {
    default = <swis password>
}

provider "swis" {
    host = var.swis_host
    username = var.swis_username
    password = var.swis_password
}

//Book 10.12.72.22 in 10.12.72.0 vlan if it has not DHCP Scope and comment is as reserved for s1slt000123
resource "swis_ipam" "static" {
    vlan_address = "10.12.72.0"
    vlan_name = "1234" //optional
    comment = "Reserved by s1slt000123"
    ip_address = "10.12.72.22" //optional
    status_code = 1 //optional
    avoid_dhcp_scope = true //optional
}

//Book random IP Address in 10.12.72.0 vlan
resource "swis_ipam" "auto" {
    vlan_address = "10.12.72.0"
    comment = "Reserved by s1slt000321"
}

//Get reserved static IP Address
output "ip_address_static" {
    value = swis_ipam.static.ip_address
}

//Get reserved auto IP Address
output "ip_address_auto" {
    value = swis_ipam.auto.ip_address
}

```

## Build

```
make build
```