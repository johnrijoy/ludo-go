package app

const (
	Version = "0.1.0"
)

// helper string
const (
	ludoBaseDir        = "ludo"
	ludoPropertiesFile = "ludo.props"
	cacheDir           = "cache"
)

// properties file
const (
	dataStoreKey           = "config.database.path"
	cacheDirKey            = "config.cache.path"
	pipedApiKey            = "config.piped.apiUrl"
	defaultPipedApi        = "https://pipedapi.kavin.rocks"
	instanceListApiKey     = "config.piped.instanceListApi"
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
