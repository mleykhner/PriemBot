package models

type DialogStatus string

const (
	DialogStatusOpen   DialogStatus = "open"
	DialogStatusActive DialogStatus = "active"
	DialogStatusClosed DialogStatus = "closed"
)
