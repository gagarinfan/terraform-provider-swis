package ipam

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrxinu/gosolar"
	"log"
)

func resourceIp() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpCreate,
		Read:   resourceIpRead,
		Update: resourceIpUpdate,
		Delete: resourceIpDelete,
		Schema: map[string]*schema.Schema{
			"vlan_address": {
				Type:        schema.TypeString,
				Description: "Vlan address",
				Required:    true,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Description: "Ip address if static",
				Optional:    true,
			},
			"status_code": {
				Type:        schema.TypeInt,
				Description: "2- free, 1- assigned",
				Optional:    true,
			},
		},
	}
}

func resourceIpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gosolar.Client)
	vlan_addres := d.Get("vlan_address").(string)
	subnetId,_ := getSubnetId(client, vlan_addres)
	log.Print("Printing subnetId:")
	log.Print(subnetId)
	subnetDHCP,_ := checkIfSubnetDHCP(client, vlan_addres)
	log.Print("Printing subnetDHCP:")
	log.Print(subnetDHCP)
	ipEntity,_ := getFreeIpEntity(client, subnetId)
	log.Print("Printing ipEntity:")
	log.Print(ipEntity)

	updateIpEntity(client, *ipEntity, 2, "")

	d.Set("ip_address", ipEntity.IPAddress)
	d.SetId(ipEntity.IPAddress)
	return nil
}

func resourceIpRead(d *schema.ResourceData, meta interface{}) error {
	log.Print("koko")
	return nil
}

func resourceIpUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gosolar.Client)
	ip_address := d.Get("ip_address").(string)
	status_code := d.Get("status_code").(int)
	ipEntity,_ := getIpEntityByAddress(client, ip_address)
	updateIpEntity(client, *ipEntity, status_code, "")

	return nil
}

func resourceIpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gosolar.Client)
	ip_address := d.Get("ip_address").(string)
	ipEntity,_ := getIpEntityByAddress(client, ip_address)

	updateIpEntity(client, *ipEntity, 2, "")

	return nil
}