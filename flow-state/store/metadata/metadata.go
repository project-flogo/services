package metadata

import "github.com/project-flogo/flow/state"

type Metadata struct {
	Username, AppName, AppVersion, HostId, FlowName, Offset, Limit, Status, Interval, FlowInstanceId string
}

type FlowRecord struct {
	Count    int32
	FlowData []*state.FlowInfo
}
