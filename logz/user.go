package logz

// User describes a user.
type User struct {
	ID       string
	Email    string
	Metadata Metadata
}

// UserExtractor describes the ability to provide a user.
type UserExtractor interface {
	ExtractUser() *User
}
