class WebSocketStore {
    public connect() {
        const ws = new WebSocket("ws://localhost:8080/ws")
        ws.addEventListener("open", () => console.log("Websocket connected"))
        ws.addEventListener("error", () => console.error("Websocket Error"))
        ws.addEventListener("message", (message: any) => {
            console.log(message)
        })
    }
}

export const webSocketStore = new WebSocketStore()