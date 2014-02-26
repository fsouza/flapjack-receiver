package main

import (
    "github.com/giganteous/goflap"

    "github.com/hoisie/redis"    // connecting to flapjacks backend
    "github.com/globocom/config" // for parsing flapjack's yaml config

    "log"
    "os" // for environment
    "fmt"
    "flag"
)

const (
    DEFAULT_CONFIG_PATH = "/etc/flapjack/flapjack_config.yaml"
    DEFAULT_PERFDATA_LOCATION = "/var/cache/nagios3/event_stream.fifo"
    DEFAULT_CMDFILE_LOCATION = "/var/lib/nagios3/rw/nagios.cmd"
    DEFAULT_PIDFILE = "/var/run/flapjack/flapjack-nagios-receiver.pid"
)
var (
    FLAPJACK_ENV = "production"
    configfile, fifo_perf, fifo_cmd, pidfile, logfile string
    daemonize bool
)

func init() {
    flag.StringVar(&configfile, "config", DEFAULT_CONFIG_PATH,
        "PATH to the config file to use")
    flag.StringVar(&fifo_perf, "perfdata", DEFAULT_PERFDATA_LOCATION,
        "Path to the nagios perfdata named pipe")
    flag.StringVar(&fifo_cmd, "commandfile", DEFAULT_CMDFILE_LOCATION,
        "Path to the nagios commandfile named pipe")
    flag.BoolVar(&daemonize, "daemonize", false, "Daemonize")
    flag.StringVar(&pidfile, "pidfile", "", "Path to the pidfile to write to")
    flag.StringVar(&logfile, "logfile", "", "Path to the logfile to write to")
}

func main() {
    flag.Parse()

    if t := os.Getenv("FLAPJACK_ENV"); t != "" {
        FLAPJACK_ENV = t
    }

    fmt.Println("Going to parse ", configfile)
    config.ReadAndWatchConfigFile(configfile)
    fmt.Println("Some settings from config (env=", FLAPJACK_ENV, "):")

    host, _ := config.GetString(FLAPJACK_ENV + ":redis:host")
    db, _ := config.GetInt(FLAPJACK_ENV + ":redis:db")
    port, _ := config.GetInt(FLAPJACK_ENV + ":redis:port")
    //fmt.Println("redis port: ", config.GetInt(FLAPJACK_ENV + ":redis:port"))
    //fmt.Println("redis db: ", config.GetInt(FLAPJACK_ENV + ":redis:db"))

    address := fmt.Sprintf("%s:%d", host, port)
    fmt.Println("Going to connect to ", address)
    client := &redis.Client{Addr: address, Db: db}

    t := goflap.Event{
        Entity: "fubar.edu",
        Check: "FUBARCHECK",
        Type: "service",
        State: "warning",
        Summary: "Missed an option somewhere",
        Details: "We checked and we were really missing the option",
    }
    err := t.Add(client)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Added a bogus event to ", FLAPJACK_ENV)

}

