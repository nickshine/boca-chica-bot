package db

type ErrItemUnchanged struct {
	message string
}

func NewErrItemUnchanged() *ErrItemUnchanged {
	return &ErrItemUnchanged{
		message: "Item exists and is unchanged",
	}
}

func (e *ErrItemUnchanged) Error() string {
	return e.message
}
