package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
)

func RegisterService(r Registration) error {
	serviceUpdateURl, err := url.Parse(r.ServiceUpdateURL)
	if err != nil {
		return err
	}
	http.Handle(serviceUpdateURl.Path, &serviceUpdateHandler{})
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(r)
	if err != nil {
		return err
	}
	res, err := http.Post(ServicesURL, "application/json", buf)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register code: %d", res.StatusCode)
	}
	return nil
}
func ShowDownService(url string) error {
	req, err := http.NewRequest(http.MethodDelete, ServicesURL, bytes.NewBuffer([]byte(url)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to deregister service with code %v", res.StatusCode)

	}
	return nil
}

type providers struct {
	services map[ServiceName][]string
	mutex    *sync.RWMutex
}

func (p *providers) Update(pat patch) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for _, patchEntry := range pat.Added {
		if _, ok := p.services[patchEntry.Name]; !ok {
			p.services[patchEntry.Name] = make([]string, 0)
		}
		p.services[patchEntry.Name] = append(p.services[patchEntry.Name], patchEntry.URL)

	}
	for _, patchEntry := range pat.Removed {
		if provUrls, ok := p.services[patchEntry.Name]; ok {
			for i := range provUrls {
				if provUrls[i] == patchEntry.URL {
					p.services[patchEntry.Name] = append(provUrls[:i], provUrls[i+1:]...)
				}
			}

		}
	}

}
func (p providers) get(name ServiceName) (string, error) {
	provider, ok := p.services[name]
	if !ok {
		return "", fmt.Errorf("NO provider available for service %s", name)
	}
	return provider[0], nil
}
func GetProvider(name ServiceName) (string, error) {
	return prov.get(name)
}

type serviceUpdateHandler struct {
}

func (sub serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	dec := json.NewDecoder(r.Body)
	var p patch
	err := dec.Decode(&p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("Updated received %v\n", p)
	prov.Update(p)
}

var prov = providers{
	services: make(map[ServiceName][]string),
	mutex:    new(sync.RWMutex),
}
