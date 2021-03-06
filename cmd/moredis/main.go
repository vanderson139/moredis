package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vanderson139/moredis/moredis"
	"github.com/vanderson139/moredis/logger"
)

// Default database connection parameters
const (
	DefaultMongoURL = "localhost:27017"
	DefaultRedisURL = "localhost:6379"
)

var (
	redisURL       string
	mongoURL       string
	params         moredis.Params
	configFilePath string
	buildMongo     bool
	start          int
	end            int
)

func init() {
	const (
		defaultFilePath = "./config.yml"
	)
	// Usage strings in PrintUsage
	flag.StringVar(&redisURL, "redis_url", "", "")
	flag.StringVar(&redisURL, "r", "", "")
	flag.StringVar(&mongoURL, "mongo_url", "", "")
	flag.StringVar(&mongoURL, "m", "", "")
	flag.Var(&params, "params", "")
	flag.Var(&params, "p", "")
	flag.StringVar(&configFilePath, "conf_file", defaultFilePath, "")
	flag.StringVar(&configFilePath, "f", defaultFilePath, "")
	flag.BoolVar(&buildMongo, "build_mongo", false, "")
	flag.BoolVar(&buildMongo, "b", false, "")
	flag.IntVar(&start, "start", 0, "")
	flag.IntVar(&start, "s", 0, "")
	flag.IntVar(&end, "end", 0, "")
	flag.IntVar(&end, "e", 0, "")
}

func main() {
	flag.Usage = PrintUsage
	flag.Parse()

	// grab connection from env or default if not in flags
	mongoURL = FlagEnvOrDefault(mongoURL, "MONGO_URL", DefaultMongoURL)
	redisURL = FlagEnvOrDefault(redisURL, "REDIS_URL", DefaultRedisURL)

	conf, err := moredis.LoadConfig(configFilePath)
    if err != nil {
        logger.Error("Error loading config.", err)
        os.Exit(1)
    }

	if buildMongo == false {
        if err := moredis.BuildCache(conf, params, redisURL, mongoURL); err != nil {
            fmt.Fprint(os.Stderr, err)
            os.Exit(1)
        }
	} else {
	    if err := moredis.BuildMongo(conf, redisURL, mongoURL, start, end); err != nil {
            fmt.Fprint(os.Stderr, err)
            os.Exit(1)
        }
	}
}

// PrintUsage is used to replace flag.Usage, which is pretty terrible.
func PrintUsage() {
	var usage = `Usage of ./moredis:
  -m, -mongo_url    MongoDB URL, can also be set via the MONGO_URL environment variable
  -p, -params       JSON object with params used for substitution into queries and collection names in config.yml
  -r, -redis_url    Redis URL, can also be set via the REDIS_URL environment variable
  -f, -conf_file    Config file, defaults to ./config.yml
  -b, -build_mongo  Copy Redis keys to MongoDB
  -s, -start        Redis SCAN range start
  -e, -end          Redis SCAN range end
  -h, -help         Print this usage message
`
	fmt.Fprint(os.Stderr, usage)
}

// FlagEnvOrDefault takes in the value returned from flag parsing, and environment variable to check, and a default value.
// If the flag was not set, it tries to retrieve the value from environment variables.  If it is not found there it returns
// the default value.
func FlagEnvOrDefault(flagVal, envVar, defaultValue string) string {
	if flagVal != "" {
		return flagVal
	}
	if fromEnv := os.Getenv(envVar); fromEnv != "" {
		return fromEnv
	}
	return defaultValue
}
