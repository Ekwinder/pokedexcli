package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
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

var url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
var prev string

var cache = pokecache.Cache{
	CacheEntry: map[string]pokecache.Entry{},
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

		res, err := http.Get(url)
		if err != nil {
			fmt.Printf("Fetch Map API failed with error %s, please try later", err)
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()

		if res.StatusCode > 299 {
			fmt.Printf("Response from Map API failed with status code %d\n", res.StatusCode)
		}

		if err != nil {
			fmt.Printf("Response from Map API failed with error %s\n", err)
		}

		err = json.Unmarshal(body, &maps)
		if err != nil {
			fmt.Printf("Response parsing failed with error %s\n", err)
		}

		// update the URL to call the next page
		cache.Add(url, body)

	}
	url = maps.Next

	//
	prev = maps.Previous
	for _, v := range maps.Results {
		fmt.Println(v.Name)
	}
}
