package models

type ErrEOF struct {
}

func (e *ErrEOF) Error() string {
	return "EOF"
}
