package main

import (
	"crypto/sha256"
	"fmt"

	"github.com/dustin/go-humanize"

	"runtime"
	"time"
)

func main() {
	pow("hello world", 30)
}

func pow(prefix string, bitLength int) {
	start := time.Now()
	var hash []int
	totalHashesProcessed := 0

	numberOfCPU := runtime.NumCPU()
	closeChan := make(chan int, 1)
	solutionChan := make(chan []byte, 1)
	for idx := 0; idx < numberOfCPU; idx++ {
		hash = append(hash, 0)
		go func(hashIndex int) {
			seed := uint64(time.Now().Local().UnixNano())
			randomBytes := make([]byte, 20)
			randomBytes = append([]byte(prefix), randomBytes...)

			for {
				select {
				case <-closeChan:
					closeChan <- 1
					return
				case <-time.After(time.Nanosecond):
					count := 0
					for count < 5000 {
						count++
						seed = RandomString(randomBytes, len(prefix), seed)
						if Hash(randomBytes, bitLength) {
							hash[hashIndex] += count
							solutionChan <- randomBytes
							closeChan <- 1
							return
						}
					}
					hash[hashIndex] += count
				}
			}
		}(idx)
	}

	solution := <-solutionChan
	hashHex := sha256.Sum256(solution)

	for _, v := range hash {
		totalHashesProcessed += v
	}

	end := time.Now()
	fmt.Println("Result:", solution)
	fmt.Println("Result (hex):", "0x"+fmt.Sprintf("%x", hashHex))
	fmt.Println("time:", end.Sub(start).Seconds())
	fmt.Println("processed:",
		humanize.Comma(int64(totalHashesProcessed)))
	fmt.Printf("processed/sec: %s\n", humanize.Comma(int64(
		float64(totalHashesProcessed)/end.Sub(start).Seconds())))
}

var characterSet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandomString(str []byte, offset int, seed uint64) uint64 {
	for i := offset; i < len(str); i++ {
		seed = RandomNumber(seed)
		str[i] = characterSet[seed%62]
	}
	return seed
}

func Hash(data []byte, bits int) bool {
	bs := sha256.Sum256(data)
	nbytes := bits / 8
	nbits := bits % 8

	for idx := 0; idx < nbytes; idx++ {
		if bs[idx] > 0 {
			return false
		}
	}

	return (bs[nbytes] >> (8 - nbits)) == 0
}

func RandomNumber(seed uint64) uint64 {
	seed ^= seed << 21
	seed ^= seed >> 35
	seed ^= seed << 4
	return seed
}
