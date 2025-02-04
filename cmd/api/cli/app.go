package cli

import (
	"fmt"
	"net/http"

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

	log.Infof("listening on %s", addr)

	n := negroni.New(middleware.NewMiddleware())

	n.UseHandler(apiRouter)

	if err := http.ListenAndServe(addr, n); err != nil {
		log.Error(err.Error())
	}

	return nil
}

func (c *CLI) createConnection(ctx *cli.Context) (*db.Connection, error) {
	log.Infof("connecting to the DB at %s:%s", ctx.String(flags.DBHostFlagName), ctx.String(flags.DBPortFlagName))

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

	log.Infof("connection to the DB at %s:%s established", ctx.String(flags.DBHostFlagName), ctx.String(flags.DBPortFlagName))

	return cnn, nil
}
