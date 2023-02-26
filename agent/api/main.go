package api

import (
    "context"
    "log"
    "io"
    "time"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/redanthrax/go-remote/agent/structs"
)

func WaitForConnectionPromise(agent *structs.Agent)(waitComplete <-chan struct{}) {
    waitingComplete, done := context.WithCancel(context.Background())

    for {
        log.Println("Waiting for connection...")
        resp, err := http.Get(fmt.Sprintf("http://localhost:8080/sdp?agent=%s", agent.ID))
        if err == nil {
            defer resp.Body.Close()
            body, err2 := io.ReadAll(resp.Body)
            if err2 == nil {
                json.Unmarshal(body, &agent)
                if agent.RequestDescription.SDP != "" {
                    done()
                    return waitingComplete.Done()
                }
            } else {
                log.Println(err2)
            }
        } else {
            log.Println(err)
        }

        time.Sleep(time.Millisecond * 5000)
    }
}
