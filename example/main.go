package main

import (
	"github.com/zephyyrr/plugins"
)

func main() {
	ph := plugins.NewHandler()
	ph.Handle(PrintPlugin{})
	ph.Handle(MessangerPlugin{})

	ph.ListenAndServe()
}
