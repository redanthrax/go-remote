<script lang="ts">
export default {
    methods: {
        sendMessage() {
            const text = document.getElementById("message").value
            sendChannel.send(text)
        }
    }
}

const log = msg => {
    console.log(msg)
}

const pc = new RTCPeerConnection({
  iceServers: [
    {
      urls: 'stun:stun.l.google.com:19302'
    }
  ]
})

const sendChannel = pc.createDataChannel('agent')
sendChannel.onclose = () => console.log('sendChannel has closed')
sendChannel.onopen = () => console.log('sendChannel has opened')
sendChannel.onmessage = e => {
    let ret = document.getElementById("return")
    ret.append(e.data + " ")
}

pc.oniceconnectionstatechange = e => console.log(pc.iceConnectionState)

pc.onicecandidate = event => {
  console.log("candidate", event.candidate)
  if (event.candidate === null) {
    let localDesc = btoa(JSON.stringify(pc.localDescription))
    const requestOptions = {
        method: 'POST',
        contentType: 'application/json',
        body: localDesc
    };

    fetch('http://localhost:8080/sdp', requestOptions)
  }
}

pc.onnegotiationneeded = e => pc.createOffer().then(d => pc.setLocalDescription(d)).catch(log)

function checkAgents() {
    fetch('http://localhost:8080/sdp?type=answer')
    .then(async response => {
        const data = await response.text()
        const desc = JSON.parse(atob(data))
        console.log(desc)
        if (desc.sdp) {
            try {
                pc.setRemoteDescription(desc)
                clearInterval(checkInterval)
                console.log(pc)
          } catch (e) {
              console.log(e)
          }
        }
    });
}

const checkInterval = setInterval(() => {
    checkAgents();
}, 3000);

setInterval(function() {
  var elem = document.getElementById('return');
  elem.scrollTop = elem.scrollHeight;
}, 100);
</script>

<template>
  <main>
    <div id="agent"></div>
    <textarea cols="50" rows="10" id="message"></textarea><br />
    <button @click="sendMessage">Send</button><br />
    <br />
    <div style="overflow-wrap: everywhere;background-color: #FFF; width: 200px; height: 200px; overflow: auto; color: #000" id="return">
    </div>
  </main>
</template>
