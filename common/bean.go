package common

import (
	"time"
)

type Event struct {
	Image        string `json:"image"`
	ImageDigest  string `json:"imageDigest"`
	AppId        int    `json:"appId"`
	EnvId        int    `json:"envId"`
	PipelineId   int    `json:"pipelineId"`
	CiArtifactId int    `json:"ciArtifactId"`
	UserId       int    `json:"userId"`
	AccessKey    string `json:"accessKey"`
	SecretKey    string `json:"secretKey"`
	Token        string `json:"token"`
}

type TelemetryEvent struct {
	Id          int32     `json:"id" validate:"number"`
	UPID        string    `json:"upid"`
	ActiveSince time.Time `json:"activeSince"`
	LastActive  time.Time `json:"lastActive"`
	UserId      int32     `json:"-"` // created or modified telemetry id
}
