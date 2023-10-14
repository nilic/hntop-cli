package hntopclient

type Query struct {
	ResultCount int
	StartTime   int64
	EndTime     int64
	FrontPage   bool
	Tags        string
	Query       string
}
