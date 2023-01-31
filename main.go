package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

var currentAudioPath string
var currentAudioName string

func main() {
	currentAudioName = ""
	currentAudioPath = ""
	r := mux.NewRouter()

	r.PathPrefix("/valid").HandlerFunc(validate)
	r.HandleFunc("/invalid", invalidate)
	r.PathPrefix("/").Handler(http.FileServer(rice.MustFindBox("app").HTTPBox()))

	go filepath.Walk("./unlabelled", playWav)
	fmt.Println("Starting")
	err := http.ListenAndServe(":1235", r)
	if err != nil {
		panic(err)
	}
}

func playWav(path string, file fs.FileInfo, err error) error {
	currentAudioName = file.Name()
	currentAudioPath = path

	if err != nil {
		return err
	}

	if file.IsDir() {
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		panic(err)
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	color.Green("Playing audio...")
	done := make(chan bool)
	audio := buffer.Streamer(0, buffer.Len())
	speaker.Play(beep.Seq(audio, beep.Callback(func() {
		done <- true
	})))
	<-done
	color.Green("Fin.")
	speaker.Close()
	return nil
}

func validate(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Validating: " + currentAudioPath + currentAudioName)
	err := os.Rename(currentAudioPath, fmt.Sprintf("./positive/labelled_%s.wav", currentAudioName))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
}

func invalidate(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Invalidating: " + currentAudioPath + currentAudioName)
	err := os.Rename(currentAudioPath, fmt.Sprintf("./negative/labelled_%s.wav", currentAudioName))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
}
