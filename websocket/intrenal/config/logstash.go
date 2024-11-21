package config

import "os"

type Logstash struct {
	Enable bool
	Addr   string
}

func NewLogstash() Logstash {

	addr := os.Getenv("LOGSTASH_ADDR")
	if addr == "" {
		return Logstash{
			Enable: false,
			Addr:   addr,
		}
	}

	return Logstash{
		Enable: true,
		Addr:   addr,
	}
}
