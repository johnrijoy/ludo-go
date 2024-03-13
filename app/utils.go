package app

var (
	api_url     = "https://piapi.ggtyler.dev"
	old_api_url string
)

func SetPipedApi(val string) error {
	old_api_url = api_url
	api_url = val
	return nil
}

func GetPipedApi() string {
	return api_url
}

func GetOldPipedApi() string {
	return old_api_url
}
