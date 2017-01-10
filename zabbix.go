package zabbix

const (
	OK = iota
	Problem
)

const (
	Unclassified = iota
	Info
	Warning
	Average
	High
	Disaster
)

var Priorities = []string{
	"unclassified",
	"info",
	"warning",
	"average",
	"high",
	"disaster",
}
