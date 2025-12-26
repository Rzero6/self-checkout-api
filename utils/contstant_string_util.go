package utils

const (
	SessionIDHeader  string = "X-Session-ID"
	SessionIDMessage string = "session_id required"

	CartStatusActive     string = "ACTIVE"
	CartStatusCheckedOut string = "CHECKED OUT"

	TransactionStatusSuccess   string = "SUCCESS"
	TransactionStatusPending   string = "PENDING"
	TransactionStatusExpired   string = "EXPIRED"
	TransactionStatusCancelled string = "CANCELLED"
	TransactionStatusFailed    string = "FAILED"

	MessageStatusSuccess string = "Success"
	MessageStatusFailed  string = "Failed"
)
