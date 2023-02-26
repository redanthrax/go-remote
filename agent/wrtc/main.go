package wrtc

import (
    "time"
    "os"
    "log"

	"github.com/pion/webrtc/v3"
    "github.com/pion/webrtc/v3/pkg/media"
	"github.com/pion/webrtc/v3/pkg/media/ivfreader"
)


const (
	audioFileName   = "output.ogg"
	videoFileName   = "output.ivf"
	oggPageDuration = time.Millisecond * 20
)

func InitiateWebRTC(pc *webrtc.PeerConnection) {

    videoTrack, videoTrackErr := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion")
    if videoTrackErr != nil {
        panic(videoTrackErr)
    }

    rtpSender, videoTrackErr := pc.AddTrack(videoTrack)
    if videoTrackErr != nil {
        panic(videoTrackErr)
    }

    go func() {
        rtcpBuf := make([]byte, 1500)
        for {
            if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
                return
            }
        }
    }()

    go func() {
        // Open a IVF file and start reading using our IVFReader
        file, _ := os.Open(videoFileName)
        ivf, header, _ := ivfreader.NewWith(file)

        ticker := time.NewTicker(time.Millisecond * time.Duration((float32(header.TimebaseNumerator)/float32(header.TimebaseDenominator))*1000))
        for ; true; <-ticker.C {
            frame, _, _ := ivf.ParseNextFrame()

            videoTrack.WriteSample(media.Sample{Data: frame, Duration: time.Second})
        }
    }()

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
                    log.Println(sendErr)
				}
			}
		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			log.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
		})
	})
}

func CompleteWebrtcConnection(pc *webrtc.PeerConnection)(webrtc.SessionDescription) {
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

    return answer
}
