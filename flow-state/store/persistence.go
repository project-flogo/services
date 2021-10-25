package store

type Persistence struct {
	PersistenceType string                 `json:"type"`
	Settings        map[string]interface{} `json:"settings"`
}
