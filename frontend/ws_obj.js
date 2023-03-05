
const wsUrl = "ws://localhost:8000/ws"

function isJson(data) {
    try {
        JSON.parse(data)
    } catch (e) {
        return false
    }
    return true
}

function getWebsocket(user_id, token, vueContext) {
    let ws = new WebSocket(wsUrl + `?user_id=${user_id}&token=${token}`)
    ws.onopen = function (e) {
        console.log(e)
        vueContext.wsConn = ws
    }
    ws.onclose = function (e) {
        console.log("ws connection closed")
        ws = null;
        vueContext.wsConn = null
    }
    ws.onmessage = function (e) {
        let message = e.data
        console.log(message)

        if (isJson(message)) {
            let obj = JSON.parse(message)

            if (obj.online_clients != undefined) {
                vueContext.clients = obj.online_clients
                return
            }
        }

        vueContext.messages.unshift(message)
    }
    ws.onerror = function (e) {
        console.log("ws error: " + e.data)
    }

    return ws
}

export default getWebsocket
