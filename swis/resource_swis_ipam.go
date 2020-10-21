package swis

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mrxinu/gosolar"
	"errors"
	"log"
	"strings"
	"time"
	"math/rand"
	"strconv"
)

func resourceIp() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpCreate,
		Read:   resourceIpRead,
		Update: resourceIpUpdate,
		Delete: resourceIpDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIpImport,
		},
		Schema: map[string]*schema.Schema{
			"vlan_address": {
				Type		: 	schema.TypeString,
				Description	: 	"Vlan address",
				Required	:	true,
			},
			"vlan_name": {
				Type		:	schema.TypeString,
				Description	: 	"Vlan name",
				Optional	:	true,
				Computed	:	true,
			},
			"vlan_mask": {
				Type		:	schema.TypeInt,
				Description	: 	"Vlan mask",
				Optional	:	true,
				Default		:	24,
			},
			"comment": {
				Type		:	schema.TypeString,
				Description	:	"Server name",
				Required	:	true,
			},
			"status_code": {
				Type		:	schema.TypeInt,
				Description	:	"2- free, 1- assigned",
				Optional	:	true,
				Default		:	1,
			},
			"ip_address": {
				Type		:	schema.TypeString,
				Description	:	"Ip address if static",
				Optional	:	true,
				Computed	:	true,
			},
			"avoid_dhcp_scope": {
				Type		:	schema.TypeBool,
				Description	:	"If true will not set ip address in vlan with DHCP scope",
				Optional	:	true,
				Default		:	true,
			},
			
		},
	}
}

func resourceIpCreate(d *schema.ResourceData, meta interface{}) error {
	//Sleep for random 1-3s to avoid locking same IP Entity by two or more resources
	log.Print("############# CREATING #############")
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(3)
	time.Sleep(time.Duration(n)*time.Second)

	client := meta.(*gosolar.Client)
	//declare vars
	vlan_address := d.Get("vlan_address").(string)
	comment := d.Get("comment").(string)
	avoid_dhcp_scope := d.Get("avoid_dhcp_scope").(bool)
	ip_address := d.Get("ip_address").(string)
	status_code := d.Get("status_code").(int)
	vlan_name := d.Get("vlan_name").(string)
	vlan_mask := d.Get("vlan_mask").(int)
	
	comuptedVlanName, VlanNameErr := getVlanName(client, vlan_address)
	if VlanNameErr != nil {
		return VlanNameErr
	}
	if vlan_name != "" && vlan_name != comuptedVlanName {
		vlanNameMismatch := errors.New("There is mismatch in vlan name that you've provided (" + vlan_name + ") and comupted value which is " + comuptedVlanName)
		return vlanNameMismatch
	} else if vlan_name == "" {
		d.Set("vlan_name", comuptedVlanName)
	}
	
	subnetId, getSubnetErr := getSubnetId(client, vlan_address)
	if getSubnetErr != nil {
		return getSubnetErr
	}
	if avoid_dhcp_scope {
		log.Print("#### Detected avoid_dhcp_scope flag ####")
		subnetDHCP,getSubnetDhcpErr := checkIfSubnetDHCP(client, vlan_address)
		if getSubnetDhcpErr != nil {
			return getSubnetDhcpErr
		}
		log.Print(subnetDHCP)
		if subnetDHCP && ip_address == ""{
			d.Set("ip_address", "dhcp")
			d.SetId("dhcp")
			return nil
		} else if subnetDHCP && ip_address != "" {
			subnetMismathErr := errors.New("You are trying to get static IP from subnet with DHCP Scope!")
			return subnetMismathErr
		}
	} 
	if ip_address == "" {
		log.Print("#### Detected auto_ip ####")
		ipEntity,getIpError := getFreeIpEntity(client, subnetId)
		if getIpError != nil {
			return getIpError
		}
		updateErr := updateIpEntity(client, *ipEntity, status_code, comment)
		if updateErr != nil {
			return updateErr
		}
		d.Set("ip_address", ipEntity.IPAddress)
		d.SetId(ipEntity.IPAddress)
		return nil
	} else {
		log.Print("#### Detected static IP ####")
		ipError := validateAddresses(ip_address)
		if ipError != nil {
			return ipError
		}
		ipSubnetError := validateAddresInSubnet(vlan_address, vlan_mask, ip_address)
		if ipSubnetError != nil {
			return ipSubnetError
		}
		ipEntity,getIpError := getIpEntityByAddress(client, ip_address)
		if getIpError != nil {
			return getIpError
		}
		if ipEntity.Status != 2 {
			statusErr := errors.New("You are trying to get IP that is already assigned!")
			return statusErr
		}
		updateErr := updateIpEntity(client, *ipEntity, status_code, comment)
		if updateErr != nil {
			return updateErr
		}

		d.Set("ip_address", ipEntity.IPAddress)
		d.SetId(ipEntity.IPAddress)

		return nil
	}
}

func resourceIpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gosolar.Client)
	
	//declare vars
	id := d.Id()
	vlan_address := d.Get("vlan_address").(string)
	comment := d.Get("comment").(string)
	avoid_dhcp_scope := d.Get("avoid_dhcp_scope").(bool)
	ip_address := d.Get("ip_address").(string)
	status_code := d.Get("status_code").(int)
	vlan_name := d.Get("vlan_name").(string)

	//Validate if it's dhcp error to handle
	if id == "dhcp" && ip_address == "dhcp" {
		log.Print("This is DHCP scope to handle")
		return nil
	}
	log.Print("########## PERFORMING READ ##########")
	log.Print("VLAN name: " + vlan_name)
	log.Print("ID: " + id)
	log.Print("VLAN address: " + vlan_address)
	log.Print("Comment: " + comment)
	log.Print(avoid_dhcp_scope)
	log.Print("IP address: " + ip_address)
	log.Print("Status code: " + strconv.Itoa(status_code))

	comuptedVlanName, VlanNameErr := getVlanName(client, vlan_address)
	if VlanNameErr != nil {
		return VlanNameErr
	}
	if vlan_name != "" && vlan_name != comuptedVlanName {
		vlanNameMismatch := errors.New("There is mismatch in vlan name that you've provided (" + vlan_name + ") and comupted value which is " + comuptedVlanName)
		return vlanNameMismatch
	}

	ipEntity,_ := getIpEntityByAddress(client, ip_address)

	//Validate if provided ip address is assigned to this machine
	if !strings.Contains(ipEntity.Comments, comment) && ipEntity.Status != 2 {
		assignError := errors.New("IP address " + ip_address + " is not assigned to " + comment)
		return assignError
	}

	//Validate if provided subnet is DHCP
	dhcpScope,dhcpErr := checkIfSubnetDHCP(client, vlan_address)
	if dhcpErr != nil {
		log.Fatal(dhcpErr)
		return  dhcpErr
	}
	//Validate if subnet is DHCP AND avoid_dhcp_scope is true
	if avoid_dhcp_scope && dhcpScope {
		dhcpError := errors.New("avoid_dhcp_flag set to true, but subnet HAS dhcp scope")
		return dhcpError
	}
	
	d.Set("ip_address", ipEntity.IPAddress)

	return nil
}

func resourceIpUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gosolar.Client)
	log.Print("############# UPDATING #############")

	//declare vars
	comment := d.Get("comment").(string)
	vlan_address := d.Get("vlan_address").(string)
	vlan_name := d.Get("vlan_name").(string)

	comuptedVlanName, VlanNameErr := getVlanName(client, vlan_address)
	if VlanNameErr != nil {
		return VlanNameErr
	}
	if vlan_name != "" && vlan_name != comuptedVlanName {
		d.Set("vlan_name", comuptedVlanName)
		vlanNameMismatch := errors.New("There is mismatch in vlan name that you've provided (" + vlan_name + ") and comupted value which is " + comuptedVlanName)
		return vlanNameMismatch
	}
	ip_address := d.Id()
	if ip_address == "dhcp" {
		return nil
	}
	status_code := d.Get("status_code").(int)

	log.Print("########### STATE ##########")
	log.Print(d.State())
	
	ipEntity,_ := getIpEntityByAddress(client, ip_address)
	updateIpEntity(client, *ipEntity, status_code, comment)
	
	d.Set("ip_address", ip_address)

	return nil
}

func resourceIpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gosolar.Client)
	if d.Id() == "dhcp" {
		return nil
	}
	ip_address := d.Get("ip_address").(string)
	ipEntity,_ := getIpEntityByAddress(client, ip_address)

	updateIpEntity(client, *ipEntity, 2, "")
	return nil
}

func resourceIpImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*gosolar.Client)
	log.Print("########## PERFORMING IMPORT ##########")
	ip_address := d.Id()
	ipEntity, ipEntityErr :=getIpEntityByAddress(client, ip_address)
	if ipEntityErr != nil {
		return nil, ipEntityErr
	}
	vlanAddress, vlanAddressErr := getSubnetAddress(client, ipEntity.SubnetId)
	if vlanAddressErr != nil {
		return nil, vlanAddressErr
	}
	comuptedVlanName, VlanNameErr := getVlanName(client, vlanAddress)
	if VlanNameErr != nil {
		return nil, VlanNameErr
	}
	d.Set("comment", ipEntity.Comments)
	d.Set("status_code", ipEntity.Status)
	d.Set("vlan_address", vlanAddress)
	d.Set("ip_address", ip_address)
	d.Set("vlan_name", comuptedVlanName)

	return []*schema.ResourceData{d}, nil
}