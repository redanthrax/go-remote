package main

import (
    "io"
    "log"
    "net/http"
    "strconv"
    "encoding/json"
	"github.com/pion/webrtc/v3"
    "github.com/redanthrax/go-remote/server/signal"
)

func enableCors(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
}

type Agent struct {
    ID int
    Ready bool
    RequestDescription webrtc.SessionDescription
    AccessDescription webrtc.SessionDescription
}

func main() {
    log.Println("Starting API server...")

    agents := []Agent{}

    http.HandleFunc("/agent", func(w http.ResponseWriter, r *http.Request) {
        enableCors(&w)
        switch r.Method {
            case "POST":
                log.Println("Received agent post")
                defer r.Body.Close()
                body, _:= io.ReadAll(r.Body)
                agent := Agent{}
                json.Unmarshal(body, &agent)
                if len(agents) == 0 {
                    agents = append(agents, agent)
                } else {
                    hasAgent := false
                    for i := range(agents) {
                        if agents[i].ID == agent.ID {
                            hasAgent = true
                            agents[i] = agent
                        }
                    }

                    if !hasAgent {
                        agents = append(agents, agent)
                    }
                }
            case "GET":
                log.Println("Getting agents")
                data, _ := json.Marshal(agents)
                io.WriteString(w, string(data))
        }
    })

    http.HandleFunc("/sdp", func(w http.ResponseWriter, r *http.Request) {
        enableCors(&w)
        switch r.Method {
            case "POST":
                agentId, _ := strconv.Atoi(r.URL.Query().Get("agent"))
                log.Println("Received sdp post for agent", agentId)
                defer r.Body.Close()
                body, _ := io.ReadAll(r.Body)
                description := webrtc.SessionDescription{}
                signal.Decode(string(body), &description)
                for i := range(agents) {
                    if agents[i].ID == agentId {
                        if description.Type == webrtc.SDPTypeOffer {
                            agents[i].RequestDescription = description
                        } else {
                            agents[i].AccessDescription = description
                        }
                    }                        
                }

            case "GET":
                agentId, _ := strconv.Atoi(r.URL.Query().Get("agent"))
                for i := range(agents) {
                    if agents[i].ID == agentId {
                        jsonData, _ := json.Marshal(agents[i])
                        io.WriteString(w, string(jsonData))
                    }
                }
        }
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
