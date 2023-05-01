package env

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
)

type Env struct {
	SMTPEmail    string `json:"SMTP_USER_MAME"`
	SMTPPassword string `json:"SMTP_PASSWORD"`
}

var once sync.Once
var env Env

func GetEnv() *Env {
	once.Do(func() {
		data, err := os.ReadFile(".env")
		if err != nil {
			panic(err)
		}

		mapEnv := make(map[string]string)
		for _, line := range strings.Split(string(data), "\n") {
			kv := strings.SplitN(line, "=", 2)
			if len(kv) != 2 {
				continue
			}

			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])

			mapEnv[key] = value
		}
		serialized, err := json.Marshal(mapEnv)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(serialized, &env)
		if err != nil {
			panic(err)
		}
	})

	return &env
}
