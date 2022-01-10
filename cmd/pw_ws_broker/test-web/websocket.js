;(() => {
    console.log("Websockets!")

    function dial() {
        const conn = new WebSocket(`ws://localhost:3000/planes`, ["planes"])
        conn.addEventListener("close", ev => {
            console.error(`WebSocket Disconnected code: ${ev.code}, reason: ${ev.reason}`)
            if (ev.code !== 1001) {
                console.info("Reconnecting in 1s")
                setTimeout(dial, 1000)
            }
        })
        conn.addEventListener("open", ev => {
            console.info("websocket connected")
        })

        // This is where we handle messages received.
        conn.addEventListener("message", ev => {
            if (typeof ev.data !== "string") {
                console.error("unexpected message type", typeof ev.data)
                return
            }

            console.info(ev.data)
        })
    }

    dial()

})()