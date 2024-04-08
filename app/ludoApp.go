package app

import (
	"io"
	"log"
	"path/filepath"
)

var appLog = log.New(io.Discard, "App:", log.LstdFlags|log.Lmsgprefix)

func Init() error {
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
	return nil
}

func Close() error {
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

	return nil
}

func MediaPlayer() *VlcPlayer {
	return &vlcPlayer
}

func AudioDb() *AudioDatastore {
	return &audioDb
}

// App Functions
