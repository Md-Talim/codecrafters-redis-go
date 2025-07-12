package config

import "flag"

type Config struct {
	Dir        string
	DBFilename string
	Port       string
	ReplicaOf  string
}

var instance *Config

func Load() *Config {
	if instance != nil {
		return instance
	}

	dir := flag.String("dir", "/tmp", "Directory for RDB file")
	dbfilename := flag.String("dbfilename", "dump.rdb", "RDB filename")
	port := flag.String("port", "6379", "Port to listen on")
	replicaOf := flag.String("replicaof", "", "Make this instance a replica of <host> <port>")

	flag.Parse()

	instance = &Config{
		Dir:        *dir,
		DBFilename: *dbfilename,
		Port:       *port,
		ReplicaOf:  *replicaOf,
	}

	return instance
}

func Get() *Config {
	if instance == nil {
		return Load()
	}
	return instance
}

func (c *Config) GetParameter(param string) (string, bool) {
	switch param {
	case "dir":
		return c.Dir, true
	case "dbfilename":
		return c.DBFilename, true
	case "port":
		return c.Port, true
	default:
		return "", false
	}
}

func (c *Config) IsReplica() bool {
	return c.ReplicaOf != ""
}
