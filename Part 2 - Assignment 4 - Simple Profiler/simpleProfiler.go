package main

import (
	"fmt"
	"runtime"
)

type profile_anchor struct {
	TSCElapsed         uint64
	TSCElapsedChildren uint64
	hitCount           uint64
	label              string
}

type profiler struct {
	Anchors        [4096]profile_anchor
	numAnchorsUsed int
	startTSC       uint64
	endTSC         uint64
}

var GlobalProfiler profiler
var GlobalProfilerParent int

type profile_block struct {
	label       string
	startTSC    uint64
	anchorIndex int
	parentIndex int
}

// func (p *profile_block) profile_block_constructor(label_ string, anchorIndex_ int) {
// 	p.anchorIndex = anchorIndex_
// 	p.label = label_
// 	p.startTSC = ReadCPUTimer()
// }

// func (p *profile_block) profile_block_destructor() {
// 	Elapsed := ReadCPUTimer() - p.startTSC
// 	anchor := &GlobalProfiler.Anchors[p.anchorIndex]
// 	anchor.TSCElapsed += Elapsed
// 	anchor.hitCount++
// 	anchor.label = p.label
// }

func getCurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(2)
	return fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
}
func getCurrentLineNumber() string {
	_, _, lineNo, _ := runtime.Caller(2)
	return fmt.Sprintf("%v", lineNo)
}

func TimeBlock(name string) *profile_block {

	GlobalProfiler.numAnchorsUsed++
	newProfileBlock := profile_block{name, ReadCPUTimer(), GlobalProfiler.numAnchorsUsed, GlobalProfilerParent}

	GlobalProfilerParent = GlobalProfiler.numAnchorsUsed

	return &newProfileBlock
}

func TimeFunction() *profile_block {
	return TimeBlock(getCurrentFuncName())
}

func TimeFunctionEnd(p *profile_block) {

	Elapsed := ReadCPUTimer() - p.startTSC

	GlobalProfilerParent = p.parentIndex

	parent := &GlobalProfiler.Anchors[p.parentIndex]
	parent.TSCElapsedChildren += Elapsed
	anchor := &GlobalProfiler.Anchors[p.anchorIndex]
	anchor.TSCElapsed += Elapsed
	anchor.hitCount++
	anchor.label = p.label
}

func PrintTimeElapsed(totalTSCElapsed uint64, anchor *profile_anchor) {
	elapsed := anchor.TSCElapsed - anchor.TSCElapsedChildren
	percent := 100.0 * (float64(elapsed) / float64(totalTSCElapsed))
	fmt.Printf("  %s: hits:%v: elapsed:%v percent:(%.2f)", anchor.label, anchor.hitCount, elapsed, percent)
	if anchor.TSCElapsedChildren > 0 {
		percentWithChildren := 100.0 * (float64(anchor.TSCElapsed) / float64(totalTSCElapsed))
		fmt.Printf(", %.2f w/children", percentWithChildren)
	}
	fmt.Println()
}

func BeginProfile() {
	GlobalProfiler.startTSC = ReadCPUTimer()
	// GlobalProfiler.numAnchorsUsed = -1
}

func EndAndPrintProfile() {
	GlobalProfiler.endTSC = ReadCPUTimer()
	CPUFreq := EstimateCPUTimerFreq()

	TotalCPUElapsed := GlobalProfiler.endTSC - GlobalProfiler.startTSC

	if CPUFreq != 0 {
		fmt.Printf("\nTotal time: %0.4fms (CPU freq %v)\n", 1000.0*float64(TotalCPUElapsed)/float64(CPUFreq), CPUFreq)
	}

	for anchorIndex := range GlobalProfiler.numAnchorsUsed + 1 {
		anchor := &GlobalProfiler.Anchors[anchorIndex]
		if anchor.TSCElapsed > 0 {
			PrintTimeElapsed(TotalCPUElapsed, anchor)
		}
	}
}
