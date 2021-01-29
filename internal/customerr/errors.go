package customerr

type Err interface {
	Error() string
}

type OutOfBoundsError struct {
	Err error
	Op  string
}

func (r *OutOfBoundsError) Error() string {
	return r.Err.Error()
}

func (r *OutOfBoundsError) Operation() string {
	return r.Op
}

type ConversionError struct {
	Err error
	Op  string
}

func (r *ConversionError) Error() string {
	return r.Err.Error()
}

func (r *ConversionError) Operation() string {
	return r.Op
}

type UUIDError struct {
	Err error
	Op  string
}

func (r *UUIDError) Error() string {
	return r.Err.Error()
}

func (r *UUIDError) Operation() string {
	return r.Op
}
