package metadata

import "github.com/project-flogo/flow/state"

type Metadata struct {
	Username, AppId, HostId, FlowName, Offset, Limit, Status string
}

type FlowRecord struct {
	Count    int32
	FlowData []*state.FlowInfo
}
