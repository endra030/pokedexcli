package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/nannabat/pokedexcli/internal/pokecache"
	"github.com/nannabat/pokedexcli/internal/utils"
)

const (
	URL             = "https://pokeapi.co/api/v2/"
	cacheExpiration = 5 * time.Minute
)

var (
	globalCache    *pokecache.Cache
	cacheOne       sync.Once
	caughtPokemons map[string]Pokemon
)

func GlobalCache() *pokecache.Cache {
	cacheOne.Do(func() {
		globalCache = pokecache.NewCache(cacheExpiration)
	})
	return globalCache
}

func CaughtPokemon() map[string]Pokemon {
	if caughtPokemons == nil {
		caughtPokemons = map[string]Pokemon{}
	}

	cacheOne.Do(func() {
		caughtPokemons = map[string]Pokemon{}
	})
	return caughtPokemons
}

func commandExit(c *config, p ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
func commandHelp(c *config, p ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for cmd, cmdDetails := range registry {
		hStmt := fmt.Sprintf("%s: %s", cmd, cmdDetails.description)
		fmt.Println(hStmt)

	}
	return nil
}
func commandMap(c *config, p ...string) error {
	cache := GlobalCache()
	fullUrl := URL + "location-area?offset=0&limit=20"
	var reader io.Reader
	locResp := LocationResponse{}
	if c.Next != "" {
		fullUrl = c.Next

	}
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return err
	}
	cacheData, inCache := cache.Get(fullUrl)
	if inCache {
		fmt.Println("Response from cache")
		reader = bytes.NewReader(cacheData)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(&locResp)
		if err != nil {
			return err
		}
	} else {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil
		}
		var buf bytes.Buffer
		reader = io.TeeReader(resp.Body, &buf)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(&locResp)
		if err != nil {
			return err
		}
		bBytes := buf.Bytes()
		//fmt.Printf("updating cache with key:%v\n", fullUrl)
		cache.Add(fullUrl, bBytes)

	}

	if locResp.Next != nil {
		c.Next = *locResp.Next
	} else {
		c.Next = ""
	}
	if locResp.Previous != nil {
		c.Previous = *locResp.Previous
	} else {
		c.Previous = ""
	}
	results := locResp.Results
	// fmt.Printf("c.Next:%v\n", c.Next)
	// fmt.Printf("c.Previous:%v\n", c.Previous)
	for _, r := range results {
		fmt.Println(r.Name)
	}
	return err

}

func commandMapb(c *config, p ...string) error {
	cache := GlobalCache()
	var reader io.Reader
	locResp := LocationResponse{}
	if c.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	fullUrl := c.Previous

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return err
	}
	cacheData, inCache := cache.Get(fullUrl)
	if inCache {
		//fmt.Println("Data is from cache")
		reader = bytes.NewReader(cacheData)
		reader = bytes.NewReader(cacheData)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(&locResp)
		if err != nil {
			return err
		}
	} else {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		reader = io.TeeReader(resp.Body, &buf)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(&locResp)
		if err != nil {
			return err
		}
		bBytes := buf.Bytes()
		cache.Add(fullUrl, bBytes)

	}

	if err != nil {
		return err
	}
	if locResp.Next != nil {
		c.Next = *locResp.Next
	} else {
		c.Next = ""
	}
	if locResp.Previous != nil {
		c.Previous = *locResp.Previous
	} else {
		c.Previous = ""
	}
	results := locResp.Results
	for _, r := range results {
		fmt.Println(r.Name)
	}
	return err

}

func commandExplore(c *config, p ...string) error {
	if len(p) <= 0 {
		return errors.New("no palce is passed to explore Pokemons")
	}
	fullUrl := URL + "location-area/" + p[0]
	//fmt.Printf("fullUrl:%s", fullUrl)
	var reader io.Reader
	cache := GlobalCache()
	locAreaResp := LocationAreaResponse{}
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return err
	}
	cacheData, inCache := cache.Get(fullUrl)
	if inCache {
		//fmt.Println("Data is from cache")
		reader = bytes.NewReader(cacheData)
		reader = bytes.NewReader(cacheData)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(&locAreaResp)
		if err != nil {
			return err
		}
	} else {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		reader = io.TeeReader(resp.Body, &buf)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(&locAreaResp)
		if err != nil {
			return err
		}
		bBytes := buf.Bytes()
		cache.Add(fullUrl, bBytes)

	}
	pokemonEncounters := locAreaResp.PokemonEncounters
	fmt.Println("Exploring pastoria-city-area...")
	pokemonsList := []string{}
	for _, encounter := range pokemonEncounters {
		pokemonName := encounter.Pokemon.Name
		pokemonsList = append(pokemonsList, pokemonName)
	}
	fmt.Println("Found Pokemon:")
	if len(pokemonsList) > 0 {
		for _, p := range pokemonsList {
			fmt.Printf("- %s\n", p)
		}
	}

	return err

}

