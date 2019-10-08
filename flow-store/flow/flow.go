package flow

type Metdata struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	CreationDate string `json:"creationDate"`
}

type Flow struct {
	Metdata *Metdata    `json:"metdata"`
	Flow    interface{} `json:"flow"`
}
