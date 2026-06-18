package carefulness

type Gone struct {
}

var ErrGone Gone

func (e Gone) Error() string {
	return "Gone"
}

func (e Gone) JSONError() JSONError {
	return JSONError{Error: e.Error()}
}
