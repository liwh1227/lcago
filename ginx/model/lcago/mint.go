package lcago

type MintRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type MintResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Data   `json:"data"`
}

type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
