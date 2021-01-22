package v0

import "testing"

func Test(t *testing.T) {
	response, _ := httpRequest("127.0.0.1", 8080)
	t.Log(response)
}
