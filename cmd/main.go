package main

import (
	"social/internal/app"
	"social/internal/config"
)

func main() {
	cnf := config.New()

	app.Run(cnf)
}
