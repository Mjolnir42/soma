package msg

import "github.com/mjolnir42/soma/lib/proto"

// Authorization struct
type Authorization struct {
	AuthUser     string
	RemoteAddr   string
	Section      string
	Action       string
	TeamID       string
	OncallID     string
	MonitoringID string
	CapabilityID string
	RepositoryID string
	BucketID     string
	GroupID      string
	ClusterID    string
	NodeID       string
	Grant        *proto.Grant
}
