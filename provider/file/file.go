package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"

	"gopkg.in/yaml.v2"
)

// Provider file provider
type Provider struct {
	path string
}

// NewProvider returns new FileProvider
func NewProvider(path string) *Provider {
	return &Provider{path: path}
}

func (fp *Provider) Read() (conf map[string]interface{}, err error) {
	var content []byte
	content, err = ioutil.ReadFile(fp.path)
	if err != nil {
		return
	}

	fType := filepath.Ext(fp.path)
	conf = make(map[string]interface{})
	switch fType {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(content, &conf)
	case ".json":
		err = json.Unmarshal(content, &conf)
	case ".init", ".toml":
		fallthrough
	default:
		err = fmt.Errorf("File type %v unsupported", fType)
	}
	return
}

// Watch file and automate update
func (fp *Provider) Watch(watcher func(map[string]interface{})) {
	m := make(map[string]interface{})
	watcher(m)

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer w.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-w.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file: ", event.Name)
				}
			case err := <-w.Errors:
				log.Println("error: ", err)
			}
		}
	}()

	err = w.Add(fp.path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
