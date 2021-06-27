package v0

//func Get(c *gin.Context) {
//	id, _ := strconv.Atoi(c.Param("id"))
//
//	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
//	if resultAdmin.Err != nil {
//		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
//		return
//	}
//
//	resultNode := dbNode.Get(node.ID, &core.Node{Model: gorm.Model{ID: uint(id)}})
//	if resultNode.Err != nil {
//		log.Println(resultNode.Err)
//		c.JSON(http.StatusInternalServerError, common.Error{Error: resultNode.Err.Error()})
//		return
//	}
//	response, err := httpRequest(resultNode.Node[0].IP, resultNode.Node[0].Port)
//	if err != nil {
//		log.Println(err)
//		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, pci.Result{PCI: response})
//}
