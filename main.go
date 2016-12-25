package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	http.HandleFunc("/hello", hello)

	http.HandleFunc("/cgi-bin/", func(w http.ResponseWriter, r *http.Request) {
		cmdName := strings.SplitN(r.URL.Path, "/", 3)[2]
		var (
			cmdOut []byte
			err    error
		)
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		cmdArgs := []string{}
		if cmdOut, err = exec.Command(filepath.Join(dir, cmdName), cmdArgs...).Output(); err != nil {
			fmt.Fprintln(os.Stderr, "There was an error running git rev-parse command: ", err)
			os.Exit(1)
		}

		w.Write([]byte(cmdOut))

	})

	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]

		data, err := query(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})

	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
}

func query(city string) (weatherData, error) {
	appid := os.Getenv("APPID")
	api_url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?APPID=%s&q=%s", appid, city)
	fmt.Println(api_url)
	resp, err := http.Get(api_url)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var d weatherData

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}

	return d, nil
}

type commandData struct {
	Name string `json:"name"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}
