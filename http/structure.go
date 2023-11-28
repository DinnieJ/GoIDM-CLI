package http

type Metadata struct {
	FileName       string
	ContentLength  uint64
	ContentType    string
	Url            string
	SupportPartial bool
}
