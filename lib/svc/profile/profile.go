package profile

import (
	"encoding/json"
	"log"
	"os"
)

func GetProfileFile(confDir string, name string) string {
	return confDir + "/" + name + ".json"
}

func LoadProfile(file string) (info Info, err error) {
	log.Printf("Load profile file: %s", file)
	f, err := os.Open(file)
	if err != nil {
		log.Println("Unable to load profile file")
		return
	}
	err = json.NewDecoder(f).Decode(&info)
	if err != nil {
		log.Println("Unable to decode profile file")
	}
	return
}

func (info Info) SaveProfile(file string) (err error) {
	log.Printf("Save profile file: %s", file)
	f, err := os.Create(file)
	if err != nil {
		log.Println("Unable to create profile file")
		return
	}
	err = json.NewEncoder(f).Encode(info)
	if err != nil {
		log.Println("Unable to encode to profile file")
	}
	return
}
