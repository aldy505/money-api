package out

type Err struct {
	Error string `json:"error" xml:"error" form:"error"`
}

type Msg struct {
	Message string `json:"message" xml:"message" form:"message"`
}
