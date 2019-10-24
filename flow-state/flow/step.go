package flow

type StepInfo struct {
	FlowID   string    `json:"flowID"`
	ID       int64     `json:"id"`
	FlowURI  string    `json:"flowURI"`
	StepData *StepData `json:"stepData"`
	Status   int64     `json:"status"`
	State    int64     `json:"state"`
	Date     string    `json:"date"`
}

type StepData struct {
	Status    int64       `json:"status"`
	State     int64       `json:"state"`
	WqChanges []*WqChange `json:"wqChanges"`
	TdChanges []*TdChange `json:"tdChanges"`
	LdChanges []*LdChange `json:"ldChanges"`
	Attrs     []*Attr     `json:"attrs"`
}

type Snapshot struct {
	ID           int64         `json:"id"`
	FlowID       string        `json:"flowID"`
	State        int64         `json:"state"`
	Status       int64         `json:"status"`
	Date         string        `json:"date"`
	SnapshotData *SnapshotData `json:"SnapshotData"`
}

type SnapshotData struct {
	ID          string       `json:"id"`
	State       int64        `json:"state"`
	Status      int64        `json:"status"`
	Attrs       []*Attribute `json:"attrs"`
	FlowURI     string       `json:"flowUri"`
	WorkQueue   []*WorkItem  `json:"workQueue"`
	RootTaskEnv *RootTaskEnv `json:"rootTaskEnv"`
}

type RootTaskEnv struct {
	ID        interface{} `json:"id"`
	TaskId    interface{} `json:"taskId"`
	TaskDatas []*TaskData `json:"taskDatas"`
	LinkDatas []*LinkData `json:"linkDatas"`
}

type SnapshotInfo struct {
	ID           int64         `json:"id"`
	FlowID       string        `json:"flowID"`
	State        int64         `json:"state"`
	Status       int64         `json:"status"`
	Snapshot     *SnapshotData `json:"snapshot"`
	Date         string        `json:"date"`
	SnapshotData *SnapshotData `json:"snapshotData"`
}

type Attr struct {
	ChgType int
	Att     Attribute `json:"Attribute"`
}

type Attribute struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type RollUpObj struct {
	ID                 int64         `json:"id"`
	Status             int64         `json:"status"`
	State              int64         `json:"state"`
	FlowURI            string        `json:"flowUri"`
	Snapshot           *SnapshotData `json:"snapshot"`
	IgnoredTaskIds     []interface{}
	IgnoredLinkds      []int64
	IgnoredWorkitemIds []int64
	IgnoredAttrs       []string
}
