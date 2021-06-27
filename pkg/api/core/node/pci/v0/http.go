package v0

//func httpRequest(ip string, port uint) ([]libvirtxml.NodeDevice, error) {
//	var res pci.Node
//
//	response, err := client.Get("http://"+ip+":"+strconv.Itoa(int(port))+"/api/v1/pci", "")
//	if err != nil {
//		return []libvirtxml.NodeDevice{}, err
//	}
//
//	log.Println(response)
//
//	if err = json.Unmarshal([]byte(response), &res); err != nil {
//		return []libvirtxml.NodeDevice{}, err
//	}
//
//	return res.Data.PCI, nil
//}
//
//func httpRequest1(ip string, port uint) (string, error) {
//	response, err := client.Get("http://"+ip+":"+strconv.Itoa(int(port))+"/api/v1/pci", "")
//	if err != nil {
//		return "", err
//	}
//
//	return response, nil
//}
