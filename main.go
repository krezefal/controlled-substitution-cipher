package main

import (
	"fmt"
	"math/rand"
	"time"
)

// 1 subblock = 8 bit; 1 key = 8 bit value; 1 row of substitution table designed for 8 bit value.
// Amount of keys on 1 round should be = amount of subblocks for better encryption.
// E.g.: 8 rounds is max available value for 256 capacity of keys and 32 amount of subblocks (256 / 8 = 32)

const rounds = 8
const subblocksNum = 32
const keysCapacity = rounds * subblocksNum

func encrypt(rounds int, byteRow []byte, keys [keysCapacity]uint8, subsTables [32][256]uint8) []byte {

	result := make([]byte, len(byteRow))
	copy(result, byteRow)

	for r := 0; r < rounds; r++ {
		for i := 0; i < len(result); i += subblocksNum {
			result[i] = result[i] ^ keys[r*32]
			tableNum := result[i] % 32
			for j := 1; j < subblocksNum; j++ {
				result[i+j] = result[i+j] ^ keys[j+r*32]
				result[i+j] = subsTables[tableNum][result[i+j]]
				tableNum = (result[i+j] % 32) ^ tableNum
			}
		}
	}

	return result
}

func decrypt(rounds int, byteRow []byte, keys [keysCapacity]uint8, subsTables [32][256]uint8) []byte {

	result := make([]byte, len(byteRow))
	copy(result, byteRow)

	for r := rounds - 1; r >= 0; r-- {
		for i := 0; i < len(result); i += subblocksNum {
			tableNum := result[i] % 32
			result[i] = result[i] ^ keys[r*32]
			for j := 1; j < subblocksNum; j++ {
				tmp := tableNum
				tableNum = (result[i+j] % 32) ^ tableNum
				result[i+j] = subsTables[(tmp+16)%32][result[i+j]]
				result[i+j] = result[i+j] ^ keys[j+r*32]
			}
		}
	}

	return result
}

func main() {

	var keys [keysCapacity]uint8
	var subsTables [32][256]uint8

	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)

	for i := range keys {
		keys[i] = uint8(r.Intn(256))
	}

	fmt.Println("ROUND KEYS:", keys)

	for i := 0; i < 16; i++ {
		availableValues := make([]uint8, 256)
		for i := range availableValues {
			availableValues[i] = uint8(i)
		}
		j := 0
		for len(availableValues) != 0 {
			randIdx := r.Intn(len(availableValues))
			randVal := availableValues[randIdx]
			availableValues = remove(availableValues, randIdx)
			subsTables[i][j] = randVal
			j++
		}
	}
	for i := 16; i < 32; i++ {
		for j := 1; j < 256; j++ {
			subsTables[i][subsTables[i-16][j]] = uint8(j)
		}
	}

	img, byteRow := getFile("examples/ford.bmp")

	encRow := encrypt(rounds, byteRow, keys, subsTables)
	saveFile(img, encRow, "examples/enc_res.bmp")

	decRow := decrypt(rounds, encRow, keys, subsTables)
	saveFile(img, decRow, "examples/dec_res.bmp")

}
