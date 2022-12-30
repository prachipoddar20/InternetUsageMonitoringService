package models

type UsageDetails struct {
	Username        string `json:"username"`
	LastDayUsage    string `json:"lastDayUsage"`
	Last7DaysUsage  string `json:"last7DaysUsage"`
	Last30DaysUsage string `json:"last30DaysUsage"`
}

type Usage struct {
	Time     string `json:"time"`
	Upload   string `json:"upload"`
	Download string `json:"download"`
}

type UserDetails struct {
	Username        string `json:"username"`
	LastHourUsage   Usage  `json:"lastHourUsage"`
	Last6HourUsage  Usage  `json:"last6HourUsage"`
	Last24HourUsage Usage  `json:"last24HourUsage"`
}

type UsageRecord struct {
	Username     string
	MacAddress   string
	StartTime    string
	UsageTime    string
	UploadSize   string
	DownloadSize string
}
