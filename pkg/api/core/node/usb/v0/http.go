package v0

//func httpRequest(ip string, port uint) ([]libvirtxml.NodeDevice, error) {
//	var res usb.Node
//
//	response, err := client.Get("http://"+ip+":"+strconv.Itoa(int(port))+"/api/v1/usb", "")
//	if err != nil {
//		return []libvirtxml.NodeDevice{}, err
//	}
//
//	if err = json.Unmarshal([]byte(response), &res); err != nil {
//		return []libvirtxml.NodeDevice{}, err
//	}
//
//	return res.Data.USB, nil
//}
