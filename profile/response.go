package profile

type Response struct {
	ID         uint64
	ChanID     uint64
	Time       uint64
	Size       uint64
	IsSucceed  bool
	StatusCode int
}

func NewResponse(ID, chanID, time, size uint64, isSucceed bool, statusCode int) *Response {
	return &Response{ID, chanID, time, size, isSucceed, statusCode}
}
