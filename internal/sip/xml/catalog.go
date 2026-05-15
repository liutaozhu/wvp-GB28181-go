package xml

import "encoding/xml"

// CatalogQuery is the XML for catalog query request (GB28181)
type CatalogQuery struct {
	XMLName  xml.Name `xml:"Query"`
	CmdType  string   `xml:"CmdType"`
	SN       int      `xml:"SN"`
	DeviceID string   `xml:"DeviceID"`
}

// CatalogResponse is the XML for catalog response from device
type CatalogResponse struct {
	XMLName   xml.Name        `xml:"Response"`
	CmdType   string          `xml:"CmdType"`
	SN        int             `xml:"SN"`
	DeviceID  string          `xml:"DeviceID"`
	SumNum    int             `xml:"SumNum"`
	ItemList  []CatalogItem   `xml:"Item"`
}

// CatalogItem represents a single channel in catalog
type CatalogItem struct {
	DeviceID             string  `xml:"DeviceID"`
	Name                 string  `xml:"Name"`
	Manufacturer         string  `xml:"Manufacturer"`
	Model                string  `xml:"Model"`
	Owner                string  `xml:"Owner"`
	CivilCode            string  `xml:"CivilCode"`
	Block                string  `xml:"Block"`
	Address              string  `xml:"Address"`
	Parental             int     `xml:"Parental"`
	ParentID             string  `xml:"ParentID"`
	SafetyWay            int     `xml:"SafetyWay"`
	RegisterWay          int     `xml:"RegisterWay"`
	CertNum              string  `xml:"CertNum"`
	Certifiable          int     `xml:"Certifiable"`
	ErrCode              int     `xml:"ErrCode"`
	EndTime              string  `xml:"EndTime"`
	Secrecy              int     `xml:"Secrecy"`
	IPAddress            string  `xml:"IPAddress"`
	Port                 int     `xml:"Port"`
	Password             string  `xml:"Password"`
	Status               string  `xml:"Status"`
	Longitude            float64 `xml:"Longitude"`
	Latitude             float64 `xml:"Latitude"`
	PTZType              int     `xml:"PTZType"`
	PositionType         int     `xml:"PositionType"`
	RoomType             int     `xml:"RoomType"`
	UseType              int     `xml:"UseType"`
	SupplyLightType      int     `xml:"SupplyLightType"`
	DirectionType        int     `xml:"DirectionType"`
	Resolution           string  `xml:"Resolution"`
	BusinessGroupID      string  `xml:"BusinessGroupID"`
	DownloadSpeed        string  `xml:"DownloadSpeed"`
	SvcSpaceSupportMod   int     `xml:"SvcSpaceSupportMod"`
	SvcTimeSupportMode   int     `xml:"SvcTimeSupportMode"`
}

// MarshalCatalogQuery creates XML for catalog query
func MarshalCatalogQuery(sn int, deviceID string) ([]byte, error) {
	query := CatalogQuery{
		CmdType:  "Catalog",
		SN:       sn,
		DeviceID: deviceID,
	}
	return xml.Marshal(query)
}
