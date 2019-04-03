package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"hash"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/mem"

	"github.com/shirou/gopsutil/cpu"
)

var (
	hasher      string
	wordlist    string
	hashMethods int
	verbose     bool
	Vverbose    bool
	Benchmark   bool
	InfoHw      bool
	HelpFlag    bool
	err         error

	wg         sync.WaitGroup
	hashFind   string
	count      uint16
	hashToFind string
)

func init() {
	flag.StringVar(&hasher, "H", "", "Hash to test")
	flag.BoolVar(&HelpFlag, "h", false, "Show this help")
	flag.StringVar(&wordlist, "w", "", "Wordlist for hash generating")
	flag.BoolVar(&verbose, "v", false, "Verbose program")
	flag.BoolVar(&Vverbose, "vv", false, "More more verbose program")
	flag.BoolVar(&Benchmark, "b", false, "Benchmark mode")
	flag.BoolVar(&InfoHw, "I", false, "Info hardware")
	flag.IntVar(&hashMethods, "m", 0, "Hash methods\n1 \t: MD5\n2 \t: SHA1\n3 \t: SHA256\n")

}

const (
	DEFAULT_BENCHMARK_VALUE float64 = 10000000
	KB                      float64 = 1000
	MB                      float64 = 1000000
	GB                      float64 = 1000000000
	TB                      float64 = 1000000000000
)

func main() {

	flag.Parse()

	if Benchmark == true {
		BenchHash(hashMethods)
		os.Exit(0)
	} else if InfoHw == true {
		InfoHardWare()
		os.Exit(0)
	} else if HelpFlag == true {
		flag.Usage()
		os.Exit(0)
	}
	if hasher == "" || hashMethods == 0 || hashMethods > 3 {
		flag.Usage()
		os.Exit(0)
	}

	hashToFind = strings.ToLower(hasher)

	wg.Add(1)
	start := time.Now()
	fmt.Printf("Hash mode %s\n", GetHashName(hashMethods))

	go ReadWordListFile(wordlist)

	wg.Wait()

	fmt.Println(hashFind)
	if verbose == true {
		fmt.Printf("Password count : %d\n", count)
	}
	fmt.Printf("Time elapsed \t : %s\n", time.Since(start))

}

func ReadWordListFile(wordlits string) {
	defer wg.Done()
	f, _ := os.Open(wordlist)
	scanner := bufio.NewScanner(f)

	hasher := GetHash(hashMethods)

	var line string
	var result string
	found := false
	for scanner.Scan() {

		line = scanner.Text()
		hasher.Write([]byte(line))
		result = hex.EncodeToString(hasher.Sum(nil))
		count++
		if Vverbose == true {
			fmt.Printf("Password : %s \t hash : %s\n", line, result)
		}
		if result == hashToFind {
			hashFind = "hash found : " + line
			found = true
		}
		hasher.Reset()

	}

	if found == false {
		hashFind = "hash not found"
	}

}

// GetHash func
func GetHash(choice int) hash.Hash {

	switch choice {
	case 1:
		tmp := md5.New()
		return tmp
	case 2:
		tmp := sha1.New()
		return tmp
	case 3:
		tmp := sha256.New()
		return tmp

	}
	return nil
}

// GetHashName func
func GetHashName(choice int) string {
	switch choice {
	case 1:
		tmp := "md5"
		return tmp
	case 2:
		tmp := "sha1"
		return tmp
	case 3:
		tmp := "sha256"
		return tmp

	}
	return ""
}

func BenchHash(choice int) {

	start := time.Now()
	var hasher hash.Hash
	if choice == 0 {
		choice = 1
	}
	hasher = GetHash(choice)
	fmt.Printf("Benchmark for %s methods\n", GetHashName(choice))
	for i := 0; i < int(DEFAULT_BENCHMARK_VALUE); i++ {
		hasher.Write([]byte(strconv.Itoa(i)))
		result := hex.EncodeToString(hasher.Sum(nil))
		if verbose == true {
			fmt.Printf("%d : %s\n", i, result)
		}

	}
	stop := time.Since(start)
	fmt.Printf("Time elapsed : %v for %.f hash\n", stop, DEFAULT_BENCHMARK_VALUE)

	hTime := DEFAULT_BENCHMARK_VALUE / stop.Seconds()
	var result string
	if hTime > KB && hTime < MB {
		result = strconv.FormatFloat(hTime/1000, 'f', 6, 64) + " KH"
	} else if hTime > MB && hTime < GB {
		result = strconv.FormatFloat(hTime/1000000, 'f', 6, 64) + " MH"
	} else if hTime > GB && hTime < TB {
		result = strconv.FormatFloat(hTime/1000000000, 'f', 6, 64) + " GH"
	}
	fmt.Printf("Bench : %s/s\n", result)

}

func InfoHardWare() {
	C, err := cpu.Info()
	if err != nil {
		panic(err)
	} else {

		fmt.Println("CPU INFO :")
		fmt.Printf("CPU \t: %v\nCores \t: %v\nFamily \t: %v\n", C[0].ModelName, C[0].Cores, C[0].Family)

	}
	M, err := mem.VirtualMemory()
	if err != nil {
		panic(err)
	} else {

		fmt.Println("MEMORY INFO :")
		mem := [2]uint64{M.Total, M.Used}
		for i := 0; i < 2; i++ {
			var result string
			if mem[i] > uint64(KB) && mem[i] < uint64(MB) {
				result = strconv.FormatFloat(float64(mem[i])/1000, 'f', 6, 64) + " KB"
			} else if mem[i] > uint64(MB) && mem[i] < uint64(GB) {
				result = strconv.FormatFloat(float64(mem[i])/1000000, 'f', 6, 64) + " MB"
			} else if mem[i] > uint64(GB) && mem[i] < uint64(TB) {
				result = strconv.FormatFloat(float64(mem[i])/1000000000, 'f', 6, 64) + " GB"
			}
			if i == 0 {
				fmt.Printf("Total : %s\n", result)
			} else if i == 1 {
				fmt.Printf("Used : %s\n", result)
			}

		}

	}
}
