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
	vlcPlayer VlcPlayer
	props     properties.Properties
}

func (app *AppContext) Init() error {
	props, err := loadProperties()
	if err != nil {
		return err
	}
	app.props = *props
	setPipedConfig(props)

	if err := app.vlcPlayer.InitPlayer(); err != nil {
		return err
	}
	return nil
}

func (app *AppContext) Close() error {
	if err := app.vlcPlayer.ClosePlayer(); err != nil {
		return err
	}
	return nil
}

func (app *AppContext) VlcPlayer() *VlcPlayer {
	return &app.vlcPlayer
}

func loadProperties() (*properties.Properties, error) {
	appLog.Println(os.UserHomeDir())
	appLog.Println(os.UserCacheDir())
	appLog.Println(os.UserConfigDir())
	localDr, _ := os.UserCacheDir()
	ludoCfg := filepath.Join(localDr, ludoDir, ludoPropertiesFile)
	appLog.Println(ludoCfg)

	if _, err := os.Stat(ludoCfg); os.IsNotExist(err) {
		appLog.Println("properties file does not exist")
	}

	prop, err := properties.LoadFile(ludoCfg, properties.UTF8)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.New("error in loading properties file")
		}

		appLog.Println("Loading default props")
		prop = properties.MustLoadString(defaultProp)
		if err := createPropertiesFile(prop, ludoCfg); err != nil {
			return nil, err
		}

	}

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
