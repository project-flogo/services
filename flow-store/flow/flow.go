package flow

type Metadata struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	CreationDate string `json:"creationDate"`
}

type Flow struct {
	Metadata *Metadata   `json:"metadata"`
	Flow     interface{} `json:"flow"`
}
