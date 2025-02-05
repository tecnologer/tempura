package flags

import (
	"github.com/tecnologer/tempura/pkg/contants/envvarname"
	"github.com/urfave/cli/v2"
)

const (
	VerboseFlagName     = "verbose"
	DBHostFlagName      = "db-host"
	DBPortFlagName      = "db-port"
	DBUsernameFlagName  = "db-username"
	DBPasswordFlagName  = "db-password"
	DBNameFlagName      = "db-name"
	DBSSLModeFlagName   = "db-ssl-mode"
	APIPortFlagName     = "api-port"
	APIReadTimeoutName  = "api-read-timeout"
	APIWriteTimeoutName = "api-write-timeout"
	APIIdleTimeoutName  = "api-idle-timeout"
)

func Verbose() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  VerboseFlagName,
		Usage: "enable verbose output.",
	}
}

func DBHost() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     DBHostFlagName,
		Aliases:  []string{"H"},
		Usage:    "database host.",
		Required: true,
		EnvVars:  []string{envvarname.DBHost},
	}
}

func DBPort() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     DBPortFlagName,
		Aliases:  []string{"p"},
		Usage:    "database port.",
		Required: true,
		EnvVars:  []string{envvarname.DBPort},
	}
}

func DBUsername() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     DBUsernameFlagName,
		Aliases:  []string{"u"},
		Usage:    "database username.",
		Required: true,
		EnvVars:  []string{envvarname.DBUsername},
	}
}

func DBPassword() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     DBPasswordFlagName,
		Aliases:  []string{"P"},
		Usage:    "database password.",
		Required: true,
	}
}

func DBName() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     DBNameFlagName,
		Aliases:  []string{"d"},
		Usage:    "database name.",
		Required: true,
		EnvVars:  []string{envvarname.DBName},
	}
}

func DBSSLMode() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    DBSSLModeFlagName,
		Aliases: []string{"ssl-mode"},
		Usage:   "database ssl mode.",
		Value:   "disable",
		EnvVars: []string{envvarname.DBSSLMode},
	}
}

func APIPort() *cli.IntFlag {
	return &cli.IntFlag{
		Name:    APIPortFlagName,
		Aliases: []string{"a"},
		Usage:   "port to listen on",
		Value:   8080,
		EnvVars: []string{envvarname.APIPort},
	}
}

func APIReadTimeout() *cli.IntFlag {
	return &cli.IntFlag{
		Name:    APIReadTimeoutName,
		Aliases: []string{"rt"},
		Usage:   "read timeout for the API server",
		Value:   30,
	}
}

func APIWriteTimeout() *cli.IntFlag {
	return &cli.IntFlag{
		Name:    APIWriteTimeoutName,
		Aliases: []string{"wt"},
		Usage:   "write timeout for the API server",
		Value:   15,
	}
}

func APIIdleTimeout() *cli.IntFlag {
	return &cli.IntFlag{
		Name:    APIIdleTimeoutName,
		Aliases: []string{"it"},
		Usage:   "idle timeout for the API server",
		Value:   60,
	}
}
