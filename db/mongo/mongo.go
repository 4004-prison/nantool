package mongo

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config Config
type Config struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	PoolSize int    `yaml:"poolSize"`
	Database string `yaml:"database"`

	ConnTimeout     time.Duration `yaml:"connTimeout"`
	MaxConnIdleTime time.Duration `yaml:"maxConnIdleTime"`
	MaxOpenConns    int           `yaml:"maxOpenConns"`

	uri string
}

const (
	authTemplate = "%s:%s@"
	uriTemplate  = "mongodb://%s%s%s"
)

var client *mongo.Client
var mu sync.Mutex

func (c *Config) fill() error {
	auth := ""
	if c.User != "" && c.Password != "" {
		auth = fmt.Sprintf(authTemplate, c.User, c.Password)
	} else if c.User == "" || c.Password == "" {
		return errors.New("the mongo username or password cannot be empty")
	}
	if c.Port != "" {
		c.Port = ":" + c.Port
	}
	c.uri = fmt.Sprintf(uriTemplate, auth, c.Host, c.Port)
	return nil
}

func New(config *Config) (*mongo.Client, error) {
	if client == nil {
		mu.Lock()
		defer mu.Unlock()
		if client == nil {
			err := config.fill()
			if err != nil {
				return nil, err
			}
			cliOpt := options.Client().
				SetMaxPoolSize(uint64(config.PoolSize)).
				SetConnectTimeout(config.ConnTimeout).
				SetMaxConnIdleTime(config.MaxConnIdleTime).
				SetMaxConnecting(uint64(config.MaxOpenConns)).
				ApplyURI(config.uri)
			client, err = mongo.Connect(context.Background(), cliOpt)
			if err != nil {
				return nil, err
			}
		}
	}
	return client, nil
}
