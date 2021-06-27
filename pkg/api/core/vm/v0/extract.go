package v0

//func extractImaCon(vmTemplate template.Template) (config.ImaCon, vm.GetImaCon, error) {
//	var getImaCon vm.GetImaCon
//
//	for _, imaCon := range config.Conf.ImaCon {
//		response, err := client.Get("http://"+imaCon.IP+":"+strconv.Itoa(int(imaCon.Port))+"/api/v1/storage/uuid/"+
//			vmTemplate.UUID, "")
//		if err != nil {
//			log.Println(err)
//		}
//
//		err = json.Unmarshal([]byte(response), &getImaCon)
//		if err != nil {
//			log.Println(err)
//		}
//
//		// Status OKである場合、結果とnilを返す
//		if getImaCon.Status == http.StatusOK {
//			return imaCon, getImaCon, nil
//		}
//	}
//
//	// Imageが見つからなかった場合、errorを返す
//	return config.ImaCon{}, vm.GetImaCon{}, fmt.Errorf("Error: image not found... ")
//}
