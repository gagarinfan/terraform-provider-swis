provider "ipam" {
    host = 
    username = 
	password = 
}

resource "ipam" "default" {
    vlan_address = "192.168.0.0"
    status_code = 2
    //ip_address = "192.168.0.0"
}

output "changed_ip" {
    value = ipam.default.ip_address
}

output "status" {
    value = ipam.default.status_code
}