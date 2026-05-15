package gb28181

// GB28181 standard constants
const (
	// GB code length
	GBCodeLength = 20

	// Center code length (first 8 digits)
	CenterCodeLength = 8

	// Industry code length (2 digits after center)
	IndustryCodeLength = 2

	// Type code length (3 digits)
	TypeCodeLength = 3

	// Network code length (1 digit)
	NetworkCodeLength = 1

	// Serial number length (6 digits)
	SerialNumberLength = 6
)

// Device type codes (GB/T 28181)
const (
	DeviceTypeDVR      = "111"
	DeviceTypeNVR      = "118"
	DeviceTypeCamera   = "131"
	DeviceTypeIPC      = "132"
	DeviceTypeGroup    = "215"
	DeviceTypeVirtual  = "216"
	DeviceTypeCenter   = "300"
)

// Industry code types
const (
	IndustrySocialSecurityRoad     = "00"
	IndustrySocialSecurityCommunity = "01"
	IndustryTrafficRoad            = "04"
	IndustryTrafficBayonet         = "05"
	IndustryCityManagement         = "08"
	IndustryEmergencyManagement    = "09"
	IndustryPeopleAirDefense       = "10"
	IndustrySafetyProduction       = "11"
	IndustryEnvironmentalProtection = "12"
	IndustryWaterConservancy       = "13"
	IndustryNaturalResources       = "14"
	IndustryHealthCommission       = "15"
	IndustryMarketRegulation       = "16"
	IndustryTransportation         = "17"
	IndustryAgricultureRural       = "18"
	IndustryCultureTourism         = "19"
	IndustryJustice               = "20"
	IndustryEducation              = "21"
	IndustryForestry               = "22"
)

// Network identification types
const (
	NetworkPublicSecurityVideo = "0"
	NetworkIndustrySpecific    = "2"
	NetworkPublicInternet      = "7"
)
