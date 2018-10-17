package main

import (
	"bufio"
	"container/ring"
	"encoding/binary"
	"fmt"
	"os"
)

func main() {

	f, err := os.Open("Enc1.dat")
	if err != nil {
		panic(err.Error())
	}

	reader := bufio.NewReader(f)
	reader.Discard(6)

	xor32Key := []uint32{0x3F08A79B, 0xE25CC287, 0x93D27AB9, 0x20DEA7BF}
	xor32KeyRingBuffer := ring.New(4)

	for i := 0; i < xor32KeyRingBuffer.Len(); i++ {
		xor32KeyRingBuffer.Value = xor32Key[i]
		xor32KeyRingBuffer = xor32KeyRingBuffer.Next()
	}

	encryptionKeys := make([]uint32, 12)

	for i := 0; i < 12; i++ {
		xor32DecryptionKeysByteArray := make([]byte, 4)
		reader.Read(xor32DecryptionKeysByteArray)
		encryptionKeys[i] = binary.LittleEndian.Uint32(xor32DecryptionKeysByteArray) ^ xor32KeyRingBuffer.Value.(uint32)
		xor32KeyRingBuffer = xor32KeyRingBuffer.Next()
	}

	fmt.Printf("%v\n\n", encryptionKeys)

	keys := findOtherKey(encryptionKeys)

	decryptionKeys := encryptionKeys
	for i := 0; i < len(keys); i++ {
		decryptionKeys[i+4] = keys[i]
	}

	fmt.Printf("%v\n\n", decryptionKeys)

}

func findOtherKey(encryptionKeys []uint32) []uint32 {
	keys := make([]uint32, 4)
	for i := 0; i < 4; i++ {
		keys[i] = findKey(encryptionKeys[i], encryptionKeys[i+4])
	}
	return keys
}

func findKey(modulus, key uint32) uint32 {
	fmt.Printf("Finding mod %v, key %v: ", modulus, key)

	var i uint32
	for i = 0; i < modulus; i++ {
		if key*i%modulus == 1 {
			fmt.Printf("%v\n\n", i)
			return i
		}
	}
	return 0
}
