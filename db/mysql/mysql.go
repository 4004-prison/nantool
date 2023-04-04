package mysql

import (
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	defaultHost         = "127.0.0.1"
	defaultPort         = "3306"
	defaultMaxIdelConns = 10
	defaultMaxOpenConns = 20
)

const (
	dsnFormat = "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"
)

var db *gorm.DB

// Config mysql db configuration
type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"Port"`
	DB       string `yaml:"db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`

	MaxIdleConns int `yaml:"maxIdleConns"`
	MaxOpenConns int `yaml:"maxOpenConns"`

	dsn string
}

// New return a gorm db. it's a singleton.
func New(config *Config) (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}
	err := config.fill()
	if err != nil {
		return nil, err
	}
	db, err = gorm.Open(
		mysql.New(mysql.Config{
			DSN:                       config.dsn,
			DefaultStringSize:         256,
			DisableDatetimePrecision:  true,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			SkipInitializeWithVersion: false,
		}))
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)

	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = defaultMaxOpenConns
	}
	sqlDB.SetMaxOpenConns(defaultMaxOpenConns)
	return db, nil
}

func (c *Config) fill() error {
	if c.User == "" {
		return errors.New("the connetion user name cannot be empty")
	}
	if c.Password == "" {
		return errors.New("the connetion password cannot be empty")
	}
	if c.Host == "" {
		c.Host = defaultHost
	}
	if c.Port == "" {
		c.Port = defaultPort
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = defaultMaxIdelConns
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = defaultMaxOpenConns
	}
	c.dsn = fmt.Sprintf(dsnFormat, c.User, c.Password, c.Host, c.Port, c.DB)
	return nil
}
