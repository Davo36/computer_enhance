package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/dterei/gotsc"
)

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

func ReadCPUTimer() uint64 {
	return gotsc.BenchStart()
}

func EstimateCPUTimerFreq() uint64 {

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

// func PrintTimeElapsed(label string, totalTSCElapsed uint64, begin uint64, end uint64) {
// 	elapsed := end - begin
// 	percent := 100.0 * (float64(elapsed) / float64(totalTSCElapsed))
// 	fmt.Printf("  %s: %v (%.2f)\n", label, elapsed, percent)
// }
