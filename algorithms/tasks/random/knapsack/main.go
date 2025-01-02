package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	catalogue       []Item
	maxWeight       int
	maxNumber       int
	populationSize  = 100
	mutationRate    = 0.01
	generationLimit = 20000
)

type Item struct {
	Weight int
	Value  int
}

type Individual struct {
	Genes   []bool
	Fitness int
}

func (ind *Individual) calculateFitness() {
	totalWeight := 0
	totalValue := 0
	for i, gene := range ind.Genes {
		if gene {
			totalWeight += catalogue[i].Weight
			totalValue += catalogue[i].Value
		}
	}
	if totalWeight > maxWeight {
		totalValue = 0
	}
	ind.Fitness = totalValue
}

func generateIndividual() Individual {
	genes := make([]bool, len(catalogue))
	numItemsToPick := rand.Intn(maxNumber + 1)

	for i := 0; i < numItemsToPick; i++ {
		genes[i] = true
	}

	rand.Shuffle(len(genes), func(i, j int) {
		genes[i], genes[j] = genes[j], genes[i]
	})

	individual := Individual{Genes: genes}
	individual.calculateFitness()
	return individual
}

func generateInitPopulation() []Individual {
	population := make([]Individual, populationSize)
	for i := range population {
		population[i] = generateIndividual()
	}
	return population
}

func selection(population []Individual) []Individual {
	totalFitness := 0
	for _, ind := range population {
		totalFitness += ind.Fitness
	}

	if totalFitness == 0 {
		return generateInitPopulation()
	}

	selected := make([]Individual, 0, populationSize/2)
	for len(selected) < populationSize/2 {
		threshold := rand.Intn(totalFitness)
		sum := 0
		for _, ind := range population {
			sum += ind.Fitness
			if sum > threshold {
				selected = append(selected, ind)
				break
			}
		}
	}
	return selected
}

func crossover(parent1, parent2 Individual) Individual {
	point := rand.Intn(len(parent1.Genes))
	genes := append([]bool{}, parent1.Genes[:point]...)
	genes = append(genes, parent2.Genes[point:]...)
	child := Individual{Genes: genes}
	child.calculateFitness()
	return child
}

func reproduce(parents []Individual) []Individual {
	children := make([]Individual, len(parents))
	for i := 0; i < len(parents); i += 2 {
		if i+1 < len(parents) {
			children[i] = crossover(parents[i], parents[i+1])
			children[i+1] = crossover(parents[i+1], parents[i])
		} else {
			children[i] = parents[i]
		}
	}
	return children
}

func mutate(children []Individual) {
	for i := range children {
		for j := range children[i].Genes {
			if rand.Float64() < mutationRate {
				children[i].Genes[j] = !children[i].Genes[j]
			}
		}

		trueCount := 0
		for _, gene := range children[i].Genes {
			if gene {
				trueCount++
			}
		}

		for trueCount > maxNumber {
			for j := range children[i].Genes {
				if children[i].Genes[j] {
					children[i].Genes[j] = false
					trueCount--
				}
				if trueCount <= maxNumber {
					break
				}
			}
		}

		children[i].calculateFitness()
	}
}

func evolve(population []Individual) []Individual {
	parents := selection(population)
	children := reproduce(parents)
	mutate(children)
	newPopulation := append(population[:populationSize/4], children...)
	sort.Slice(newPopulation, func(i, j int) bool {
		return newPopulation[i].Fitness > newPopulation[j].Fitness
	})
	if len(newPopulation) > populationSize {
		newPopulation = newPopulation[:populationSize]
	}
	return newPopulation
}

func main() {
	catalogue = make([]Item, 0)
	file, err := os.Open("long.txt")
	if err != nil {
		fmt.Println("error opening file")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	isFirst := true
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 2 {
			weight, errWeight := strconv.Atoi(parts[0])
			value, errValue := strconv.Atoi(parts[1])
			if errWeight == nil && errValue == nil {
				if isFirst {
					isFirst = false
					maxWeight = weight
					maxNumber = value
					continue
				} else {
					catalogue = append(catalogue, Item{
						Weight: weight,
						Value:  value,
					})
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("error reading file")
		return
	}

	population := generateInitPopulation()

	bestIndividual := population[0]
	for generation := 0; generation < generationLimit; generation++ {
		population = evolve(population)
		if population[0].Fitness > bestIndividual.Fitness {
			bestIndividual = population[0]
		}
	}

	fmt.Println(bestIndividual.Fitness)
}
