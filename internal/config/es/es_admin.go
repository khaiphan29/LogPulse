package esconfig

import (
   "os"
   "encoding/json"
)

type IndexConfig struct {
   Name string `json:"name"`
   Mapping string `json:"mapping"`
}

type ClientConfig struct {
   Host string
   Port string
}

func LoadClientConfig(prefix string) *ClientConfig {
   return &ClientConfig{
      Host: os.Getenv(prefix + "_HOST"),
      Port: os.Getenv(prefix + "_PORT"),
   }
}

func LoadIndexConfig(path string) *IndexConfig {
   data, err := os.ReadFile(path)
   if err != nil {
      panic(err)
   }

   var config IndexConfig
   json.Unmarshal(data, &config)

   return &config
}
