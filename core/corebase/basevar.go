package corebase


type ProcessWhiteList struct {
	ProcessWhitelist	[]ProcessWhite	`json:"process_whitelist"`
}

type ProcessWhite struct {
	ProcessName string	`json:"process_name"`
	ProcessFileHash	string	`json:"process_file_hash"`
}