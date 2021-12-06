package service

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	animeI "github.com/jcaladwizeline/academy-go-q42021/api/model"

	"github.com/fatih/structs"
)

func GetAllAnimes() []animeI.AnimeStruct {
	// open csv
	f, err := os.Open("test.csv")
	if err != nil {
		log.Println("Unable to read test.csv", err)
		return nil
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println("Unable to close test.csv", err)
		}
	}(f)

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil
	}
	var response = make([]animeI.AnimeStruct, len(records))
	for row, content := range records {

		animeID, err := strconv.Atoi(content[0])
		if err != nil {
			return nil
		}

		singleRow := animeI.AnimeStruct{
			AnimeID:  animeID,
			Title:    content[1],
			Synopsis: content[2],
			Studio:   content[3],
		}

		response[row] = singleRow
	}
	return response
}

func GetAnimeById(id string) animeI.AnimeStruct {
	var s animeI.AnimeStruct
	// check params
	var idValue int
	if id != "" {
		row, err := strconv.Atoi(id)
		if err != nil {
			return s
		}
		idValue = row
	}

	// open csv
	f, err := os.Open("test.csv")
	if err != nil {
		log.Println("Unable to read test.csv", err)
		return s
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println("Unable to close test.csv", err)
		}
	}(f)

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return s
	}
	var newRecord []string
	for i := 0; i < len(records); i++ {
		value, _ := strconv.Atoi(records[i][0])
		if value == idValue {
			newRecord = records[i]
			break
		}
	}

	if id != "" && len(newRecord) == 0 {
		log.Println("Record does not exists")
		return s
	}
	if len(newRecord) > 1 {
		x := [][]string{newRecord}
		records = x
	}

	var response = make([]animeI.AnimeStruct, len(records))
	for row, content := range records {

		animeID, err := strconv.Atoi(content[0])
		if err != nil {
			return s
		}

		response[row] = animeI.AnimeStruct{
			AnimeID:  animeID,
			Title:    content[1],
			Synopsis: content[2],
			Studio:   content[3],
		}
	}
	return response[0]
}

func PostAnimeById(id string) int {
	animeData := animeByIDExternalAPI(id)
	animeValues := make([]string, 0)
	for _, v := range structs.Values(animeData) {
		animeValues = append(animeValues, fmt.Sprint(v))
	}
	// open csv
	f, err := os.OpenFile("test.csv", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Unable to open test.csv", err)
		return http.StatusInternalServerError
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println("Unable to close test.csv", err)
		}
	}(f)

	csvwriter := csv.NewWriter(f)
	defer csvwriter.Flush()

	if err := csvwriter.Write(animeValues); err != nil {
		log.Fatalln("error writing record to file", err)
		return http.StatusInternalServerError
	}

	return http.StatusAccepted
}

func animeByIDExternalAPI(id string) animeI.AnimeStruct {
	url := "https://api.jikan.moe/v3/anime/" + id
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	animeID, _ := result["mal_id"].(float64)
	title := strings.Replace(result["title"].(string), ",", "", -1)
	synopsis := strings.Replace(result["synopsis"].(string), ",", "", -1)
	studio := strings.Replace(result["studios"].([]interface{})[0].(map[string]interface{})["name"].(string), ",", "", -1)
	animeData := animeI.AnimeStruct{
		AnimeID:  int(animeID),
		Title:    title,
		Synopsis: synopsis,
		Studio:   studio,
	}
	return animeData
}

func worker(t string, jobs <-chan []string, results chan<- animeI.AnimeStruct) {
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			animeID, err := strconv.Atoi(job[0])
			if err != nil {
				return
			}

			anime := animeI.AnimeStruct{
				AnimeID:  animeID,
				Title:    job[1],
				Synopsis: job[2],
				Studio:   job[3],
			}
			if t == "odd" && anime.AnimeID%2 == 0 {
				results <- anime
			} else if t == "even" && anime.AnimeID%2 != 0 {
				results <- anime
			}
		}
	}
}

func WorkerPool(numJobs int, itemsPerWorker int, jobType string) ([]animeI.AnimeStruct, error) {
	// open csv
	f, err := os.Open("test.csv")
	if err != nil {
		log.Println("Unable to read test.csv", err)
		return nil, err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println("Unable to close test.csv", err)
		}
	}(f)

	csvReader := csv.NewReader(f)

	animes := make([]animeI.AnimeStruct, 0)
	jobs := make(chan []string, itemsPerWorker)
	result := make(chan animeI.AnimeStruct, numJobs)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		worker(jobType, jobs, result)
	}()

	for j := 1; j <= numJobs; j++ {
		rStr, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		jobs <- rStr
	}

	close(jobs)
	wg.Wait()
	close(result)

	for a := range result {
		animes = append(animes, a)
	}

	return animes, nil
}
