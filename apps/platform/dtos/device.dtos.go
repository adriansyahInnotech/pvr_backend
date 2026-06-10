package dtos

type Device []struct {
	Params struct {
		Sn string `json:"sn"`
	} `json:"params"`
	Data struct {
		AreaID   string `json:"area_id"`
		Timezone int    `json:"timezone"`
	} `json:"data"`
}

type DeleteDevice struct {
	SN string `json:"sn"`
}
