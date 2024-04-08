package rest

type resp struct {
	Item any   `json:"item"`
	Err  error `json:"error"`
}

type ErrWarn struct {
	msg string
}

func (warn ErrWarn) Error() string {
	return warn.msg
}

func Warn(msg string) ErrWarn {
	return ErrWarn{msg: msg}
}
