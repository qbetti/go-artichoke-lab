package main

import (
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/qbetti/go-artichoke/artichoke/generator"
	"github.com/qbetti/go-artichoke/artichoke/pas"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

const NUM_REPETITIONS = 10
const GROUP = "G1"
const DATA_DIR = "data"

func main() {

	CreateDirIfNotExist(DATA_DIR)

	actionSizes := []int{1, 500, 1000, 10000, 100000}
	var genDatas [][][]string
	var verifDatas [][][]string
	var spaceDatas [][][]string

	for _, actionSize := range actionSizes {
		genData := runGenerationExperiment(100, 2000, 100, actionSize)
		genDatas = append(genDatas, genData)

		verifData := runVerificationExperiment(100, 2000, 100, actionSize)
		verifDatas = append(verifDatas, verifData)

		spaceData := runSpaceExperiment(100, 2000, 100, actionSize)
		spaceDatas = append(spaceDatas, spaceData)
	}
	joinedGenDatas := joinExperimentsData(genDatas...)
	writeToCsv(joinedGenDatas, DATA_DIR + "/go-art-generation-time.csv")

	joinedVerifDatas := joinExperimentsData(verifDatas...)
	writeToCsv(joinedVerifDatas, DATA_DIR + "/go-art-verification-time.csv")

	joinedSpaceDatas := joinExperimentsData(spaceDatas...)
	writeToCsv(joinedSpaceDatas, DATA_DIR + "/go-art-space.csv")


	overheadDatas := runOverheadExperiment(10, []float64{1, 2, 5}, 8)
	writeToCsv(overheadDatas, DATA_DIR + "/go-art-overhead.csv")

	ratioWriteVerifDatas := runRatioWriteVerifExperiment(100, []float64{1,2,5}, 9)
	writeToCsv(ratioWriteVerifDatas, DATA_DIR + "/go-art-ratio-write-verif.csv")

	space256byteActionData := runSpaceExperiment(100, 2000, 100, 256)
	writeToCsv(space256byteActionData, DATA_DIR + "/got-art-space-256b.csv")
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func runGenerationExperiment(startingActionNb int, endingActionNb int, step int, actionSize int) [][]string {
	var category = fmt.Sprintf("%d-byte actions", actionSize)
	var data = [][]string{{"Length", category}}

	peerKey, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println(err)
	}
	groupKey := make([]byte, 32)
	rand.Read(groupKey)

	var action = generator.GenerateRandomAction(actionSize)

	for actionNb := startingActionNb; actionNb <= endingActionNb; actionNb += step {
		fmt.Printf("Generation for %d %d-byte actions\n", actionNb, actionSize)

		seq := pas.NewPeerActionSequence()
		var totalTime int64 = 0
		for j := 0; j < NUM_REPETITIONS; j++ {
			start := time.Now()
			for i := 0; i < actionNb; i++ {
				seq.Append(action, peerKey, GROUP, groupKey)
			}
			duration := time.Since(start)
			totalTime += duration.Milliseconds()
		}
		avgTime := totalTime / NUM_REPETITIONS
		data = append(data, []string{strconv.Itoa(actionNb), fmt.Sprintf("%v", avgTime)} )
	}

	return data
}

func runVerificationExperiment(startingActionNb int, endingActionNb int, step int, actionSize int) [][]string {
	var category = fmt.Sprintf("%d-byte actions", actionSize)
	var data = [][]string{{"Length", category}}

	peerKey, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println(err)
	}
	groupKey := make([]byte, 32)
	rand.Read(groupKey)

	var action = generator.GenerateRandomAction(actionSize)

	for actionNb := startingActionNb; actionNb <= endingActionNb; actionNb += step {
		fmt.Printf("Verification for %d %d-byte actions\n", actionNb, actionSize)

		seq := pas.NewPeerActionSequence()
		for i := 0; i < actionNb; i++ {
			seq.Append(action, peerKey, GROUP, groupKey)
		}

		var totalTime int64 = 0
		for j := 0; j < NUM_REPETITIONS; j++ {
			start := time.Now()
			seq.Verify()
			duration := time.Since(start)
			totalTime += duration.Milliseconds()
		}

		avgTime := totalTime / NUM_REPETITIONS
		data = append(data, []string{strconv.Itoa(actionNb), fmt.Sprintf("%v", avgTime)} )
	}

	return data
}

