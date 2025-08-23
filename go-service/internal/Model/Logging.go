package model

import "time"

type Logging struct {
	Id        string    `json:"id"`
	TimeStamp time.Time `json:"time_stamp"`
	Service   string    `json:"service" validate:"required"`
	Level     string    `json:"level" validate:"required"`
	Message   string    `json:"message" `
	MetaData  MetaData  `json:"meta_data" `
}

type MetaData struct {
	SourceIp string `json:"source_ip"`
	Region   string `json:"region" validate:"required"`
}
