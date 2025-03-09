package main

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

var registry map[string]cliCommand

func init() {
	{
		registry = map[string]cliCommand{
			"exit": {
				name:        "exit",
				description: "Exit the Pokedex",
				callback:    commandExit,
			},
			"help": {
				name:        "help",
				description: "Displays a help message",
				callback:    commandHelp,
			},
			"map": {
				name:        "map",
				description: "Displays name of 20 new locations each time",
				callback:    commandMap,
			},
			"mapb": {
				name:        "map",
				description: "Displays previous 20 locations each time",
				callback:    commandMapb,
			},
			"explore": {
				name:        "explore",
				description: "Explore Pokemons in a given area",
				callback:    commandExplore,
			},
			"catch": {
				name:        "catch",
				description: "Attempt to catch a Pikachu",
				callback:    commandCatch,
			},
			"inspect": {
				name:        "inspect",
				description: "Inspects details of Pokemon passed",
				callback:    commandInspect,
			},
			"pokedex": {
				name:        "pokedex",
				description: "Lists all the caught Pokemons",
				callback:    commandPokedex,
			},
		}
	}
}
