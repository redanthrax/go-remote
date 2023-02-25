//https://gist.github.com/thesubtlety/be6e7ec9c19083473bed4cae11c8160d
//https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getcursorpos

package main

import (
    "log"
    "time"
    "context"
    "net/http"
    "bytes"
    //"io"
    "encoding/json"
    "fmt"
    "io"

	"github.com/pion/webrtc/v3"
    //"github.com/redanthrax/go-remote/agent/signal"
)

type Agent struct {
    ID int
    Ready bool
    RequestDescription webrtc.SessionDescription
    AccessDescription webrtc.SessionDescription
}

var agent Agent

func WaitForConnectionPromise()(waitComplete <-chan struct{}) {
    waitingComplete, done := context.WithCancel(context.Background())

    for {
        log.Println("Waiting for connection...")
        resp, _ := http.Get(fmt.Sprintf("http://localhost:8080/sdp?agent=%d", agent.ID))
        body, _ := io.ReadAll(resp.Body)
        json.Unmarshal(body, &agent)
        if agent.RequestDescription.SDP != "" {
            done()
            return waitingComplete.Done()
        }

        time.Sleep(time.Millisecond * 5000)
    }
}

func InitiateWebRTC(pc *webrtc.PeerConnection) {
    pc.OnICECandidate(func(c *webrtc.ICECandidate) {
        log.Println("On ice candidate")
		if c == nil {
			return
		}
	})

    pc.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		log.Printf("Peer Connection State has changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			log.Println("Peer Connection has gone to failed exiting")
		}
	})

    pc.OnDataChannel(func(d *webrtc.DataChannel) {
		log.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			log.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", d.Label(), d.ID())

			for range time.NewTicker(1 * time.Second).C {
				// Send the message as text
				sendErr := d.SendText("mystring")
				if sendErr != nil {
					panic(sendErr)
				}
			}
		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			log.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
		})
	})

}

func main() {
    log.Println("Starting agent service")
    log.Println("Registering agent")

    agent.ID = 1
    agent.Ready = false

    jsonData, _ := json.Marshal(agent)
    http.Post("http://localhost:8080/agent", "application/json", bytes.NewBuffer(jsonData))

    waitComplete := WaitForConnectionPromise()

    <-waitComplete

    log.Println("Request description obtained")

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

    InitiateWebRTC(pc)

    err = pc.SetRemoteDescription(agent.RequestDescription)
    if err != nil {
        panic(err)
    }

    log.Println("Creating answer")
    // Create an answer
	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

    
	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(pc)

	// Sets the LocalDescription, and starts our UDP listeners
	err = pc.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

    log.Println("Gathering complete")

    agent.AccessDescription = answer

    jsonData, _ = json.Marshal(agent)
    log.Println(string(jsonData))
    http.Post("http://localhost:8080/agent", "application/json", bytes.NewBuffer(jsonData))

    select { }
}

