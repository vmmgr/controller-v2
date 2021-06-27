package v2

func getCPUMode(mode uint) string {
	//デフォルトではcustomModeになる
	//https://access.redhat.com/documentation/ja-jp/red_hat_enterprise_linux/6/html/virtualization_administration_guide/sect-libvirt-dom-xml-cpu-model-top

	if mode == 1 {
		return "host-model"
	} else if mode == 2 {
		return "host-passthrough"
	}
	return "custom"
}

func getArchConvert(arch uint) string {
	if arch == 32 {
		return "i686"
	}
	return "x86_64"
}
