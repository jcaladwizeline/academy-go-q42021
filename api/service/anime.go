package service

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/jcaladwizeline/academy-go-q42021/api/model"

	"github.com/fatih/structs"
	"github.com/spf13/viper"
)

func GetAllAnimes() ([]model.Anime, int, error) {
	// open csv
	f, err := os.Open(viper.GetString("Files.Name"))
	if err != nil {
		log.Println("Unable to read test.csv", err)

		return nil, http.StatusInternalServerError, errors.New("unable to read test.csv")
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
		return nil, http.StatusInternalServerError, err
	}
	var response = make([]model.Anime, len(records))
	for row, content := range records {
		animeID, err := strconv.Atoi(content[0])
		if err != nil {
			log.Println("Unable to convert animeID from string into integer", err)

			return nil, http.StatusInternalServerError, errors.New("unable to convert animeID from string into integer")
		}

		response[row] = model.Anime{
			AnimeID:  animeID,
			Title:    content[1],
			Synopsis: content[2],
			Studio:   content[3],
		}
	}

	return response, http.StatusAccepted, nil
}

func GetAnimeById(id int) (model.Anime, error) {
	var s model.Anime

	// open csv
	f, err := os.Open(viper.GetString("Files.Name"))
	if err != nil {
		log.Println("Unable to read test.csv", err)

		return s, err
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

		return s, err
	}
	var newRecord []string
	for i := 0; i < len(records); i++ {
		value, _ := strconv.Atoi(records[i][0])
		if value == id {
			newRecord = records[i]

			break
		}
	}

	if len(newRecord) == 0 {
		log.Println("Record does not exists")

		return s, errors.New("record does not exists")
	}

	animeID, err := strconv.Atoi(newRecord[0])
	if err != nil {
		return s, err
	}

	return model.Anime{
		AnimeID:  animeID,
		Title:    newRecord[1],
		Synopsis: newRecord[2],
		Studio:   newRecord[3],
	}, nil
}

func PostAnimeById(id string) (model.Anime, error) {
	var s model.Anime
	animeData, err := animeByIDExternalAPI(id)
	if err != nil {
		return s, err
	}
	animeValues := make([]string, 1)
	for _, v := range structs.Values(animeData) {
		animeValues = append(animeValues, fmt.Sprint(v))
	}
	// open csv
	f, err := os.OpenFile(viper.GetString("Files.Name"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Unable to open test.csv", err)

		return s, err
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

		return s, err
	}

	return animeData, nil
}

func animeByIDExternalAPI(id string) (model.Anime, error) {
	var s model.Anime
	url := viper.GetString("ExternalApis.Url") + id
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)

		return s, err
	}
	if resp.StatusCode == 404 {
		log.Println("anime not found")

		return s, errors.New("anime not found")
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	animeID := result["mal_id"].(float64)
	title := strings.Replace(result["title"].(string), ",", "", -1)
	synopsis := strings.Replace(result["synopsis"].(string), ",", "", -1)
	studio := strings.Replace(result["studios"].([]interface{})[0].(map[string]interface{})["name"].(string), ",", "", -1)

	return model.Anime{
		AnimeID:  int(animeID),
		Title:    title,
		Synopsis: synopsis,
		Studio:   studio,
	}, err
}

func worker(t string, jobs <-chan []string, results chan<- model.Anime) {
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

			anime := model.Anime{
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

func WorkerPool(numJobs int, itemsPerWorker int, jobType string) ([]model.Anime, error) {
	// open csv
	f, err := os.Open(viper.GetString("Files.Name"))
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

	animes := make([]model.Anime, 0)
	jobs := make(chan []string, itemsPerWorker)
	result := make(chan model.Anime, numJobs)

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
