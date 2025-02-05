package cli

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tecnologer/tempura/cmd/api/handler"
	"github.com/tecnologer/tempura/cmd/api/middleware"
	"github.com/tecnologer/tempura/cmd/api/router"
	"github.com/tecnologer/tempura/cmd/flags"
	"github.com/tecnologer/tempura/pkg/dao/db"
	"github.com/tecnologer/tempura/pkg/utils/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/negroni"
)

type CLI struct {
	*cli.App
	connection *db.Connection
}

func NewCLI(versionValue string) *CLI {
	newCLI := &CLI{}

	newCLI.setupApp(versionValue)

	return newCLI
}

func (c *CLI) setupApp(versionValue string) {
	c.App = &cli.App{
		Name:        "tempura-api",
		Version:     versionValue,
		Usage:       "Runs the API HTTP server for Tempura",
		Description: "",
		Action:      c.run,
		Before:      c.beforeRun,
		Flags: []cli.Flag{
			flags.Verbose(),
			flags.DBHost(),
			flags.DBPort(),
			flags.DBName(),
			flags.DBUsername(),
			flags.DBPassword(),
			flags.DBSSLMode(),
			flags.APIPort(),
			flags.APIWriteTimeout(),
			flags.APIReadTimeout(),
			flags.APIIdleTimeout(),
		},
		EnableBashCompletion: true,
	}
}

func (c *CLI) beforeRun(ctx *cli.Context) error {
	// Disable color globally.
	if ctx.Bool(flags.VerboseFlagName) {
		log.SetLevel(log.DebugLevel)
	}

	connection, err := c.createConnection(ctx)
	if err != nil {
		return fmt.Errorf("create connection: %w", err)
	}

	c.connection = connection

	return nil
}

func (c *CLI) run(ctx *cli.Context) error {
	log.Info("api started")

	var (
		apiHandler = handler.NewHandler(c.connection)
		apiRouter  = router.New(apiHandler)
		addr       = fmt.Sprintf(":%d", ctx.Int(flags.APIPortFlagName))
	)

	negroniMiddleware := negroni.New(middleware.NewMiddleware())
	negroniMiddleware.UseHandler(apiRouter)

	server := &http.Server{
		Addr:         addr,
		Handler:      negroniMiddleware,
		ReadTimeout:  time.Duration(ctx.Int(flags.APIReadTimeoutName)) * time.Second,
		WriteTimeout: time.Duration(ctx.Int(flags.APIWriteTimeoutName)) * time.Second,
		IdleTimeout:  time.Duration(ctx.Int(flags.APIIdleTimeoutName)) * time.Second,
	}

	log.Infof(
		"listening on %s with read time out %d secods, write time out %d seconds and idle time out %d seconds",
		addr,
		ctx.Int(flags.APIReadTimeoutName),
		ctx.Int(flags.APIWriteTimeoutName),
		ctx.Int(flags.APIIdleTimeoutName),
	)

	if err := server.ListenAndServe(); err != nil {
		log.Error(err.Error())
	}

	return nil
}

func (c *CLI) createConnection(ctx *cli.Context) (*db.Connection, error) {
	dbConfig := &db.Config{
		Host:     ctx.String(flags.DBHostFlagName),
		Port:     ctx.String(flags.DBPortFlagName),
		User:     ctx.String(flags.DBUsernameFlagName),
		Password: ctx.String(flags.DBPasswordFlagName),
		DBName:   ctx.String(flags.DBNameFlagName),
		SSLMode:  ctx.String(flags.DBSSLModeFlagName),
	}

	cnn, err := db.NewConnection(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("create new connection: %w", err)
	}

	return cnn, nil
}