func runSpaceExperiment(startingActionNb int, endingActionNb int, step int, actionSize int) [][]string {
	var category = fmt.Sprintf("%d-byte actions", actionSize)
	var data = [][]string{{"Size (bytes)", category}}

	peerKey, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println(err)
	}
	groupKey := make([]byte, 32)
	rand.Read(groupKey)

	var action = generator.GenerateRandomAction(actionSize)

	for actionNb := startingActionNb; actionNb <= endingActionNb; actionNb += step {
		fmt.Printf("Space requirements for %d %d-byte actions\n", actionNb, actionSize)

		seq := pas.NewPeerActionSequence()
		for i := 0; i < actionNb; i++ {
			seq.Append(action, peerKey, GROUP, groupKey)
		}

		serializedHistory := seq.Serialize()
		size := len([]byte(serializedHistory))

		data = append(data, []string{strconv.Itoa(actionNb), fmt.Sprintf("%v", size)} )
	}

	return data
}

func runOverheadExperiment(actionNb int, units []float64, maxLogFactor int) [][]string {
	var data = [][]string{{"Action size (bytes)", "Overhead"}}

	peerKey, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println(err)
	}
	groupKey := make([]byte, 32)
	rand.Read(groupKey)

	for actionSizeLogFactor := 0; actionSizeLogFactor <= maxLogFactor; actionSizeLogFactor++ {
		for _, actionSizeUnit := range units {
			if actionSizeLogFactor == maxLogFactor {
				if actionSizeUnit != 1 {
					break
				}
			}

			actionSize := int(actionSizeUnit * math.Pow10(actionSizeLogFactor))

			fmt.Printf("Space overhead for %d %d-byte actions\n", actionNb, actionSize)
			var action = generator.GenerateRandomAction(actionSize)

			seq := pas.NewPeerActionSequence()
			for i := 0; i < actionNb; i++ {
				seq.Append(action, peerKey, GROUP, groupKey)
			}

			serializedHistory := seq.Serialize()
			size := len([]byte(serializedHistory))
			baseSize := actionSize * actionNb
			overhead :=  float64(size - baseSize) / float64(baseSize)

			data = append(data, []string{strconv.Itoa(actionSize), fmt.Sprintf("%v", overhead)})
		}


	}

	return data
}

func runRatioWriteVerifExperiment(actionNb int, units []float64, maxLogFactor int) [][]string {
	var data = [][]string{{"Action size (bytes)", "Write time", "Verif time", "Ratio write-verif"}}

	peerKey, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println(err)
	}
	groupKey := make([]byte, 32)
	rand.Read(groupKey)

	for actionSizeLogFactor := 0; actionSizeLogFactor <= maxLogFactor; actionSizeLogFactor++ {
		for _, actionSizeUnit := range units {
			if actionSizeLogFactor == maxLogFactor {
				if actionSizeUnit != 1 {
					break
				}
			}

			actualActionNb := 0
			if actionSizeLogFactor < 1 {
				actualActionNb = 1000
			} else if actionSizeLogFactor < 5 {
				actualActionNb = 500
			} else if actionSizeLogFactor < 6 {
				actualActionNb = 100
			} else if actionSizeLogFactor < 8 {
				actualActionNb = 10
			} else {
				actualActionNb = 2
			}
			
			actionSize := int(actionSizeUnit * math.Pow10(actionSizeLogFactor))

			
			
			fmt.Printf("Ratio Write/Verif for %d %d-byte actions\n", actualActionNb, actionSize)
			var action = generator.GenerateRandomAction(actionSize)

			var writeTotalTime int64 = 0
			var verifTotalTime int64 = 0

			for j := 0; j < NUM_REPETITIONS; j++ {
				seq := pas.NewPeerActionSequence()
				start := time.Now()
				for i := 0; i < actualActionNb; i++ {
					seq.Append(action, peerKey, GROUP, groupKey)
				}
				writeDur := time.Since(start)
				writeTotalTime += writeDur.Nanoseconds()

				start = time.Now()
				seq.Verify()
				verifDur := time.Since(start)
				verifTotalTime += verifDur.Nanoseconds()
			}

			writeTime := float64(writeTotalTime) / float64(NUM_REPETITIONS)
			verifTime := float64(verifTotalTime) / float64(NUM_REPETITIONS)

			ratio := writeTime / verifTime
			data = append(data, []string{
				strconv.Itoa(actionSize),
				fmt.Sprintf("%v", writeTime),
				fmt.Sprintf("%v", verifTime),
				fmt.Sprintf("%v", ratio)})
		}
	}

	return data
}


func joinExperimentsData(datas... [][]string) [][]string {
	var joinedDatas [][]string

	for i, data := range datas {
		for j, values := range data {
			if i == 0 {
				joinedDatas = append(joinedDatas, values)
			} else {
				joinedDatas[j] = append(joinedDatas[j], values[1:]...)
			}
		}
	}
	return joinedDatas
}

func writeToCsv(data [][]string, fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal("error creating file:", err)
	}

	w := csv.NewWriter(f)
	for _, record := range data {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	w.Flush()
}
