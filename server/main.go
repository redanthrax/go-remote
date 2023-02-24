package main

import (
    "log"
    "net/http"
    //"fmt"
    "io"
	"github.com/pion/webrtc/v3"
    "github.com/redanthrax/go-remote/server/signal"
)

func enableCors(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
    log.Println("Starting API server...")

    agent := webrtc.SessionDescription{}
    client := webrtc.SessionDescription{}

    http.HandleFunc("/sdp", func(w http.ResponseWriter, r *http.Request) {
        enableCors(&w)
        switch r.Method {
            case "POST":
                log.Println("Received agent post")
                defer r.Body.Close()
                body, _ := io.ReadAll(r.Body)
                description := webrtc.SessionDescription{}
                signal.Decode(string(body), &description)

                if description.Type == webrtc.SDPTypeOffer {
                    client = description
                } else {
                    agent = description
                }

            case "GET":
                rType := r.URL.Query().Get("type")
                log.Println("Get of type", rType)
                if rType == "offer" {
                    io.WriteString(w, signal.Encode(client))
                } else {
                    io.WriteString(w, signal.Encode(agent))
                }
        }
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
