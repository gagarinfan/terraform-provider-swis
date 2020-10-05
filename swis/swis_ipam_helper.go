package swis

import (
	"github.com/mrxinu/gosolar"
	"encoding/json"
	"log"
	"errors"
	"strconv"
	"net"
)

type Subnet struct {
	SubnetId 		int    `json:"subnetid"`
	Uri			    string `json:"uri"`
	CIDR			int    `json:"cidr"`
	GroupTypeText	string `json:"grouptypetext"`
	Address			string `json:"address"`
	VlanName		string `json:"vlan"`
}

type IPEntity struct {
	IpNodeId	int	   `json:"ipnodeid"`
	SubnetId	int	   `json:"subnetid"`
	IPAddress	string `json:"ipaddress"`
	Comments	string `json:"comments"`
	Status		int    `json:"status"`
	Uri			string `json:"uri"`
}

// Check if subnet given by address has DHCP Scope
func checkIfSubnetDHCP(client *gosolar.Client, subnetAddress string) (bool, error) {
	var subnetInfo[] Subnet

	query := "SELECT Vlan,SubnetId,Uri,GroupTypeText,CIDR FROM IPAM.Subnet WHERE  Address='" + subnetAddress + "'AND GroupTypeText='DHCP Scope'"
	res, err := client.Query(query, nil)
	if err != nil {
		log.Fatal(err)
        return false, err
	}
	
	jsonErr := json.Unmarshal(res, &subnetInfo)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return false, jsonErr
	}

	if len(subnetInfo) != 0 && subnetInfo[0].GroupTypeText == "DHCP Scope" {
		log.Print("This subnet has DHCP Scope!")
		return true, nil
	} else {
		log.Print("This subnet does not have dhcp scope")
		return false, nil
	}
}

// Get Subnet ID by it's address
func getSubnetId(client *gosolar.Client, subnetAddress string) (int, error) {
	var subnetInfo[] Subnet

	query := "SELECT Vlan,Address,SubnetId,Uri,CIDR,GroupTypeText FROM IPAM.Subnet WHERE Address='" + subnetAddress + "'"
	res, err := client.Query(query, nil)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	
	jsonErr := json.Unmarshal(res, &subnetInfo)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return 0, jsonErr
	}

	if len(subnetInfo) == 0 {
		subnetErr := errors.New("Could not find provided subnet!")
		return 0, subnetErr
	}

	log.Print("Subnet address is " + subnetInfo[0].Address + " and it's scope: " + subnetInfo[0].GroupTypeText)
	return subnetInfo[0].SubnetId, nil
}

// Get Subnet Address by it's ID
func getSubnetAddress(client *gosolar.Client, subnetId int) (string, error) {
	var subnetInfo[] Subnet
	query := "SELECT Vlan,Address,SubnetId,Uri,CIDR,GroupTypeText FROM IPAM.Subnet WHERE SubnetId='" + strconv.Itoa(subnetId) + "'"
	res, err := client.Query(query, nil)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	
	jsonErr := json.Unmarshal(res, &subnetInfo)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return "", jsonErr
	}

	if len(subnetInfo) == 0 {
		subnetErr := errors.New("Could not find provided subnet!")
		return "", subnetErr
	}

	log.Print("Subnet address is " + subnetInfo[0].Address + " and it's scope: " + subnetInfo[0].GroupTypeText)
	return subnetInfo[0].Address, nil
}

// Get VLAN name by subnet address
func getVlanName(client *gosolar.Client, subnetAddress string) (string, error) {
	var subnetInfo[] Subnet

	query := "SELECT Vlan,Address,SubnetId,Uri,CIDR,GroupTypeText FROM IPAM.Subnet WHERE  Address='" + subnetAddress + "'"
	res, err := client.Query(query, nil)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	
	jsonErr := json.Unmarshal(res, &subnetInfo)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return "", jsonErr
	}

	if len(subnetInfo) == 0 {
		subnetErr := errors.New("Could not find provided subnet!")
		return "", subnetErr
	}

	log.Print("Vlan name is " + subnetInfo[0].VlanName + " and it's scope: " + subnetInfo[0].GroupTypeText)
	return subnetInfo[0].VlanName, nil
}

// Get first free IP Entity in given Subnet by it's ID
func getFreeIpEntity(client *gosolar.Client, subnetId int) (*IPEntity, error) {
	var ipEntity[] IPEntity
	query := "SELECT TOP 1 IpNodeId,IPAddress,Comments,Status,Uri FROM IPAM.IPNode WHERE SubnetId='" + strconv.Itoa(subnetId) + "' and status=2 AND IPOrdinal BETWEEN 11 AND 254"
	res, err := client.Query(query, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	jsonErr := json.Unmarshal(res, &ipEntity)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return nil, jsonErr
	}
	if len(ipEntity) == 0 {
		ipNullErr := errors.New("Provided IP: " + ipEntity[0].IPAddress + " does not exist or empty!")
		return nil, ipNullErr
	}
	if ipEntity[0].Status != 2 {
		ipStatusErr := errors.New("Provided IP should have status 2 (free), but unexpectedly found " + strconv.Itoa(ipEntity[0].Status))
		return nil, ipStatusErr
	}
	return &ipEntity[0], nil
}

// Update IP Entity
func updateIpEntity(client *gosolar.Client, ipEntity IPEntity, status int, comment string) error {
	log.Print("I am going to book IP address: " + ipEntity.IPAddress + " with comment: " + comment + " and status " + strconv.Itoa(status))
	if status == 2 {
		comment = ""
	}
	request := map[string]interface{} {
		"Status"	: status,
		"Comments"	: comment,
	}
	log.Print(request)
	_, err := client.Update(ipEntity.Uri, request)

	if err != nil {
		return err
	} else {
		log.Print(ipEntity.IPAddress + " has been successfully claimed!")
		log.Print(ipEntity)
		return nil
	}
}

// Get IP Entity by it's address
func getIpEntityByAddress(client *gosolar.Client, ipEntityAddress string) (*IPEntity, error) {
	var ipEntity[] IPEntity
	query := "SELECT IpNodeId,SubnetId,IPAddress,Comments,Status,Uri FROM IPAM.IPNode WHERE IPAddress='" + ipEntityAddress + "'"
	res, err := client.Query(query, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	jsonErr := json.Unmarshal(res, &ipEntity)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return nil, jsonErr
	}
	log.Print(ipEntity[0])
	return &ipEntity[0], nil
}

// Valides if address is in proper IPv4 format
func validateAddresses(ip_address string) error {
	log.Print("#### VALIDATING IP ADDRESS ####")
	if net.ParseIP(ip_address) == nil {
		ipv4Error := errors.New("Provided IP " + ip_address + " is not valid IP Address!")
		return ipv4Error
	} else {
		return nil
	}
}

// Validate if given IP Address belongs to given subnet
func validateAddresInSubnet(vlan_address string, mask int, ip_address string) error {
	subnet := vlan_address + "/" + strconv.Itoa(mask)
	ip := ip_address + "/" + strconv.Itoa(mask)
	_,subnet_parsed,_ := net.ParseCIDR(subnet)
	ip_parsed,_,_ := net.ParseCIDR(ip)

	if subnet_parsed.Contains(ip_parsed) {
		return nil
	} else {
		ipv4SubnetError := errors.New("Provided IP " + ip_address + " does not belong to " + vlan_address + "/" + strconv.Itoa(mask) + " subnet!")
		return ipv4SubnetError
	}
}