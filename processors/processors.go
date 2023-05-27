package processors

const NOT_LOGGED_IN = "not logged in, login or register to continue"

type NotLoggedInError struct{}

func (err *NotLoggedInError) Error() string {
	return NOT_LOGGED_IN
}

func NewNotLoggedInError() error {
	return &NotLoggedInError{}
}
