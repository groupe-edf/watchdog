// Schemes: https
// Host: watchdog.org
// BasePath: /api/v1
// Version: 1.0.0
// Contact: Habib MAALEM<habib.maalem@watchdog.org> https://www.watchdog.org
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Security:
// - bearer
// SecurityDefinitions:
// bearer:
//   type: apiKey
//   name: Authorization
//   in: header
//
// swagger:meta
// go:generate swagger generate spec
package main

import (
	"context"

	"github.com/groupe-edf/watchdog/cmd/watchdog-server/commands"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	commands.Execute(ctx)
}
