package main

import (
	"github.com/sirupsen/logrus"

	"github.com/thedeltaflyer/redirector/database"
	"github.com/thedeltaflyer/redirector/logging"
	"github.com/thedeltaflyer/redirector/server"

	flag "github.com/spf13/pflag"
)

var (
	Debug  = false          // Debug mode option
	Bind   = ":8080"        // Bind host and/or port
	DbPath = "./db/db.bolt" // Path to the BoltDB file
)

// main initializes the logger, enables debug mode if specified, initializes the database, and starts the HTTP server.
// It logs the startup and shutdown information and ensures proper closure of database resources on exit.
func main() {
	logger := logging.GetLogger()
	logger.Info("Starting redirector...")

	if Debug {
		logger.SetLevel(logrus.DebugLevel)
		logger.Info("Debug mode enabled")
	}

	database.InitDB(DbPath, true)
	defer database.CloseDB()

	logger.Infof("Starting redirector on %q", Bind)
	server.Run(Bind, Debug)

	logger.Info("Redirector stopped")
}

func init() {
	flag.BoolVarP(&Debug, "debug", "d", Debug, "Debug mode")
	flag.StringVarP(&Bind, "bind", "b", Bind, "Address/port to bind to")
	flag.StringVarP(&DbPath, "db", "s", DbPath, "Path to database file")
	flag.Parse()
}
