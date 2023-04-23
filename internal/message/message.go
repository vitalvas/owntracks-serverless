package message

type MessageLocation struct {
	Acc       int64    `json:"acc,omitempty" dynamodbav:"acc,omitempty"`
	Alt       int64    `json:"alt,omitempty" dynamodbav:"alt,omitempty"`
	Batt      int64    `json:"batt,omitempty" dynamodbav:"batt,omitempty"`
	Bs        int64    `json:"bs,omitempty" dynamodbav:"bs,omitempty"`
	Cog       int64    `json:"cog,omitempty" dynamodbav:"cog,omitempty"`
	Lat       float64  `json:"lat" dynamodbav:"lat"`
	Lon       float64  `json:"lon" dynamodbav:"lon"`
	Rad       int64    `json:"rad,omitempty" dynamodbav:"rad,omitempty"`
	T         string   `json:"t,omitempty" dynamodbav:"t,omitempty"`
	Tid       string   `json:"tid,omitempty" dynamodbav:"tid,omitempty"`
	Tst       int64    `json:"tst" dynamodbav:"tst"`
	Vac       int64    `json:"vac,omitempty" dynamodbav:"vac,omitempty"`
	Vel       int64    `json:"vel,omitempty" dynamodbav:"vel,omitempty"`
	P         float64  `json:"p,omitempty" dynamodbav:"p,omitempty"`
	Poi       string   `json:"poi,omitempty" dynamodbav:"poi,omitempty"`
	Conn      string   `json:"conn,omitempty" dynamodbav:"conn,omitempty"`
	Tag       string   `json:"tag,omitempty" dynamodbav:"tag,omitempty"`
	Inregions []string `json:"inregions,omitempty" dynamodbav:"inregions,omitempty"`
	Inrids    []string `json:"inrids,omitempty" dynamodbav:"inrids,omitempty"`
	SSID      string   `json:"SSID,omitempty" dynamodbav:"SSID,omitempty"`
	BSSID     string   `json:"BSSID,omitempty" dynamodbav:"BSSID,omitempty"`
	CreatedAt int64    `json:"created_at,omitempty" dynamodbav:"created_at,omitempty"`
	M         int64    `json:"m,omitempty" dynamodbav:"m,omitempty"`

	UserID   string `json:"user_id,omitempty" dynamodbav:"user_id,omitempty"`
	DeviceID string `json:"device_id,omitempty" dynamodbav:"device_id,omitempty"`
	RemoteIP string `json:"remote_ip,omitempty" dynamodbav:"remote_ip,omitempty"`
}
