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

type AppContext struct {
	vlcPlayer  VlcPlayer
	props      properties.Properties
	audioDb    AudioDatastore
	audioCache CacheStore
}

func (app *AppContext) Init() error {
	// Load properties file
	props, err := loadProperties()
	if err != nil {
		return err
	}
	app.props = *props

	// Set Piped config
	setPipedConfig(props)

	// Load database
	localDr, _ := getLudoDir()
	dbPath := props.GetString(dataStoreKey, localDr)

	if err := app.audioDb.InitDb(dbPath); err != nil {
		return err
	}

	// load Cache
	defCachePath := filepath.Join(localDr, defaultCacheDir)
	cachePath := props.GetString(cacheDirKey, defCachePath)
	if err := app.audioCache.Init(cachePath); err != nil {
		return err
	}

	// load audio player
	if err := app.vlcPlayer.InitPlayer(&app.audioDb, &app.audioCache); err != nil {
		return err
	}
	return nil
}

func (app *AppContext) Close() error {
	if err := app.vlcPlayer.ClosePlayer(); err != nil {
		return err
	}
	localDr, _ := getLudoDir()
	dumpPath := filepath.Join(localDr, "dumps.json")

	app.audioDb.db.ExportCollection(audioDocCollection, dumpPath)

	if err := app.audioDb.CloseDb(); err != nil {
		return err
	}

	app.audioCache.Close()

	return nil
}

func (app *AppContext) VlcPlayer() *VlcPlayer {
	return &app.vlcPlayer
}

func (app *AppContext) AudioDb() *AudioDatastore {
	return &app.audioDb
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
