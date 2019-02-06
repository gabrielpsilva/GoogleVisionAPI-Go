package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
)

//ImageTexts are Text Object
type ImageTexts struct {
	Locale      string `json:"locale"`
	Description string `json:"description"`
	Rect        Rect   `json:"bounding_poly"`
}

//Rect are Text Square
type Rect struct {
	Vertices []Vertice `json:"vertices"`
}

//Vertice are Dots of vertice
type Vertice struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

func main() {

	// [START setting_port]
	//_ = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "C:\\Users\\gabri\\Documents\\creds\\mlvision-ocr12.json")
	http.HandleFunc("/up", ReceiveFile)
	http.HandleFunc("/", hello)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
	// [END setting_port]
}

func hello(w http.ResponseWriter, r *http.Request) {

	RespondWithJSON(w, 201, "HELLO WORLD!")

}

//ReceiveFile Receive file from post
func ReceiveFile(w http.ResponseWriter, r *http.Request) {

	var Buf bytes.Buffer

	//file, header, err := r.FormFile("fileupload")
	file, _, err := r.FormFile("fileupload")
	if err != nil {
		fmt.Println(err)
	}
	defer func() { _ = file.Close() }()

	//name := header.Filename
	//secret := r.Header.Get("secret")

	_, _ = io.Copy(&Buf, file)
	defer Buf.Reset()

	ret := SubmitToOcr(&Buf)

	if ret == nil {
		respondWithError(w, http.StatusInternalServerError, "ERROR 01")
		return
	}

	RespondWithJSON(w, http.StatusCreated, ret)

	return
}

//SubmitToOcr Performs OCR
func SubmitToOcr(fileBuffer *bytes.Buffer) []ImageTexts {

	ctx := context.Background()
	fmt.Println(ctx)

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	fmt.Println(ctx, client)

	image, err := vision.NewImageFromReader(fileBuffer)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		log.Fatalf("Failed finding text: %v", err)
	}

	fmt.Println(annotations)
	var allTexts []ImageTexts
	for _, annotation := range annotations {

		var itxt ImageTexts
		var bounding Vertice

		itxt.Locale = annotation.GetLocale()
		itxt.Description = annotation.GetDescription()

		var rect Rect
		for _, bonds := range annotation.GetBoundingPoly().GetVertices() {
			bounding.X = bonds.GetX()
			bounding.Y = bonds.GetY()
			rect.Vertices = append(rect.Vertices, bounding)
		}
		itxt.Rect = rect
		allTexts = append(allTexts, itxt)
	}
	return allTexts
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, map[string]string{"error": msg})
}

//RespondWithJSON Submit response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}