func commandCatch(c *config, p ...string) error {
	if len(p) < 1 {
		return errors.New("no Pokemon name specified to catch")
	}
	var reader io.Reader
	cache := GlobalCache()
	cp := CaughtPokemon()
	fullUrl := URL + "pokemon/" + p[0]
	pokemonResp := PokemonResp{}
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return err
	}
	cacheData, inCache := cache.Get(fullUrl)
	if inCache {
		//fmt.Println("Data is from cache")
		reader = bytes.NewReader(cacheData)
		reader = bytes.NewReader(cacheData)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(&pokemonResp)
		if err != nil {
			return err
		}
	} else {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		reader = io.TeeReader(resp.Body, &buf)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(&pokemonResp)
		if err != nil {
			return err
		}
		bBytes := buf.Bytes()
		cache.Add(fullUrl, bBytes)

	}
	catchChance := utils.CatchPokemon(pokemonResp.BaseExperience)
	fmt.Printf("Throwing a Pokeball at %s...	\n", p[0])
	// fmt.Printf("catchPokemonz:%v\n", catchChance)
	// fmt.Printf("Base experience:%v\n", pokemonResp.BaseExperience)
	if catchChance {

		fmt.Printf("%s was caught!\n", p[0])
		fmt.Println("You may now inspect it with the inspect command.")
		cp[p[0]] = Pokemon{
			Name: p[0],
			URL:  fullUrl,
		}
		//fmt.Printf("caughtPokemons:%v\n", cp)
	} else {
		fmt.Printf("%s escaped!\n", p[0])
		//fmt.Printf("caughtPokemons:%v\n", cp)

	}

	return nil
}

func commandInspect(c *config, p ...string) error {
	pokemonToInspect := p[0]
	cp := CaughtPokemon()
	cache := GlobalCache()
	cacheData, inCache := cache.Get(cp[pokemonToInspect].URL)
	pokemonResp := PokemonResp{}
	req, err := http.NewRequest("GET", cp[pokemonToInspect].URL, nil)
	if err != nil {
		return err
	}
	var reader io.Reader
	_, ok := cp[pokemonToInspect]
	if ok {
		if inCache {
			//fmt.Println("Data is from cache")
			reader = bytes.NewReader(cacheData)
			decoder := json.NewDecoder(reader)
			err := decoder.Decode(&pokemonResp)
			if err != nil {
				return err
			}
		} else {
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			var buf bytes.Buffer
			reader = io.TeeReader(resp.Body, &buf)
			decoder := json.NewDecoder(reader)
			err = decoder.Decode(&pokemonResp)
			if err != nil {
				return err
			}
			bBytes := buf.Bytes()
			cache.Add(cp[pokemonToInspect].URL, bBytes)

		}
		fmt.Printf("Name: %s\n", pokemonResp.Name)
		fmt.Printf("Height: %d\n", pokemonResp.Height)
		fmt.Printf("Weight: %d\n", pokemonResp.Weight)
		fmt.Printf("Stats:\n")
		for _, statItem := range pokemonResp.Stats {
			fmt.Printf("  -%s: %d\n", statItem.Stat.Name, statItem.Base_stat)
		}
		fmt.Println("Types:")
		for _, t := range pokemonResp.Types {
			fmt.Printf(" -%s\n", t.Type.Name)

		}

	} else {
		fmt.Println("you have not caught that pokemon")
	}

	return nil
}

func commandPokedex(c *config, p ...string) error {
	cp := CaughtPokemon()
	if len(cp) > 0 {
		fmt.Println("Your Pokedex:")
		for pName, _ := range cp {
			fmt.Printf("- %s\n", pName)

		}

	} else {
		fmt.Println("No Pokemons caught")
	}
	return nil

}
