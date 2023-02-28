//https://gist.github.com/thesubtlety/be6e7ec9c19083473bed4cae11c8160d
//https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getcursorpos

package main

import (
    "log"
    "time"
    //"net/http"
    //"bytes"
    //"encoding/json"
    "io/ioutil"

    "github.com/redanthrax/go-remote/agent/structs"
    //"github.com/redanthrax/go-remote/agent/wrtc"
    //"github.com/redanthrax/go-remote/agent/api"
	//"github.com/pion/webrtc/v3"
    "gopkg.in/yaml.v3"
	//"github.com/pion/webrtc/v3/pkg/media/oggreader"
)

var agent structs.Agent
var state structs.State

func main() {
    //startup
    //get config data
    agent := structs.Agent{}
    log.Println("Checking agent config")
    yfile, err := ioutil.ReadFile("config.yml")    
    if err != nil {
        log.Fatal(err)
    }

    err = yaml.Unmarshal(yfile, &agent)
    if err != nil {
        log.Fatal(err)
    }

    log.Println(agent)

    //start looping to do stuff speed limit at 1 second for non-bog
    for {
        log.Println("Starting agent loop")

                 

        time.Sleep(time.Second * 1)
    }

    /*
    //say hello to the api and register agent - no auth
    log.Println("Connecting to API")
    jsonData, err2 := json.Marshal(agent)
    if err2 != nil {
        log.Println(err2)
    }

    http.Post("http://localhost:8080/agent", "application/json", bytes.NewBuffer(jsonData))
    waitComplete := api.WaitForConnectionPromise(&agent)
    <-waitComplete

    //startup web rtc initialization
    config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

    pc, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if cErr := pc.Close(); cErr != nil {
			log.Printf("cannot close peerConnection: %v\n", cErr)
		}
	}()

    wrtc.InitiateWebRTC(pc)
    log.Println("Setting remote description")
    err = pc.SetRemoteDescription(agent.RequestDescription)
    if err != nil {
        log.Fatal(err)
    }

    agent.AccessDescription = wrtc.CompleteWebrtcConnection(pc)
    jsonData, _ = json.Marshal(agent)
    log.Println(string(jsonData))
    http.Post("http://localhost:8080/agent", "application/json", bytes.NewBuffer(jsonData))
    */
}

