import { loginRequest, sendMessage } from "./app_methods.js"

export default {
    data() {
        return {
            login: "",
            password: "",
            token: "",
            user_id: "",
            message: "",
            showLogForm: true,
            showMessagesBlock: false,
            messages: [],
            clients: {},
            selected_client: 0,
            wsConn: null
        }
    },
    methods: {
        loginRequest,
        sendMessage,
        setClient(id) {
            this.selected_client = id
        }
    }
}