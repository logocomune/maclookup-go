package maclookup

const (
	apiURIPrefix      = "https://api.maclookup.app"
	apiMAC            = "/v2/macs/"
	companyNameSuffix = "/company/name"
	apiKeyParam       = "?apiKey="
)

type maclookupResponseAPIV2 struct {
	Success    bool   `json:"success"`
	Found      bool   `json:"found"`
	MacPrefix  string `json:"macPrefix"`
	Company    string `json:"company"`
	Address    string `json:"address"`
	Country    string `json:"country"`
	BlockStart string `json:"blockStart"`
	BlockEnd   string `json:"blockEnd"`
	BlockSize  int    `json:"blockSize"`
	BlockType  string `json:"blockType"`
	Updated    string `json:"updated"`
	IsRand     bool   `json:"isRand"`
	IsPrivate  bool   `json:"isPrivate"`
}

type errorResponseAPIV2 struct {
	Success   bool   `json:"success" xml:"success"`
	Error     string `json:"error,omitempty" xml:"error,omitempty"`
	ErrorCode int    `json:"errorCode,omitempty" xml:"errorCode,omitempty"`
	MoreInfo  string `json:"moreInfo,omitempty" xml:"moreInfo,omitempty"`
}
