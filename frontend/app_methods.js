var url = "http://localhost:8000"

import getWebsocket from "./ws_obj.js"

function loginRequest() {
    let loginForm = {
        login: this.login,
        password: this.password
    }
    let vueContext = this
    let loginData = JSON.stringify(loginForm)

    fetch(url + "/login", {
        method: 'POST',
        headers: {
            'Content-Type': 'text/plain',
        },
        body: loginData
    }).then((response) => {
        return response.json()
    }).then((data) => {
        console.log(data)

        if (data.user_id != undefined && data.token != undefined) {
            this.showLogForm = false
            this.showMessagesBlock = true
            this.user_id = data.user_id
            this.token = data.token
            getWebsocket(this.user_id, this.token, vueContext)
        }

    }).catch((err) => {
        console.log(err)
    });
}

function sendMessage() {
    if (this.wsConn == null) {
        console.log("ws connection not defined")
        return
    }

    let clientMessage = {
        to_user_id: parseInt(this.selected_client),
        from_user_id: parseInt(this.user_id),
        message_text: this.message
    }

    this.wsConn.send(JSON.stringify(clientMessage))
    this.message = ""
}

export {
    loginRequest,
    sendMessage
}