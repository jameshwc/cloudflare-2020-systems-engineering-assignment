package profile

type Request struct {
	URL    string
	Method string
}

func NewRequest(url, method string) *Request {
	return &Request{url, method}
}
