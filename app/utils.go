package app

const (
	defaultProp = "config.piped.apiUrl=https://pipedapi.kavin.rocks\nconfig.piped.instanceListApi=https://piped-instances.kavin.rocks"
	Version     = "0.1.0"
)

// helper string
const (
	ludoDir                = "ludo"
	ludoPropertiesFile     = "ludo.properties"
	defaultCacheDir        = "cache"
	defaultDataStore       = "ludo.dbquit"
	pipedUrlKey            = "config.piped.apiUrl"
	instanceListUrlKey     = "config.piped.instanceListApi"
	defaultPipedUrl        = "https://pipedapi.kavin.rocks"
	defaultInstanceListApi = "https://piped-instances.kavin.rocks"
)

// Helpers //

func trimList[T any](inputList []T, offset int, limit int) []T {
	outputList := inputList

	if offset > 0 && offset < len(inputList) {
		outputList = outputList[offset:]
	}

	if limit > 0 && offset >= 0 && limit < len(inputList)-offset {
		outputList = outputList[:limit]
	}

	return outputList
}
