package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/Ekwinder/pokedexcli/internal/pokecache"
)

type MapData struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type PokeExplore struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

var baseUrl = "https://pokeapi.co/api/v2/"

var prev string

var cache = pokecache.Cache{
	CacheEntry: map[string]pokecache.Entry{},
}
var url = baseUrl + "location-area/?offset=0&limit=20"

func getResponse(apiName string, url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Fetch %s failed with error %s, please try later", apiName, err)
		return nil
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode > 299 {
		fmt.Printf("Response from %s failed with status code %d\n", apiName, res.StatusCode)
		return nil
	}

	if err != nil {
		fmt.Printf("Response from %s failed with error %s\n", apiName, err)
		return nil
	}

	return body
}

func GetMap(isPrev bool) {

	if isPrev && len(prev) == 0 {
		fmt.Println("Error: Already on first page")
		return
	} else if isPrev {
		url = prev
	}
	maps := MapData{}
	// try a cache hit before hitting the URL??

	cachedResponse, ok := cache.Get(url)
	if ok {
		err := json.Unmarshal(cachedResponse, &maps)
		if err != nil {
			fmt.Printf("Response parsing failed with error %s\n", err)
		}

	} else {

		body := getResponse("Map API", url)
		err := json.Unmarshal(body, &maps)
		if err != nil {
			fmt.Printf("Response parsing failed with error %s\n", err)
		}
		cache.Add(url, body)

	}
	url = maps.Next

	//
	prev = maps.Previous
	for _, v := range maps.Results {
		fmt.Println(v.Name)
	}
}

func Explore(area string) {
	// See Pokemon in an area
	areaUrl := baseUrl + "location-area/" + area
	exploreData := PokeExplore{}

	//check if the key is in cache
	cachedResponse, ok := cache.Get(areaUrl)

	// if it is then
	if ok {
		err := json.Unmarshal(cachedResponse, &exploreData)
		if err != nil {
			fmt.Printf("Response parsing failed with error %s\n", err)
		}

	} else {
		body := getResponse("Explore API", areaUrl)
		err := json.Unmarshal(body, &exploreData)
		if err != nil {
			fmt.Printf("Response parsing failed with error %s\n", err)
		}
		cache.Add(areaUrl, body)
	}
	fmt.Println("Found Pokemon:")
	for _, v := range exploreData.PokemonEncounters {
		fmt.Printf(" - %s\n", v.Pokemon.Name)
	}

}

type Pokemon struct {
	BaseExperience int    `json:"base_experience"`
	Name           string `json:"name"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Height int `json:"height"`
	Weight int `json:"weight"`
	Types  []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

var Pokedex = map[string]Pokemon{}

func Catch(name string) {
	// Catch a pokemon with chance
	catchUrl := baseUrl + "pokemon/" + name
	poke := Pokemon{}
	body := getResponse("Catch", catchUrl)

	err := json.Unmarshal(body, &poke)
	if err != nil {
		fmt.Printf("Response parsing failed with error %s\n", err)
		return
	}
	prob := rand.Intn(poke.BaseExperience)
	pp := float32(poke.BaseExperience)
	fmt.Printf("Throwing a Pokeball at %s...", name)
	if prob < poke.BaseExperience-int(0.5*pp) {
		fmt.Printf("%s escaped!\n", name)
	} else {
		fmt.Printf("%s was caught!\nYou may now inspect it with the inspect command.\n", name)
		Pokedex[name] = poke
	}
}

func Inspect(name string) {
	poke, ok := Pokedex[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return
	}
	fmt.Printf("Name: %s\n", poke.Name)
	fmt.Printf("Height: %d\n", poke.Height)
	fmt.Printf("Weight: %d\n", poke.Weight)
	fmt.Println("Stats:")
	for _, stat := range poke.Stats {

		fmt.Printf(" - %s %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, tp := range poke.Types {
		fmt.Printf(" - %s\n", tp.Type.Name)
	}

}

func PokedexRecord() {
	fmt.Println("Your Pokedex:")
	for k, _ := range Pokedex {
		fmt.Printf(" - %s\n", k)

	}
}
