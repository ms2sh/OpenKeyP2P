package p2p

import "time"

type _ConnType uint8

type Role int

type VersionInfo struct {
	SupportedVersions []string `json:"supported_versions"`
	AgreedVersion     string   `json:"agreed_version,omitempty"`
}

const (
	// BigDataWrapper
	headerSize                    = 21
	fragmentSize                  = 1024
	maxWriteRetries               = 3
	retryBaseDelay  time.Duration = 1 * time.Second
	//maxFrames       int           = (10 * 1024 * 1024) / (fragmentSize - headerSize)
	TCP        _ConnType = 254
	TLS        _ConnType = 253
	UnixSocket _ConnType = 252

	// Unix Sockett
	unixSocketFrameHeaderSize  = 12
	unixSocketAckSize          = 9
	unixSocketFramePayloadSize = fragmentSize
	unixSocketMaxFrameSize     = unixSocketFrameHeaderSize + unixSocketFramePayloadSize

	Client Role = iota
	Server
)
