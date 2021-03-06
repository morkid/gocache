package gocache

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// DiskCacheConfig struct
type DiskCacheConfig struct {
	Directory string
	ExpiresIn time.Duration
}

type diskCache struct {
	Config DiskCacheConfig
}

type diskCacheContent struct {
	CreatedAt time.Time `json:"created_at"`
	Value     string    `json:"value"`
}

// NewDiskCache store cache to file on disk
func NewDiskCache(config DiskCacheConfig) *AdapterInterface {
	if config.Directory == "" {
		p, err := os.Getwd()
		if nil == err {
			config.Directory = p
		}
	} else {
		p, err := filepath.Abs(config.Directory)
		if nil == err {
			config.Directory = p
		}
	}
	if config.ExpiresIn <= 0 {
		config.ExpiresIn = 3600 * time.Second
	}

	var adapter AdapterInterface = diskCache{
		Config: config,
	}

	return &adapter
}

func (d diskCache) getDirectory() (string, error) {
	if d.Config.Directory == "" {
		return "", errors.New("Invalid cache directory")
	}

	if !dirExists(d.Config.Directory) {
		dir, err := filepath.Abs(d.Config.Directory)
		if nil != err {
			return "", err
		}
		d.Config.Directory = dir
		if err := os.MkdirAll(d.Config.Directory, 0755); nil != err {
			return "", err
		}
	}

	return d.Config.Directory, nil
}

func (d diskCache) pathOf(filename string) string {
	dir, err := d.getDirectory()
	if nil != err {
		log.Printf("Error: %s", err.Error())
		return ""
	}

	return path.Join(dir, filename)
}

func (d diskCache) isValidKey(key string) bool {
	if ok, err := regexp.Match(`[A-z0-9_]+`, []byte(key)); nil == err && ok {
		return true
	}

	return false
}

func (d diskCache) Set(key string, value string) error {
	if d.isValidKey(key) {
		filename := key + ".json"
		if absfile := d.pathOf(filename); absfile != "" {
			content := diskCacheContent{
				CreatedAt: time.Now(),
				Value:     value,
			}
			data, err := json.Marshal(content)
			if nil != err {
				return err
			}

			return ioutil.WriteFile(absfile, data, 0644)
		}

		return errors.New("Error: Cannot write data to cache directory")
	}

	return errors.New("Invalid cache key characters, key only allow regex [A-z0-9_]+")
}

func (d diskCache) Get(key string) (string, error) {
	if d.isValidKey(key) {
		filename := key + ".json"
		if absfile := d.pathOf(filename); absfile != "" && fileExists(absfile) {
			data, err := ioutil.ReadFile(absfile)
			if nil != err {
				return "", err
			}
			var content diskCacheContent
			err = json.Unmarshal(data, &content)
			if nil != err {
				return "", err
			}
			now := time.Now()
			diff := now.Sub(content.CreatedAt)
			if diff > d.Config.ExpiresIn {
				return "", errors.New("Cache has been expired")
			}

			return content.Value, nil
		}

		return "", errors.New("Error: Cannot write data to cache directory")
	}

	return "", errors.New("Invalid cache key characters, key only allow regex [A-z0-9_]+")
}

func (d diskCache) IsValid(key string) bool {
	if d.isValidKey(key) {
		filename := key + ".json"
		if absfile := d.pathOf(filename); absfile != "" && fileExists(absfile) {
			data, err := ioutil.ReadFile(absfile)
			if nil != err {
				return false
			}
			var content diskCacheContent
			err = json.Unmarshal(data, &content)
			if nil != err {
				return false
			}
			now := time.Now()
			diff := now.Sub(content.CreatedAt)
			if diff > d.Config.ExpiresIn {
				return false
			}

			return true
		}
	}

	return false
}

func (d diskCache) Clear(key string) error {
	if d.isValidKey(key) {
		filename := key + ".json"
		if absfile := d.pathOf(filename); absfile != "" && fileExists(absfile) {
			defer func(absfile string) {
				if err := os.Remove(absfile); nil != err {
					log.Println(err)
				}
			}(absfile)

			return nil
		}
	}

	return errors.New("Invalid cache key characters, key only allow regex [A-z0-9_]+")
}

func (d diskCache) ClearPrefix(keyPrefix string) error {
	dir, err := d.getDirectory()
	if nil != err {
		return err
	}

	defer func(dir string, keyPrefix string) {
		err := filepath.Walk(dir, func(filename string, info os.FileInfo, err error) error {
			if nil == err &&
				!info.IsDir() &&
				strings.HasPrefix(info.Name(), keyPrefix) &&
				strings.HasSuffix(info.Name(), ".json") {
				log.Println(filename)
				return os.Remove(filename)
			}

			return nil
		})
		if nil != err {
			log.Println(err)
		}
	}(dir, keyPrefix)

	return nil
}

func (d diskCache) ClearAll() error {
	dir, err := d.getDirectory()
	if nil != err {
		return err
	}

	defer func(dir string) {
		if err := os.RemoveAll(dir); nil != err {
			log.Println(err)
			return
		}

		_, err = d.getDirectory()
		if nil != err {
			log.Println("Cannot recreate cache directory")
		}
	}(dir)

	return nil
}

func dirExists(dirname string) bool {
	info, err := os.Stat(filepath.FromSlash(dirname))
	if nil != err {
		return false
	}
	return !os.IsNotExist(err) && info.IsDir()
}

func fileExists(filename string) bool {
	info, err := os.Stat(filepath.FromSlash(filename))
	if nil != err {
		return false
	}
	return !os.IsNotExist(err) && !info.IsDir()
}
