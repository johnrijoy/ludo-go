package app

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/magiconair/properties"
)

var appLog = log.New(io.Discard, "App:", log.LstdFlags|log.Lmsgprefix)

var props properties.Properties
var isRunning = false

func Init() error {
	if isRunning {
		return nil
	}
	// Load properties file
	lprops, err := loadProperties()
	if err != nil {
		return err
	}
	props = *lprops

	// Set Piped config
	setPipedConfig(props)

	// Load database
	localDr, _ := getLudoDir()
	dbPath := props.GetString(dataStoreKey, localDr)

	if err := audioDb.InitDb(dbPath); err != nil {
		return err
	}

	// load Cache
	defCachePath := filepath.Join(localDr, defaultCacheDir)
	cachePath := props.GetString(cacheDirKey, defCachePath)
	audioCache.isEnabled = props.GetBool(isCacheEnabledKey, true)
	if err := audioCache.Init(cachePath); err != nil {
		return err
	}

	// load audio player
	if err := vlcPlayer.InitPlayer(); err != nil {
		return err
	}

	isRunning = true
	return nil
}

func Close() error {
	if !isRunning {
		return nil
	}
	if err := vlcPlayer.ClosePlayer(); err != nil {
		return err
	}
	localDr, _ := getLudoDir()
	dumpPath := filepath.Join(localDr, "dumps.json")

	audioDb.db.ExportCollection(audioDocCollection, dumpPath)

	if err := audioDb.CloseDb(); err != nil {
		return err
	}

	audioCache.Close()

	isRunning = false
	return nil
}

func MediaPlayer() *VlcPlayer {
	return &vlcPlayer
}

func AudioDb() *AudioDatastore {
	return &audioDb
}

func loadProperties() (*properties.Properties, error) {
	appLog.Println(os.UserHomeDir())
	appLog.Println(os.UserCacheDir())
	appLog.Println(os.UserConfigDir())

	// loading config path from user cache
	localDr, err := getLudoDir()
	if err != nil {
		return nil, err
	}
	appLog.Println(localDr)

	ludoCfg := filepath.Join(localDr, ludoPropertiesFile)

	if _, err := os.Stat(ludoCfg); os.IsNotExist(err) {
		appLog.Println("properties file does not exist")
	}

	// load prop file
	prop, err := properties.LoadFile(ludoCfg, properties.UTF8)

	// if does not exist, try create it
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.New("error in loading properties file")
		}

		appLog.Println("Loading default props")
		prop = properties.NewProperties()
		prop.Set(pipedApiKey, defaultPipedApi)
		prop.Set(instanceListApiKey, defaultInstanceListApi)
		if err := createPropertiesFile(prop, ludoCfg); err != nil {
			return nil, err
		}

	}

	// logging proprties loaded
	appLog.Println("Properties loaded")
	for _, key := range prop.Keys() {
		appLog.Println(key, ":", prop.MustGetString(key))
	}

	return prop, nil
}

func createPropertiesFile(prop *properties.Properties, ludoCfg string) error {
	if err := os.MkdirAll(filepath.Dir(ludoCfg), 0755); err != nil {
		return err
	}

	file, err := os.Create(ludoCfg)
	if err != nil {
		return err
	}

	n, err := prop.Write(file, properties.UTF8)
	if err != nil {
		return err
	}

	appLog.Println("Properties created:", n)
	return nil
}

func getLudoDir() (string, error) {
	localDr, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	ludoDir := filepath.Join(localDr, ludoBaseDir)
	appLog.Println(ludoDir)

	// looking for path in ENV
	if path, ok := os.LookupEnv("LUDO_BASE_PATH"); ok {
		ludoDir = path
	}

	return ludoDir, nil
}
