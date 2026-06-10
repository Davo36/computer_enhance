package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/dterei/gotsc"
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

func readJSONFile() []byte {
	data, _ := os.ReadFile("..//Part 2 -  Assignment 1 - Generating Haversine Input//pairsData.json")
	return data
}

func parseJSONData(data []byte) []pair {

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

func readOSTimer() uint64 {

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	qpc := kernel32.NewProc("QueryPerformanceCounter")

	var counter int64
	ret, _, _ := qpc.Call(uintptr(unsafe.Pointer(&counter)))
	if ret == 0 {
		fmt.Println("QueryPerformanceCounter failed")
		return 0
	}
	// fmt.Printf("QueryPerformanceCounter   : %d\n", counter)

	return uint64(counter)
}

func getOSTimerFrequency() uint64 {

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	qpf := kernel32.NewProc("QueryPerformanceFrequency")

	var freq int64
	ret, _, _ := qpf.Call(uintptr(unsafe.Pointer(&freq)))
	if ret == 0 {
		fmt.Println("QueryPerformanceFrequency failed")
		return 0
	}
	// fmt.Printf("QueryPerformanceFrequency: %d ticks/sec (≈ %.2f MHz)\n", freq, float64(freq)/1000000)

	return uint64(freq)
}

func readCPUTimer() uint64 {
	return gotsc.BenchStart()
}

func estimateCPUTimerFreq() uint64 {

	MillisecondsToWait := uint64(100)
	OSFreq := getOSTimerFrequency()

	CPUStart := gotsc.BenchStart()
	OSStart := readOSTimer()
	OSEnd := uint64(0)
	OSElapsed := uint64(0)
	OSWaitTime := OSFreq * MillisecondsToWait / 1000
	for OSElapsed < OSWaitTime {
		OSEnd = readOSTimer()
		OSElapsed = OSEnd - OSStart
	}

	CPUEnd := gotsc.BenchEnd()
	CPUElapsed := CPUEnd - CPUStart

	CPUFreq := uint64(0)
	if OSElapsed > 0 {
		CPUFreq = OSFreq * CPUElapsed / OSElapsed
	}

	return CPUFreq
}

func PrintTimeElapsed(label string, totalTSCElapsed uint64, begin uint64, end uint64) {
	elapsed := end - begin
	percent := 100.0 * (float64(elapsed) / float64(totalTSCElapsed))
	fmt.Printf("  %s: %v (%.2f)\n", label, elapsed, percent)
}

func main() {

	profBegin := readCPUTimer()

	profRead := readCPUTimer()
	data := readJSONFile()
	profMiscSetup := readCPUTimer()

	profParse := readCPUTimer()
	pairs := parseJSONData(data)

	profSum := readCPUTimer()
	theSum := computeSum(pairs)
	profMiscOutput := readCPUTimer()

	fmt.Println("Pair count:", len(pairs))
	fmt.Println("Haversine sum:", theSum)

	referenceSum := getAnswer()
	fmt.Println("\nValidation:")
	fmt.Println("Reference sum:", referenceSum)
	fmt.Println("Difference:", referenceSum-theSum)

	profEnd := readCPUTimer()
	TotalCPUElapsed := profEnd - profBegin

	CPUFreq := estimateCPUTimerFreq()
	fmt.Printf("\nTotal time: %0.4fms (CPU freq %v)\n", 1000.0*float64(TotalCPUElapsed)/float64(CPUFreq), CPUFreq)

	PrintTimeElapsed("Startup", TotalCPUElapsed, profBegin, profRead)
	PrintTimeElapsed("Read", TotalCPUElapsed, profRead, profMiscSetup)
	PrintTimeElapsed("MiscSetup", TotalCPUElapsed, profMiscSetup, profParse)
	PrintTimeElapsed("Parse", TotalCPUElapsed, profParse, profSum)
	PrintTimeElapsed("Sum", TotalCPUElapsed, profSum, profMiscOutput)
	PrintTimeElapsed("MiscOutput", TotalCPUElapsed, profMiscOutput, profEnd)

}
