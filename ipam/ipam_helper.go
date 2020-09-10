package ipam

import (
	"github.com/mrxinu/gosolar"
	"encoding/json"
	"log"
	"errors"
	"strconv"
)

type Subnet struct {
	SubnetId 		int    `json:"subnetid"`
	Uri			    string `json:"uri"`
	CIDR			int    `json:"cidr"`
	GroupTypeText	string `json:"grouptypetext"`
	Address			string `json:"address"`
}

type IPEntity struct {
	IpNodeId	int	   `json:"ipnodeid"`
	SubnetId	int	   `json:"subnetid"`
	IPAddress	string `json:"ipaddress"`
	Comments	string `json:"comments"`
	Status		int    `json:"status"`
	Uri			string `json:"uri"`
}

func checkIfSubnetDHCP(client *gosolar.Client, subnetAddress string) (bool,error) {
	var subnetInfo[] Subnet

	query := "SELECT SubnetId,Uri,GroupTypeText,CIDR FROM IPAM.Subnet WHERE  Address='" + subnetAddress + "'AND GroupTypeText='DHCP Scope'"
	log.Print("Querry for checkIfSubnetDHCP is " + query)
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

func getSubnetId(client *gosolar.Client, subnetAddress string) (int, error) {
	var subnetInfo[] Subnet

	query := "SELECT Address,SubnetId,Uri,CIDR,GroupTypeText FROM IPAM.Subnet WHERE  Address='" + subnetAddress + "'"
	log.Print("Querry for getSubnetId is " + query)
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
		log.Fatal(subnetErr)
		return 0, subnetErr
	}

	log.Print("Subnet address is " + subnetInfo[0].Address + "and it's scope: " + subnetInfo[0].GroupTypeText)
	return subnetInfo[0].SubnetId, nil
}

func getFreeIpEntity(client *gosolar.Client, subnetId int) (*IPEntity, error) {
	var ipEntity[] IPEntity
	query := "SELECT TOP 1 IpNodeId,IPAddress,Comments,Status,Uri FROM IPAM.IPNode WHERE SubnetId='" + strconv.Itoa(subnetId) + "' and status=2 AND IPOrdinal BETWEEN 11 AND 254"
	log.Print("Querry for getFreeIps is " + query)
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
		ipNullErr := errors.New("Provided IP: " + ipEntity[0].IPAddress + " does not exist or empty")
		return nil, ipNullErr
	}
	if ipEntity[0].Status != 2 {
		ipStatusErr := errors.New("Provided IP should have status 2 (free), but unexpectedly found " + strconv.Itoa(ipEntity[0].Status))
		return nil, ipStatusErr
	}
	return &ipEntity[0], nil
}

func updateIpEntity(client *gosolar.Client, ipEntity IPEntity, status int, comment string) error {
	log.Print("I am going to book IP address: " + ipEntity.IPAddress)
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
		return nil
	}
}

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