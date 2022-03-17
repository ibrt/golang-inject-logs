package logz

// User describes a user for logs purposes.
type User struct {
	ID       string
	Email    string
	Metadata map[string]interface{}
}

// UserExtractor describes the ability to provide a user for logs purposes.
type UserExtractor interface {
	ExtractUser() *User
}
