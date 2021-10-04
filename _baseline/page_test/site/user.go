package site

// User is the standard user type shared between
// page handers, utilities, etc... The site package
// is a good place for this kind of definition.
type User struct {
	Name  string
	Email string
}
