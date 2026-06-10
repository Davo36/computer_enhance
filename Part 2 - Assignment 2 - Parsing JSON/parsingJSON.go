package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Parse the JSON data into an array of pair structs

type pair struct {
	x0 float64
	x1 float64
	y0 float64
	y1 float64
}

// Used to read float64s from the binary answers file
// type float64Value struct {
// 	num float64
// }

func getVal(line string, label string) (float64, error) {

	pos := strings.Index(line, label)

	if pos == -1 {
		return -99, errors.New("Value not found")
	}

	pos += 4
	pos2 := strings.Index(line[pos:], "}")
	pos3 := strings.Index(line[pos:], ",")
	pos4 := min(pos2, pos3) + pos
	if pos3 == -1 { // No comma on last line...
		pos4 = pos2 + pos
	}

	val, err := strconv.ParseFloat(line[pos:pos4], 64)
	return val, err

}

func loadJSONData() []pair {

	data, _ := os.ReadFile("..//Part 2 -  Assignment 1 - Generating Haversine Input//pairsData.json")
	lines := strings.Split(string(data), "\n")
	pairs := []pair{}
	for _, line := range lines[1 : len(lines)-2] { // TODO":Is this right?
		x0, err := getVal(line, "x0")
		if err != nil {
			log.Fatal("Can't load data.")
		}
		x1, err := getVal(line, "x1")
		if err != nil {
			log.Fatal("Can't load data.")
		}
		y0, err := getVal(line, "y0")
		if err != nil {
			log.Fatal("Can't load data.")
		}
		y1, err := getVal(line, "y1")
		if err != nil {
			log.Fatal("Can't load data.")
		}

		pairs = append(pairs, pair{x0, x1, y0, y1})
	}

	return pairs

}

func computeSum(pairs []pair) float64 {

	earthRadius := 6372.8
	theSum := 0.0
	sumCoef := 1.0 / float64(len(pairs))
	for _, thePair := range pairs {
		haversineDistance := ReferenceHaversine(thePair.x0, thePair.y0, thePair.x1, thePair.y1, earthRadius)
		theSum += sumCoef * haversineDistance
	}

	return theSum
}

func getAnswer() float64 {

	// Input file.  This is a binary file
	file, err := os.Open("..//Part 2 -  Assignment 1 - Generating Haversine Input//haversAnswers.bin")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Found this later, might have been handy...
	// stats, err := file.Stat()
	// size := stats.Size()
	// fmt.Println(size)

	// Read the binary data into the struct
	var val float64
	vals := []float64{}
	for {
		err = binary.Read(file, binary.LittleEndian, &val)
		if err != nil { // EOF
			break
		}
		vals = append(vals, val)
	}

	return vals[len(vals)-1] // The average value of the values is in the last position.
}

func main() {

	pairs := loadJSONData()
	fmt.Println("Num Pairs:", len(pairs))
	theSum := computeSum(pairs)
	fmt.Println("Haversine sum:", theSum)

	referenceSum := getAnswer()
	fmt.Println("\nValidation:")
	fmt.Println("Reference sum:", referenceSum)
	fmt.Println("Difference:", referenceSum-theSum)

}
