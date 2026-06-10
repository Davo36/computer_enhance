package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"time"
)

const maxPairs = 10_000_000
const numClusters = 64

type Cluster struct {
	xCenter int
	xRadius int
	yCenter int
	yRadius int
}

func main() {

	startTime := time.Now()

	wordPtr := flag.String("type", "cluster", "uniform or cluster")
	seedPtr := flag.Int("seed", 234089, "random seed") // DNB: NB: I'm not using a random seed.  Go automatically uses a new seed each time the program is run.
	numPairsPtr := flag.Int("pairs", 1_000, "number of pairs of points")

	flag.Parse()

	fmt.Println()
	fmt.Println("Uniform/cluster:", *wordPtr)
	fmt.Println("Seed:", *seedPtr)
	fmt.Println("Number of pairs:", *numPairsPtr)

	if *numPairsPtr > maxPairs || *numPairsPtr < 1 {
		fmt.Println("Num pairs must be between 1 and", maxPairs)
		return
	}

	MaxAllowedX := 180
	MaxAllowedY := 90
	sum := 0.0
	sumCoef := 1.0 / float64(*numPairsPtr)

	clustering := false
	clusters := []Cluster{}
	numPointsInEachCluster := 0
	if *wordPtr == "cluster" {
		clustering = true
		numPointsInEachCluster = int(math.Ceil(float64(*numPairsPtr) / float64(numClusters)))
		for range numClusters {
			clusters = append(clusters, Cluster{rand.IntN(MaxAllowedX*2) - MaxAllowedX,
				rand.IntN(MaxAllowedX) + 1,
				rand.IntN(MaxAllowedY*2) - MaxAllowedY,
				rand.IntN(MaxAllowedY) + 1})
		}
	} else if *wordPtr != "uniform" {
		fmt.Println("Usage: -type=uniform | cluster -seed=234089 -pairs=1000000")
	}

	jsonFile, _ := os.Create("pairsData.json")
	defer jsonFile.Close()
	jsonFile.WriteString("{\"pairs\":[\n")
	haverAnswersFile, _ := os.Create("haversAnswers.bin")
	defer haverAnswersFile.Close()
	haverAnswersTextFile, _ := os.Create("haversAnswers.txt")
	defer haverAnswersTextFile.Close()

	for i := 0; i < *numPairsPtr; i++ {
		var x0, y0, x1, y1 float64
		if clustering {
			clusterNumber := i / numPointsInEachCluster
			x0 = float64(rand.IntN(clusters[clusterNumber].xRadius) + clusters[clusterNumber].xCenter)
			x1 = float64(rand.IntN(clusters[clusterNumber].xRadius) + clusters[clusterNumber].xCenter)
			y0 = float64(rand.IntN(clusters[clusterNumber].yRadius) + clusters[clusterNumber].yCenter)
			y1 = float64(rand.IntN(clusters[clusterNumber].yRadius) + clusters[clusterNumber].yCenter)
		} else {
			x0 = float64(rand.IntN(MaxAllowedX*2) - MaxAllowedX)
			x1 = float64(rand.IntN(MaxAllowedX*2) - MaxAllowedX)
			y0 = float64(rand.IntN(MaxAllowedY*2) - MaxAllowedY)
			y1 = float64(rand.IntN(MaxAllowedY*2) - MaxAllowedY)
		}
		JSONSep := ",\n"
		if i == (*numPairsPtr - 1) {
			JSONSep = "\n"
		}
		fmt.Fprintf(jsonFile, "    {\"x0\":%.16f, \"y0\":%.16f, \"x1\":%.16f, \"y1\":%.16f}%s", x0, y0, x1, y1, JSONSep)

		earthRadius := 6372.8
		haversineDistance := ReferenceHaversine(x0, y0, x1, y1, earthRadius)
		sum += sumCoef * haversineDistance
		binary.Write(haverAnswersFile, binary.LittleEndian, haversineDistance)
		fmt.Fprintln(haverAnswersTextFile, haversineDistance)
	}
	fmt.Fprintf(jsonFile, "]}\n")
	binary.Write(haverAnswersFile, binary.LittleEndian, sum) // Put the average value in the last pos of the file.
	fmt.Println("Expected sum:", sum)

	fmt.Println("Time taken:", time.Since(startTime))

}
