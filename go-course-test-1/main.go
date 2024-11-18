package main

// dijkstra.go and queue.go are adjust implementation of this lib
// https://github.com/albertorestifo/dijkstra/tree/aba76f725f72b3086bd6957f94692feb8a5f9bb9

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func LoadCities(filename string) (map[int]string, map[int]map[int]float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return map[int]string{}, map[int]map[int]float64{}, fmt.Errorf("failed to open file: %s", err.Error())
	}
	defer file.Close()

	cities := make(map[int]string)
	distances := make(map[int]map[int]float64)

	scanner := bufio.NewScanner(file)
	readingCities := true
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			readingCities = false
			continue
		}

		if readingCities {
			parts := strings.Split(line, ",")
			if len(parts) != 2 {
				continue
			}
			id, err := strconv.Atoi(parts[0])
			if err != nil {
				continue
			}
			cities[id] = parts[1]
		} else {
			parts := strings.Split(line, ",")
			if len(parts) != 3 {
				continue
			}
			city1ID, err := strconv.Atoi(parts[0])
			if err != nil {
				continue
			}
			city2ID, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
			}
			distance, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				continue
			}

			if distances[city1ID] == nil {
				distances[city1ID] = make(map[int]float64)
			}
			distances[city1ID][city2ID] = distance
		}
	}

	if err := scanner.Err(); err != nil {
		return map[int]string{}, map[int]map[int]float64{}, fmt.Errorf("error reading file: %s", err.Error())
	}

	return cities, distances, nil
}

func Render(w io.Writer, name string, data interface{}) error {
	return template.Must(template.ParseGlob("*.html")).ExecuteTemplate(w, name, data)
}

type State struct {
	Cities map[int]string
	Path   string
	Cost   float64
}

func main() {
	if len(os.Args) != 2 {
		panic("provide 1 arg, filename")
	}

	filename := os.Args[1]
	cities, distances, err := LoadCities(filename)
	if err != nil {
		fmt.Println(err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		state := State{
			Cities: cities,
			Path:   "",
			Cost:   0.0,
		}
		Render(w, "index", state)
	})
	mux.HandleFunc("POST /distance", func(w http.ResponseWriter, r *http.Request) {
		city1Str := r.FormValue("city1")
		city1, err := strconv.Atoi(city1Str)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		city2Str := r.FormValue("city2")
		city2, err := strconv.Atoi(city2Str)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		g := Graph(distances)
		path, cost, err := g.Path(city1, city2)
		if err != nil {
			fmt.Println(err.Error())
		}

		var result []string
		for _, key := range path {
			if val, exists := cities[key]; exists {
				result = append(result, val)
			}
		}
		state := State{
			Cities: cities,
			Path:   strings.Join(result, " -> "),
			Cost:   cost,
		}
		Render(w, "index", state)
	})

	http.ListenAndServe(":8080", mux)

	fmt.Println("Cities:")
	fmt.Println(cities)
	fmt.Println("Distances:")
	fmt.Println(distances)
}
