package models

type Tx struct {
	Hash      string `json:"hash,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}
