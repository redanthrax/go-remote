<script lang="ts">

const pc = new RTCPeerConnection({
  iceServers: [
    {
      urls: 'stun:stun.l.google.com:19302'
    }
  ]
})

//offer data channel
const sendChannel = pc.createDataChannel('agent')
sendChannel.onclose = () => console.log('sendChannel has closed')
sendChannel.onopen = () => console.log('sendChannel has opened')
sendChannel.onmessage = e => {
    console.log(e.data)
}

pc.addTransceiver('video', {
  direction: 'sendrecv'
})
pc.addTransceiver('audio', {
  direction: 'sendrecv'
})


pc.ontrack = function (event) {
  const el = document.createElement(event.track.kind)
  el.srcObject = event.streams[0]
  el.autoplay = true
  el.controls = true

  document.getElementById('view').appendChild(el)
}

pc.oniceconnectionstatechange = e => console.log(pc.iceConnectionState)

pc.onnegotiationneeded = e => pc.createOffer().then(d => pc.setLocalDescription(d)).catch((err) => console.log(err))

export default {
    data() {
        return {
            agent: {},
            inter: {}
        }
    },
    mounted() {
        this.initiateConnection()
    },
    methods: {
        initiateConnection() {
            this.agent.ID = this.$route.query.id
            console.log("Initiating connection to agent: " + this.agent.ID)
            pc.onicecandidate = event => {
              if (event.candidate === null) {
                let localDesc = btoa(JSON.stringify(pc.localDescription))
                const requestOptions = {
                    method: 'POST',
                    body: localDesc
                };

                fetch('http://localhost:8080/sdp?agent=' + this.agent.ID, requestOptions)
                this.inter = setInterval(() => this.checkForAgent(), 2000)
              }
            }
        },
        checkForAgent() {
            console.log("Checking for agent")
            fetch('http://localhost:8080/sdp?agent=' + this.agent.ID)
                .then(response => response.json())
                .then(json => {
                    if(json.AccessDescription.sdp) {
                        this.agent = json
                        clearInterval(this.inter)
                        pc.setRemoteDescription(this.agent.AccessDescription)
                    }
                })
        }
    }
}
</script>

<template>
  <div class="agent">
    <h1>Agent: {{agent.ID}}</h1>
    <div id="view" class="view"></div>
  </div>
</template>

<style>
.agent {
    min-height: 100vh;
    align-items: center;
}
.view {
    width: 800px;
    height: 600px;
    background-color: #FFF;
}
</style>
