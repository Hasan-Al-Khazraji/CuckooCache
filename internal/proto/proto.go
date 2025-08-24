package proto

const (
	OpGet byte = 0x01
	OpSet byte = 0x02
)

const (
	StatusOK byte = iota
	StatusNotFound
	StatusErr
)

type Request struct {
	Op    byte
	Key   []byte
	Value []byte
}

type Response struct {
	Status byte
	Value  []byte
}

func OK(v []byte) *Response {
	return &Response{
		Status: StatusOK,
		Value:  v,
	}
}

func Err() *Response {
	return &Response{
		Status: StatusErr,
	}
}

func NotFound() *Response {
	return &Response{
		Status: StatusNotFound,
	}
}
