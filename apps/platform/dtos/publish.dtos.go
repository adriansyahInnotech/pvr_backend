package dtos

type Publish struct {
	DeviceID string `json:"device_id"`
	Message  string `json:"message"`
}

type PersonPublish struct {
	SeqID    string     `json:"seqId"`
	CallType string     `json:"callType"`
	Info     InfoPerson `json:"info"`
}

type InfoPerson struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	IdcardNum  string `json:"idcardNum"`
	Blacklist  int    `json:"blacklist"`
	Remark     string `json:"remark"`
	ExpireTime string `json:"expireTime"`
}
