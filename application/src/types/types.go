package types

import "time"

type User struct {
	Id   int64
	Name string
}

type EventType struct {
	Id   int64
	Name string
}

type Event struct {
	Id        int64                  `json:"id"`
	Timestamp string                 `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata"`
	UserId    int64                  `json:"user_id"`
	TypeId    int64                  `json:"type_id"`
}

type Event2 struct {
	Id        int64                  `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata"`
	UserId    int64                  `json:"user_id"`
	TypeId    int64                  `json:"type_id"`
	TypeName  string                 `json:"type"`
}
