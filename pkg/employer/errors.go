package employer

type ErrSchemaNotFound struct {
	ID string
}

func (e ErrSchemaNotFound) Error() string {
	return "schema with id " + e.ID + " not found"
}

type ErrEmployerNotFound struct {
	ID string
}

func (e ErrEmployerNotFound) Error() string {
	return "employer with id (" + e.ID + ") not found"
}

type ErrInvalidReader struct {
	InternalErr error
}

func (e ErrInvalidReader) Error() string {
	return "invalid reader: " + e.InternalErr.Error()
}

func (e ErrInvalidReader) Unwrap() error {
	return e.InternalErr
}
