package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Vertex struct {
	Key   int
	Edges map[int]float64
}

type Graph struct {
	Vertices map[int]*Vertex
}

func (g *Graph) Dijkstra(startKey int) (map[int]float64, map[int]int, error) {
	distances := make(map[int]float64)
	previous := make(map[int]int)
	for key := range g.Vertices {
		distances[key] = math.MaxFloat64
	}
	distances[startKey] = 0

	visited := make(map[int]bool)
	for len(visited) < len(g.Vertices) {
		var current *Vertex
		var minDist float64
		for _, vertex := range g.Vertices {
			if !visited[vertex.Key] && (current == nil || distances[vertex.Key] < minDist) {
				current = vertex
				minDist = distances[vertex.Key]
			}
		}

		visited[current.Key] = true
		for neighborKey, cost := range current.Edges {
			if !visited[neighborKey] {
				newDist := distances[current.Key] + cost
				if newDist < distances[neighborKey] {
					distances[neighborKey] = newDist
					previous[neighborKey] = current.Key
				}
			}
		}
	}

	return distances, previous, nil
}

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
	Cost   string
}

func main() {
	if len(os.Args) != 2 {
		panic("provide 1 arg, filename")
	}

	filename := os.Args[1]
	cities, distances, err := LoadCities(filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	graph := &Graph{Vertices: make(map[int]*Vertex)}
	for key, _ := range distances {
		vertex := &Vertex{Key: key, Edges: make(map[int]float64)}
		graph.Vertices[key] = vertex
	}

	for from, edges := range distances {
		fromVertex := graph.Vertices[from]
		for to, cost := range edges {
			fromVertex.Edges[to] = cost
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		state := State{
			Cities: cities,
			Path:   "",
		}
		Render(w, "index", state)
	})
	mux.HandleFunc("/distance", func(w http.ResponseWriter, r *http.Request) {
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

		var path []string
		resultDijkstra, previous, err := graph.Dijkstra(city1)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Shortest distances from vertex", city1)
			for key, dist := range resultDijkstra {
				fmt.Printf("Vertex %d: %f\n", key, dist)
			}

			current := city2
			for current != city1 {
				if cityName, exists := cities[current]; exists {
					path = append([]string{cityName}, path...)
				}
				current = previous[current]
			}
			if cityName, exists := cities[city1]; exists {
				path = append([]string{cityName}, path...)
			}

			fmt.Println("Shortest path from city", city1, "to city", city2, ":", strings.Join(path, " -> "))
		}

		state := State{
			Cities: cities,
			Path:   strings.Join(path, " -> "),
			Cost:   fmt.Sprintf("%.2f", resultDijkstra[city2]),
		}
		Render(w, "index", state)
	})

	http.ListenAndServe(":8080", mux)
}
