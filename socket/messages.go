package socket

type MessageKind string

var Get MessageKind = "GET"
var Set MessageKind = "SET"

type Request struct {
	Kind  MessageKind `json:"kind"`
	Key   string      `json:"key"`
	Value string      `json:"value"`
}

type Response struct {
	Kind  MessageKind `json:"kind"`
	Value string      `json:"value"`
	OK    bool        `json:"OK"`
}
