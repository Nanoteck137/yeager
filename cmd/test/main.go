package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/nanoteck137/yeager/apis"
	"github.com/nanoteck137/yeager/core/log"
)

func JsonString(v any) string {
	d, err := json.Marshal(v)
	if err != nil {
		log.Fatal("Failed to marshal object", "err", err)
	}

	return string(d)
}

func main() {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	type Track struct {
		Title    string
		Filename string
	}

	tracks := []Track{
		{
			Title:    "Highway to Hell",
			Filename: "01.Highway_to_Hell.flac",
		},
		{
			Title:    "Girls Got Rhythm",
			Filename: "02.Girls_Got_Rhythm.flac",
		},
		{
			Title:    "Walk All OverYou",
			Filename: "03.Walk_All_Over_You.flac",
		},
		{
			Title: "Hello World",
			Filename: "04.Touch_Too_Much.flac",
		},
		{
			Title: "Beating Around the Bush",
			Filename: "05.Beating_Around_the_Bush.flac",
		},
		{
			Title: "Shot Down in Flames",
			Filename: "06.Shot_Down_in_Flames.flac",
		},
		{
			Title: "Get It Hot",
			Filename: "07.Get_It_Hot.flac",
		},
		{
			Title: "If You Want Blood (You’ve Got It)",
			Filename: "08.If_You_Want_Blood_(You’ve_Got_It).flac",
		},
		{
			Title: "Love Hungry Man",
			Filename: "09.Love_Hungry_Man.flac",
		},
		{
			Title: "Night Prowler",
			Filename: "10.Night_Prowler.flac",
		},
	}

	data := apis.PostAlbumBody{
		Name:   "Highway to Hell",
		Artist: "Test Artist",
	}

	w.WriteField("data", JsonString(data))

	for _, t := range tracks {
		fw, err := w.CreateFormFile("files", t.Filename)
		if err != nil {
			log.Fatal("Failed to crate form file", "err", err)
		}

		file, err := os.Open("./work/files/" + t.Filename)
		if err != nil {
			log.Fatal("Failed to open file", "err", err)
		}
		defer file.Close()

		_, err = io.Copy(fw, file)
		if err != nil {
			log.Fatal("Failed to copy file file to form", "err", err)
		}
	}

	w.Close()

	req, err := http.NewRequest(http.MethodPost, "http://localhost:3000/api/v1/music/album", &b)
	if err != nil {
		log.Fatal("Failed to create request", "err", err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Failed to send request", "err", err)
	}

	d, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Failed to read response data", "err", err)
	}

	fmt.Printf("string(d): %v\n", string(d))
}
