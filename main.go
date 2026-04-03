package main

import (
	"log"
	"os"

	"github.com/brentjolicoeur/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	s := &state{
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}
	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}
	err = cmds.run(s, cmd)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
