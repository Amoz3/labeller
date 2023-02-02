package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
var audioFiles []string
var index int

func main() {
	index = 0
	currentAudioName = ""
	currentAudioPath = ""
	r := mux.NewRouter()

	r.PathPrefix("/valid").HandlerFunc(validate)
	r.HandleFunc("/invalid", invalidate)
	r.PathPrefix("/").Handler(http.FileServer(rice.MustFindBox("app").HTTPBox()))

	filepath.Walk("./unlabelled", playWav)
	fmt.Println("Starting")
	go play()
	err := http.ListenAndServe(":1235", r)
	if err != nil {
		panic(err)
	}
}

func play() {
	fmt.Println(len(audioFiles))
	for index < len(audioFiles) {
		fmt.Println(fmt.Sprintf("Index %d", index))
		f, err := os.Open(audioFiles[index])
		currentAudioPath = "./" + audioFiles[index]
		currentAudioName = strings.Split(f.Name(), "/")[1]
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
	}
}

func playWav(path string, file fs.FileInfo, err error) error {
	currentAudioPath = path

	if err != nil {
		return err
	}

	if file.IsDir() {
		return nil
	}

	audioFiles = append(audioFiles, path)
	// f, err := os.Open(path)
	// if err != nil {
	// 	panic(err)
	// }

	// streamer, format, err := wav.Decode(f)
	// if err != nil {
	// 	panic(err)
	// }

	// buffer := beep.NewBuffer(format)
	// buffer.Append(streamer)
	// streamer.Close()

	// speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	// color.Green("Playing audio...")
	// done := make(chan bool)
	// audio := buffer.Streamer(0, buffer.Len())
	// speaker.Play(beep.Seq(audio, beep.Callback(func() {
	// 	done <- true
	// })))
	// <-done
	// color.Green("Fin.")
	// speaker.Close()
	return nil
}

func validate(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Validating: " + currentAudioPath + currentAudioName)
	err := os.Rename(currentAudioPath, fmt.Sprintf("./positive/labelled_%s.wav", currentAudioName))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	index++
	w.WriteHeader(http.StatusOK)
}

func invalidate(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Invalidating: " + currentAudioPath + currentAudioName)
	err := os.Rename(currentAudioPath, fmt.Sprintf("./negative/labelled_%s", currentAudioName))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	index++
	w.WriteHeader(http.StatusOK)
}
