package http

type Metadata struct {
	FileName       string
	ContentLength  uint64
	ContentType    string
	Url            string
	SupportPartial bool
}

type Dictionary map[string]string

type RequestStruct struct {
	Url     string
	Method  string
	Headers Dictionary
	Query   Dictionary
	Body    []byte
}
