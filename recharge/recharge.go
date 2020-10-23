package recharge

type CPayInfo struct {
	UID   string `json:"uid"`
	SvrID int32  `json:"svrid"`
	PayID int32  `json:"payid"`
	Sign  int32  `json:"sign"`
}

//广告价格
type AdPrice struct {
	Price  float64 `json:"price"` //元
	AdID   string  `json:"adId"`
	PlatID int32   `json:"platid"`
}

//广告价格和默认值
type AllAdPrice struct {
	AdPs   []AdPrice `json:"yesAdPrice"`
	CaDeng int32     `json:"touchLightGold"`
	Video  int32     `json:"videoGold"`
	TimeHB int32     `json:"taskGold"`
}
