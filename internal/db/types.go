package db

// ItemUnchangedError is an implementation of the error interface used to signify items in the db
// exist and have no changes.
type ItemUnchangedError struct {
	message string
}

// NewItemUnchangedError returns a new instance of ItemUnchangedError.
func NewItemUnchangedError() *ItemUnchangedError {
	return &ItemUnchangedError{
		message: "Item exists and is unchanged",
	}
}

func (e *ItemUnchangedError) Error() string {
	return e.message
}
