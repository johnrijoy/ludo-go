package app

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var audioCache CacheStore

type CacheStore struct {
	isEnabled bool
	cacheDir  string
	cacheMap  map[string]string
}

var cacheLog = log.New(os.Stdout, "cacheStore: ", log.LstdFlags|log.Lmsgprefix)

func (cache *CacheStore) Init(cacheDir string) error {
	if !cache.isEnabled {
		cacheLog.Println("Caching is disabled")
		return nil
	}

	cacheLog.Println("Cache at:", cacheDir)
	cache.cacheDir = cacheDir
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}
	cmap, err := cache.buildCacheMap()
	if err != nil {
		return err
	}
	cache.cacheMap = cmap
	return nil
}

func (cache *CacheStore) Close() error {
	return nil
}

func (cache *CacheStore) CacheAudio(audio AudioDetails) {
	if !cache.isEnabled {
		cacheLog.Println("Caching is disabled")
		return
	}

	trackTitle := audio.Title
	fileName := audio.YtId
	audioStreamUrl := audio.AudioStreamUrl
	fileLoc := filepath.Join(cache.cacheDir, fileName)

	// check if file does not exist
	if _, ok := cache.LookupCache(audio.AudioBasic); !ok {
		go func() {
			err := downloadFile(fileLoc, audioStreamUrl)
			if err != nil {
				cacheLog.Println("Error in downloading file:", trackTitle, "|", fileName)
				cacheLog.Println(err)
			}
		}()
	}
}

func (cache *CacheStore) LookupCache(audio AudioBasic) (string, bool) {
	if !cache.isEnabled {
		cacheLog.Println("Caching is disabled")
		return "", false
	}

	cachePath, ok := cache.cacheMap[audio.YtId]

	if !ok {
		fileLoc := filepath.Join(cache.cacheDir, audio.YtId+".m4a")
		cachePath, ok = searchCacheDir(fileLoc)
		if !ok {
			return "", false
		}
		cache.cacheMap[audio.YtId] = cachePath
	}

	cacheLog.Println("Audio Cached at:", cachePath)
	return cachePath, true
}

func (cache *CacheStore) buildCacheMap() (map[string]string, error) {
	cacheMap := make(map[string]string)
	matches, err := filepath.Glob(filepath.Join(cache.cacheDir, "*.m4a"))
	if err != nil {
		return nil, err
	}

	cacheLog.Println(matches)
	for _, match := range matches {
		fileId := strings.Split(filepath.Base(match), ".")[0]
		cacheMap[fileId] = match
	}

	cacheLog.Println(cacheMap)
	return cacheMap, nil
}

func searchCacheDir(fileLoc string) (string, bool) {
	if fd, err := os.Stat(fileLoc); err != nil || fd.IsDir() {
		return "", false
	}
	return fileLoc, true
}

func downloadFile(filePath, fileUrl string) error {
	// intialise download client
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	// get response
	resp, err := client.Get(fileUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// check response file formate
	if resp.Header.Get("Content-Type") == "audio/mp4" {
		filePath += ".m4a"
	}
	tmpFilePath := filePath + ".tmp"

	// create temp file
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	cacheLog.Println("File created")

	// download file
	size, err := io.Copy(file, resp.Body)
	fmt.Println()
	if err != nil {
		return err
	}
	file.Close()

	cacheLog.Printf("Downloaded a file %s with size %d\n", filePath, size)
	// rename temp file and remove incase of any error
	err = os.Rename(tmpFilePath, filePath)
	os.Remove(tmpFilePath)
	return err
}
