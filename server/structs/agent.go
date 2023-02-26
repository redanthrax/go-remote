package structs

import (
	"github.com/pion/webrtc/v3"
)

type Agent struct {
    ID string `yaml:"ID"`
    ApiKey string `yaml:"ApiKey"`
    RequestDescription webrtc.SessionDescription
    AccessDescription webrtc.SessionDescription
}
