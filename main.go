package main

import (
	"fmt"
	"github.com/coreos/bbolt"
	log "github.com/go-pkgz/lgr"
	"github.com/jessevdk/go-flags"
	"os"
	"os/signal"
	"runtime"
	"shopaholic/cmd"
	"shopaholic/store/engine"
	"shopaholic/store/service"
	"syscall"
	"time"
)

// Opts with all cli commands and flags
type Opts struct {
	cmd.UserCreateCommand        `command:"user:create"`
	cmd.UserListCommand          `command:"user:list"`
	cmd.TransactionCreateCommand `command:"transaction:create"`
	cmd.TransactionListCommand   `command:"transaction:list"`

	cmd.BotPollerCommand `command:"bot:start"`

	Currency   string `long:"currency" env:"CURRENCY" default:"usd" description:"money currency"`
	DBFilename string `long:"db_filename" env:"DBFILENAME" default:"shopaholic.db" description:"database filename"`
	RedisHost  string `long:"redis_host" env:"REDIS_HOST" default:"localhost:6379" description:"redis host"`
	RedisDB    int    `long:"redis_db" env:"REDIS_DB" default:"0" description:"redis db num"`
	BotToken   string `long:"bot_token" env:"BOT_TOKEN" description:"token of the bot"`

	Dbg bool `long:"dbg" env:"DEBUG" description:"debug mode"`
}

var revision = "0.2.0"

func main() {
	fmt.Printf("shopaholic %s\n", revision)

	var opts Opts
	p := flags.NewParser(&opts, flags.Default)

	p.CommandHandler = func(command flags.Commander, args []string) error {
		setupLog(opts.Dbg)
		// commands implements CommonOptionsCommander to allow passing set of extra options defined for all commands
		c := command.(cmd.Commander)
		c.SetCommon(cmd.CommonOpts{
			Currency: opts.Currency,
			Store:    *initRedisStore(opts.RedisHost, opts.RedisDB),
			BotToken: opts.BotToken,
		})
		err := c.Execute(args)
		if err != nil {
			log.Printf("[ERROR] failed with %+v", err)
		}
		return err
	}

	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}

func initBoltStore(dbFilename string) *service.DataStore {
	b, err := engine.NewBoltDB(bbolt.Options{Timeout: 30 * time.Second}, dbFilename)
	if err != nil {
		log.Printf("[ERROR] can not initiate DB %s", dbFilename)
		os.Exit(0)
	}

	return &service.DataStore{Interface: b}
}

func initRedisStore(host string, db int) *service.DataStore {
	b, err := engine.NewRedisClient(host, db)
	if err != nil {
		log.Printf("[ERROR] can not initiate redis at %s, %d", host, db)
		os.Exit(0)
	}

	return &service.DataStore{Interface: b}
}

func setupLog(dbg bool) {
	if dbg {
		log.Setup(log.Debug, log.CallerFile, log.Msec, log.LevelBraces)
		return
	}
	log.Setup(log.Msec, log.LevelBraces, log.CallerPkg, log.CallerIgnore("logger"))
}

// getDump reads runtime stack and returns as a string
func getDump() string {
	maxSize := 5 * 1024 * 1024
	stacktrace := make([]byte, maxSize)
	length := runtime.Stack(stacktrace, true)
	if length > maxSize {
		length = maxSize
	}
	return string(stacktrace[:length])
}

func init() {
	sigChan := make(chan os.Signal)
	go func() {
		for range sigChan {
			log.Printf("[INFO] SIGQUIT detected, dump:\n%s", getDump())
		}
	}()
	signal.Notify(sigChan, syscall.SIGQUIT)
}
