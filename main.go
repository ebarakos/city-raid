package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type city struct {
	name     string
	north    *city
	south    *city
	east     *city
	west     *city
	occupied *alien
}
type alien struct {
	id       string
	location *city
}

func (a *alien) getAvailableDirections() []*city {
	var directions []*city
	if a.location.north != nil {
		directions = append(directions, a.location.north)
	}
	if a.location.south != nil {
		directions = append(directions, a.location.south)
	}
	if a.location.east != nil {
		directions = append(directions, a.location.east)
	}
	if a.location.west != nil {
		directions = append(directions, a.location.west)
	}
	return directions
}

func destroyAliens(a1 *alien, a2 *alien) {
	for i, a := range aliens {
		if a == a1 {
			aliens = append(aliens[:i], aliens[i+1:]...)
			break
		}
	}
	for i, a := range aliens {
		if a == a2 {
			aliens = append(aliens[:i], aliens[i+1:]...)
			break
		}
	}
	return
}

func (a *alien) move() bool {
	oneAlien := a
	directions := a.getAvailableDirections()
	if len(directions) == 0 {
		return false
	}
	fmt.Println()
	fmt.Println("Alien ", a.id, " moves: ")
	a.location.occupied = nil
	rand.Seed(time.Now().Unix())
	chosenDirection := directions[rand.Intn(len(directions))]
	var occupyingAlienToDestroy *alien
	otherAlien := chosenDirection.occupied
	if chosenDirection.occupied != nil {
		fmt.Println(chosenDirection.name, " has​ ​ been​ ​ destroyed​ ​ by​ ​ alien ", a.id, " and alien: ", chosenDirection.occupied.id, "!")
		occupyingAlienToDestroy = chosenDirection.occupied

		chosenDirection.destroy()
	} else {
		a.location = chosenDirection
	}
	if occupyingAlienToDestroy == nil {
		a.location.occupied = a

	} else {
		a = nil
	}
	destroyAliens(oneAlien, otherAlien)
	return true

}

func (c *city) destroy() {
	for i, iCity := range cities {
		if c.name == iCity.name {
			cities = append(cities[:i], cities[i+1:]...)
			break
		}
	}
	for i, iCity := range cities {
		if iCity.north != nil {
			if c.name == iCity.north.name {
				cities[i].north = nil
			}
		}
		if iCity.south != nil {
			if c.name == iCity.south.name {
				cities[i].south = nil
			}
		}
		if iCity.west != nil {
			if c.name == iCity.west.name {
				cities[i].west = nil
			}
		}
		if iCity.east != nil {
			if c.name == iCity.east.name {
				cities[i].east = nil
			}
		}
	}
}

func printCityState() {
	for _, c := range cities {
		cityOutput := c.name
		if c.south != nil {
			cityOutput += " south=" + c.south.name
		}
		if c.north != nil {
			cityOutput += " north=" + c.north.name
		}
		if c.west != nil {
			cityOutput += " west=" + c.west.name
		}
		if c.east != nil {
			cityOutput += " east=" + c.east.name
		}
		fmt.Print(cityOutput)
		if c.occupied != nil {
			fmt.Println(" occupied by: " + c.occupied.id)
		} else {
			fmt.Println("")
		}

	}
}

func printAlienState() {
	for _, a := range aliens {
		fmt.Print("Alien: ", a.id)
		fmt.Println(" in location: ", a.location.name)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func appendRoads(cityRef *city, direction string, cityName string) {
	index := len(cities) - 1
	for i, iCity := range cities {
		if iCity.name == cityName {
			index = i
			break
		}
	}
	if index == len(cities)-1 {
		cities = append(cities, &city{name: cityName})
		index = len(cities) - 1
	}
	// TODO: Direct string comparison not working
	if strings.Contains(direction, "north") {
		cityRef.north = cities[index]
	}
	if strings.Contains(direction, "south") {
		cityRef.south = cities[index]
	}
	if strings.Contains(direction, "east") {
		cityRef.east = cities[index]
	}
	if strings.Contains(direction, "west") {
		cityRef.west = cities[index]
	}
	return
}

func play(numAliens int, inputFileName, outputFileName string) {
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		var tempCity *city
		for i, iCity := range cities {
			if words[0] == iCity.name {
				tempCity = cities[i]
				break
			}
		}
		if tempCity == nil {
			cities = append(cities, &city{name: words[0]})
			tempCity = cities[len(cities)-1]
		}
		for _, element := range words[1:] {
			splittedValues := strings.Split(element, "=")
			if len(splittedValues) == 2 {
				appendRoads(tempCity, splittedValues[0], splittedValues[1])
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	if numAliens > len(cities) {
		numAliens = len(cities)
	}
	rand.Seed(time.Now().Unix())
	for i := 1; i <= numAliens; i++ {
		var chosenCity *city
		for {
			chosenCity = cities[rand.Intn(len(cities))]

			if chosenCity.occupied == nil {
				aliens = append(aliens, &alien{id: strconv.Itoa(i), location: chosenCity})
				chosenCity.occupied = aliens[len(aliens)-1]
				break
			}

		}

	}
	fmt.Println("Initial city state: ")
	printCityState()
	fmt.Println("Initial alien state: ")
	printAlienState()
	for i := 1; i <= 10000; i++ {
		//TODO: also break if all aliens are trapped?
		if len(aliens) == 0 {
			fmt.Println("No more aliens")
			break
		}
		for _, a := range aliens {
			a.move()
			printCityState()
			printAlienState()

		}
	}

	fileOutput := ""
	for _, c := range cities {
		cityOutput := c.name
		if c.south != nil {
			cityOutput += " south=" + c.south.name
		}
		if c.north != nil {
			cityOutput += " north=" + c.north.name
		}
		if c.west != nil {
			cityOutput += " west=" + c.west.name
		}
		if c.east != nil {
			cityOutput += " east=" + c.east.name
		}

		fileOutput += cityOutput + "\n"
	}
	d1 := []byte(fileOutput)
	error := ioutil.WriteFile(outputFileName, d1, 0644)
	check(error)
}

var cities = make([]*city, 0)
var aliens = make([]*alien, 0)

func main() {
	numAliens := flag.Int("aliens", 10000, "number of aliens")
	inputFileName := flag.String("input", "input", "path to input file")
	outputFileName := flag.String("output", "output", "path to output file")
	flag.Parse()
	play(*numAliens, *inputFileName, *outputFileName)
}
