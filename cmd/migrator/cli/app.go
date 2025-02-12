package cli

import (
	"fmt"

	"github.com/tecnologer/tempura/cmd/flags"
	"github.com/tecnologer/tempura/pkg/dao/db"
	"github.com/tecnologer/tempura/pkg/utils/log"
	"github.com/urfave/cli/v2"
)

type CLI struct {
	*cli.App
	migrator *Migrator
}

func NewCLI(versionValue string) *CLI {
	newCLI := &CLI{
		migrator: NewMigrator(),
	}

	newCLI.setupApp(versionValue)

	return newCLI
}

func (c *CLI) setupApp(versionValue string) {
	c.App = &cli.App{
		Name:        "migrator",
		Version:     versionValue,
		Usage:       "Migrates the models into the DB",
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
			flags.SkipModels(),
			flags.SkipScripts(),
		},
		EnableBashCompletion: true,
	}
}

func (c *CLI) beforeRun(ctx *cli.Context) error {
	// Disable color globally.
	if ctx.Bool(flags.VerboseFlagName) {
		log.SetLevel(log.DebugLevel)
	}

	if ctx.Bool(flags.ImporterSkipModels) && ctx.Bool(flags.ImporterSkipScripts) {
		return fmt.Errorf("cannot skip both models and scripts")
	}

	return nil
}

func (c *CLI) run(ctx *cli.Context) error {
	log.Info("migrator started")

	dbCnn, err := c.createConnection(ctx)
	if err != nil {
		return fmt.Errorf("create db connection: %w", err)
	}

	defer func() {
		if err != nil {
			if rErr := dbCnn.Rollback(); rErr != nil {
				log.Warn("transaction rollback failed")
			}
		}
	}()

	if !ctx.Bool(flags.ImporterSkipModels) {
		err = c.migrateModels(dbCnn)
		if err != nil {
			return fmt.Errorf("migrate models: %w", err)
		}
	}

	if !ctx.Bool(flags.ImporterSkipScripts) {
		err = c.migrateScripts(dbCnn)
		if err != nil {
			return fmt.Errorf("migrate scripts: %w", err)
		}
	}

	log.Info("committing transaction for scripts")

	if rErr := dbCnn.Commit(); rErr != nil {
		log.Warn("transaction commit failed")
		return fmt.Errorf("commit: %w", rErr)
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
	log.Info("beginning transaction")

	err = cnn.BeginTransaction()
	if err != nil {
		return nil, fmt.Errorf("begin transaction")
	}

	log.Infof("transaction started")

	return cnn, nil
}

func (c *CLI) migrateModels(dbCnn *db.Connection) error {
	log.Info("migrating models")

	err := c.migrator.Run(dbCnn.DB)
	if err != nil {
		return fmt.Errorf("run migrator: %w", err)
	}

	log.Info("models migrated")

	return nil
}

func (c *CLI) migrateScripts(dbCnn *db.Connection) error {
	log.Info("migrating scripts")

	err := c.migrator.RunScripts(dbCnn.DB)
	if err != nil {
		log.Warn("transaction rollback")

		if rErr := dbCnn.Rollback(); rErr != nil {
			log.Warn("transaction rollback failed")

			err = fmt.Errorf("rollback: %w", rErr)
		}

		return fmt.Errorf("run scripts: %w", err)
	}

	log.Info("scripts migrated")

	return nil
}
