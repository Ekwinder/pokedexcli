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
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode > 299 {
		fmt.Printf("Response from %s failed with status code %d\n", apiName, res.StatusCode)
	}

	if err != nil {
		fmt.Printf("Response from %s failed with error %s\n", apiName, err)
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
	BaseExperience int `json:"base_experience"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
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
	}
	prob := rand.Intn(poke.BaseExperience)
	pp := float32(poke.BaseExperience)

	if prob < poke.BaseExperience-int(0.5*pp) {
		fmt.Printf("%s escaped!\n", name)
	} else {
		fmt.Printf("%s was caught!\n", name)
		Pokedex[name] = poke
	}
	fmt.Println(Pokedex)
}

func Inspect(name string){
	fmt.Println("inspecting")
}