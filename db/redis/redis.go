package redis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	defaultPoolSize = 20
)

var client *redis.ClusterClient
var mu sync.Mutex

// Config config
type Config struct {
	Addrs           []string `yaml:"address"`
	Username        string   `yaml:"username"`
	Password        string   `yaml:"password"`
	MaxRetries      int
	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration

	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PoolSize           int           `yaml:"poolSize"`
	MinIdleConns       int           `yaml:"minIdleConns"`
	MaxConnAge         time.Duration `yaml:"minIdleAge"`
	PoolTimeout        time.Duration `yaml:"poolTimeout"`
	IdleTimeout        time.Duration `yaml:"idleTimeout"`
	IdleCheckFrequency time.Duration `yaml:"idleCheckFrequency"`

	TLS *TLS
}

func (c *Config) fill() error {
	if len(c.Addrs) == 0 {
		return errors.New("the redis connection address cannot be empty")
	}
	if c.PoolSize == 0 {
		c.PoolSize = defaultPoolSize
	}
	if c.MinIdleConns == 0 {
		c.MinIdleConns = c.PoolSize / 5
	}
	return nil
}

// TLS tls
type TLS struct {
	ClientCertFile string
	clientKeyFile  string
	CACertFile     string
}

// NewClient new redis cluster client
func NewClient(conf *Config) (*redis.ClusterClient, error) {
	if client == nil {
		mu.Lock()
		defer mu.Unlock()
		if client == nil {
			if err := conf.fill(); err != nil {
				return nil, err
			}
			var (
				tlsConfig *tls.Config
				err       error
			)
			if conf.TLS != nil {
				tlsConfig, err = NewTLSConfig(conf.TLS.ClientCertFile, conf.TLS.clientKeyFile, conf.TLS.CACertFile)
			}
			if err != nil {
				return nil, err
			}
			client = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:              conf.Addrs,
				Username:           conf.Username,
				Password:           conf.Password,
				MaxRetries:         conf.MaxRetries,
				MinRetryBackoff:    conf.MinRetryBackoff,
				MaxRetryBackoff:    conf.MaxRetryBackoff,
				DialTimeout:        conf.DialTimeout,
				ReadTimeout:        conf.ReadTimeout,
				WriteTimeout:       conf.WriteTimeout,
				PoolSize:           conf.PoolSize,
				MinIdleConns:       conf.MinIdleConns,
				MaxConnAge:         conf.MaxConnAge,
				PoolTimeout:        conf.PoolTimeout,
				IdleTimeout:        conf.IdleTimeout,
				IdleCheckFrequency: conf.IdleCheckFrequency,
				TLSConfig:          tlsConfig,
			})

			if err := client.Ping(context.Background()).Err(); err != nil {
				return nil, err
			}
		}
	}
	return client, nil
}

// NewTLSConfig generates a TLS configuration used to authenticate on server with
// certificates.
// Parameters are the three pem files path we need to authenticate: client cert, client key and CA cert.
func NewTLSConfig(clientCertFile, clientKeyFile, caCertFile string) (*tls.Config, error) {
	tlsConfig := tls.Config{}

	// Load client cert
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return &tlsConfig, err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return &tlsConfig, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig.RootCAs = caCertPool

	tlsConfig.BuildNameToCertificate()
	return &tlsConfig, err
}
