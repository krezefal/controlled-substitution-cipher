package main

import (
	"log"
	"math/rand"
)

func encrypt(rounds int, byteRow []byte, keys [16]uint8, subsTables [32][256]uint8) []byte {

	var result = make([]byte, len(byteRow))

	for r := 0; r < rounds; r++ {
		for i := 0; i < len(byteRow); i += 4 {
			result[i] = byteRow[i] ^ keys[0 + r * 4]
			tableNum := result[i] % 32
			for j := 1; j < 4; j++ {
				result[i + j] = byteRow[i + j] ^ keys[j + r * 4]
				result[i + j] = subsTables[tableNum][result[i + j]]
				tableNum = (result[i + j] % 32) ^ tableNum
			}
		}
	}

	return result
}

func decrypt(rounds int, byteRow []byte, keys [16]uint8, subsTables [32][256]uint8) []byte {

	var result = make([]byte, len(byteRow))

	for r := rounds - 1; r >= 0; r++ {
		for i := 0; i < len(byteRow); i += 4 {
			tableNum := result[i] % 32
			result[i] = byteRow[i] ^ keys[0 + r * 4]
			for j := 1; j < 4; j++ {
				tmp := tableNum
				tableNum = (result[i + j] % 32) ^ tableNum
				result[i + j] = subsTables[tmp + 16][result[i + j]]
				result[i + j] = byteRow[i + j] ^ keys[j + r * 4]
			}
		}
	}

	return result
}

func main() {

	rounds := 1
	var keys [16]uint8 // max 4 rounds with different keys
	var subsTables [32][256]uint8

	for i := range keys {
        keys[i] = uint8(rand.Intn(256))
    }

	for i := 0; i < 16; i++ {
		availableValues := make([]uint8, 256)
		for i := range availableValues {
			availableValues[i] = uint8(i)
		}
		j := 0
		for len(availableValues) != 0 {
			randIdx := rand.Intn(len(availableValues))
			randVal := availableValues[randIdx]
			availableValues = remove(availableValues, randIdx)
			subsTables[i][j] = randVal
			j++
		}
	}
	for i := 16; i < 32; i++ {
		for j := 1; j < 256; j++ {
			subsTables[i][subsTables[i][j]] = uint8(j)
		}
	}

	img, byteRow := getFile("examples/ford.bmp")

	encRow := encrypt(rounds, byteRow, keys, subsTables)
	err := saveFile(img, encRow, "examples/enc_res.bmp")
	if err != nil {
		log.Println(err.Error())
	}

	decRow := decrypt(rounds, encRow, keys, subsTables)
	err = saveFile(img, decRow, "examples/dec_res.bmp")
	if err != nil {
		log.Println(err.Error())
	}

}
