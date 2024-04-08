package rest

type resp struct {
	Item any   `json:"item"`
	Err  error `json:"error"`
}
