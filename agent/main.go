//https://gist.github.com/thesubtlety/be6e7ec9c19083473bed4cae11c8160d
//https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getcursorpos

package main

import (
    "log"
    "time"
    //"context"
    "net/http"
    "bytes"
    "io"
    "math/rand"

	"github.com/pion/webrtc/v3"
    "github.com/redanthrax/go-remote/agent/signal"
)

func main() {
    log.Println("Starting agent service")

    config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
    
    // Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if cErr := peerConnection.Close(); cErr != nil {
			log.Printf("cannot close peerConnection: %v\n", cErr)
		}
	}()

    peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
        log.Println("On ice candidate")
		if c == nil {
			return
		}
	})

    peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		log.Printf("Peer Connection State has changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			log.Println("Peer Connection has gone to failed exiting")
            return
		}
	})

    peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		log.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			log.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", d.Label(), d.ID())

			for range time.NewTicker(5 * time.Millisecond).C {
				// Send the message as text
				sendErr := d.SendText(RandomString(5))
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

    log.Println("Getting offers from API")
    //get the offer from the api
    resp, _ := http.Get("http://localhost:8080/sdp?type=offer")
    body, _ := io.ReadAll(resp.Body)
    resp.Body.Close()
    description := webrtc.SessionDescription{}
    signal.Decode(string(body), &description)
    log.Println("Setting remote description")
    err = peerConnection.SetRemoteDescription(description)
    if err != nil {
        panic(err)
    }

    log.Println("Creating answer")
    // Create an answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

    //data, _ := json.Marshal(answer)
    encoded := signal.Encode(answer)

    _, er := http.Post("http://localhost:8080/sdp", "application/json", bytes.NewBuffer([]byte(encoded)))

    if er != nil {
        log.Fatal(er)
    }

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete
    log.Println("Gathering complete")

    log.Println(*peerConnection.LocalDescription())

    select { }
}

func RandomString(n int) string {
    var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
 
    s := make([]rune, n)
    for i := range s {
        s[i] = letters[rand.Intn(len(letters))]
    }
    return string(s)
}
