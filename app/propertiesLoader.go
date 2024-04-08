package app

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/magiconair/properties"
)

var propLog = log.New(io.Discard, "props:", log.LstdFlags|log.Lmsgprefix)

var props properties.Properties

func loadProperties() (*properties.Properties, error) {
	propLog.Println(os.UserHomeDir())
	propLog.Println(os.UserCacheDir())
	propLog.Println(os.UserConfigDir())

	// loading config path from user cache
	localDr, err := getLudoDir()
	if err != nil {
		return nil, err
	}
	propLog.Println(localDr)

	ludoCfg := filepath.Join(localDr, ludoPropertiesFile)

	if _, err := os.Stat(ludoCfg); os.IsNotExist(err) {
		propLog.Println("properties file does not exist")
	}

	// load prop file
	prop, err := properties.LoadFile(ludoCfg, properties.UTF8)

	// if does not exist, try create it
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.New("error in loading properties file")
		}

		propLog.Println("Loading default props")
		prop = properties.NewProperties()
		prop.Set(pipedApiKey, defaultPipedApi)
		prop.Set(instanceListApiKey, defaultInstanceListApi)
		if err := createPropertiesFile(prop, ludoCfg); err != nil {
			return nil, err
		}

	}

	// logging proprties loaded
	propLog.Println("Properties loaded")
	for _, key := range prop.Keys() {
		propLog.Println(key, ":", prop.MustGetString(key))
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

	propLog.Println("Properties created:", n)
	return nil
}

func getLudoDir() (string, error) {
	localDr, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	ludoDir := filepath.Join(localDr, ludoBaseDir)
	propLog.Println(ludoDir)

	// looking for path in ENV
	if path, ok := os.LookupEnv("LUDO_BASE_PATH"); ok {
		ludoDir = path
	}

	return ludoDir, nil
}
