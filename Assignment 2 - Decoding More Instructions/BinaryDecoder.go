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
type InstructionByte struct {
	Byte1 uint8
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

var registerMemoryFieldEncodingTable = []string{
	"[bx + si]",
	"[bx + di]",
	"[bp + si]",
	"[bp + di]",
	"[si]",
	"[di]",
	"[bp]", // Direct address
	"[bx]"}

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
	instruction := InstructionByte{}
	instructions := []InstructionByte{}
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
	IP := 0
	for IP < len(instructions) {

		// Get string representations of the bytes to make accessing the bits nice and easy.
		byte1Str := fmt.Sprintf("%08b", instructions[IP].Byte1)
		// fmt.Println(byte1Str)

		if byte1Str[:6] == "100010" {
			// Mov command.  As in last assignment.
			// fmt.Println("Standard move.")
			d, _ := strconv.ParseInt(byte1Str[6:7], 2, 64) // Direction.  From or to.
			w, _ := strconv.ParseInt(byte1Str[7:], 2, 64)  // Width.  8 or 16 bits.

			IP++
			byte2Str := fmt.Sprintf("%08b", instructions[IP].Byte1)
			mod, _ := strconv.ParseInt(byte2Str[:2], 2, 64) // Kind of move: 11 = register to register.
			reg, _ := strconv.ParseInt(byte2Str[2:5], 2, 64)
			rm, _ := strconv.ParseInt(byte2Str[5:], 2, 64)

			instructionStr := "mov "
			fromRegisterLabel := ""
			toRegisterLabel := ""

			if mod == 3 { // "11"
				// Register to register move.
				// Create the instruction string
				fromRegisterLabel = registerTable[rm][w]
				toRegisterLabel = registerTable[reg][w]
			} else if mod == 0 {
				// Source address calculation. No displacement
				// Create the instruction string
				fromRegisterLabel = registerMemoryFieldEncodingTable[rm]
				toRegisterLabel = registerTable[reg][w]
			} else if mod == 1 {
				//Source address calculation plus 8-bit displacement
				IP++
				dispLo := fmt.Sprintf("%08b", instructions[IP].Byte1)
				displacementValue, _ := strconv.ParseInt(dispLo, 2, 64)
				// Create the instruction string
				fromRegisterLabel = registerMemoryFieldEncodingTable[rm]
				if displacementValue != 0 {
					fromRegisterLabel = fromRegisterLabel[:len(fromRegisterLabel)-1] + " + " + strconv.Itoa(int(displacementValue)) + "]"
				}
				toRegisterLabel = registerTable[reg][w]
			} else if mod == 2 {
				//Source address calculation plus 16-bit displacement
				IP++
				dispLo := fmt.Sprintf("%08b", instructions[IP].Byte1)
				IP++
				dispHi := fmt.Sprintf("%08b", instructions[IP].Byte1)
				displacementValue, _ := strconv.ParseInt(dispHi+dispLo, 2, 64)
				fromRegisterLabel = registerMemoryFieldEncodingTable[rm]
				fromRegisterLabel = fromRegisterLabel[:len(fromRegisterLabel)-1] + " + " + strconv.Itoa(int(displacementValue)) + "]"
				toRegisterLabel = registerTable[reg][w]
			}
			if d == 1 {
				// Swap direction
				fromRegisterLabel, toRegisterLabel = toRegisterLabel, fromRegisterLabel
			}
			instructionStr += fromRegisterLabel + ", " + toRegisterLabel + "\n"
			outFile.WriteString(instructionStr)
			fmt.Print(instructionStr)

		} else if byte1Str[:4] == "1011" {
			// Move command.  Immediate to register move.
			// fmt.Println("Immediate move.")
			w, _ := strconv.ParseInt(byte1Str[4:5], 2, 64) // Width.  8 or 16 bits.
			reg, _ := strconv.ParseInt(byte1Str[5:], 2, 64)
			instructionStr := "mov "
			displacementValue := int64(0)
			IP++
			dispLo := fmt.Sprintf("%08b", instructions[IP].Byte1)
			fmt.Println(instructions[IP].Byte1, dispLo)
			if w == 0 {
				displacementValue, _ = strconv.ParseInt(dispLo, 2, 64)
			} else if w == 1 {
				IP++
				dispHi := fmt.Sprintf("%08b", instructions[IP].Byte1)
				displacementValue, _ = strconv.ParseInt(dispHi+dispLo, 2, 64)
			}
			fromRegisterLabel := strconv.Itoa(int(displacementValue))
			toRegisterLabel := registerTable[reg][w]
			instructionStr += toRegisterLabel + ", " + fromRegisterLabel + "\n"
			outFile.WriteString(instructionStr)
			fmt.Print(instructionStr)
		}

		IP++
	}
}

func main() {

	processFile("listing_0039_more_movs", "Davids_more_moves.asm")
	// processFile("listing_0040_challenge_movs", "Davids_many_register_mov.asm")

}
