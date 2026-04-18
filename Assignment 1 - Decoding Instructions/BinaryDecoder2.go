package main

// This works for both files i.e. a single instruction and multiple instructions

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Each instruction consists of 2 bytes.
type Instruction struct {
	Byte1 uint8
	Byte2 uint8
}

var registerTable = [][]string{
	{"al", "ax"},
	{"cl", "cx"},
	{"dl", "dx"},
	{"bl", "bx"},
	{"ah", "sp"},
	{"ch", "bp"},
	{"dh", "si"},
	{"bh", "di"},
}

// Read in a binary file and decode it into a list of instructions
func processFile(inputFileName string, outputFileName string) {

	// Input file.  This is a binary file
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Found this later, might have been handy...
	// stats, err := file.Stat()
	// size := stats.Size()
	// fmt.Println(size)

	// Read the binary data into the struct
	instruction := Instruction{}
	instructions := []Instruction{}
	for {
		err = binary.Read(file, binary.LittleEndian, &instruction)
		if err != nil { // EOF
			break
		}
		instructions = append(instructions, instruction)
		// fmt.Println(instruction)
	}

	// Create the output file, this is a text file.
	outFile, _ := os.Create(outputFileName)
	defer outFile.Close()
	outFile.WriteString("bits 16\n\n")
	for _, instruction := range instructions {
		// Get string representations of the bytes to make accessing the bits nice and easy.
		byte1Str := fmt.Sprintf("%08b", instruction.Byte1)
		byte2Str := fmt.Sprintf("%08b", instruction.Byte2)
		// Create the instruction string
		instructionStr := "mov "
		w, _ := strconv.ParseInt(byte1Str[7:], 2, 64)
		fromRegister, _ := strconv.ParseInt(byte2Str[5:], 2, 64)
		fromRegisterLabel := registerTable[fromRegister][w]
		toRegister, _ := strconv.ParseInt(byte2Str[2:5], 2, 64)
		toRegisterLabel := registerTable[toRegister][w]
		instructionStr += fromRegisterLabel + ", " + toRegisterLabel + "\n"
		outFile.WriteString(instructionStr)
		fmt.Print(instructionStr)
	}
}

func main() {

	processFile("listing_0037_single_register_mov", "Davids_single_register_mov.asm")
	processFile("listing_0038_many_register_mov", "Davids_many_register_mov.asm")

}
