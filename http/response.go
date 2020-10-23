package http

type Response struct {
	ID         uint64
	ChanID     uint64
	Time       uint64
	IsSucceed  bool
	StatusCode int
	Size       int
}

func NewResponse(ID, chanID, time uint64, isSucceed bool, statusCode, size int) *Response {
	return &Response{ID, chanID, time, isSucceed, statusCode, size}
}
