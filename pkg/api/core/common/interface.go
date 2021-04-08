package common

type Error struct {
	Error string `json:"error"`
}

type Result struct {
	UUID   string `json:"uuid"`
	Result string `json:"result"`
}
