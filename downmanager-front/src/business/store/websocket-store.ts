import { linkStore } from "./link-store";
import type { UpdateProgress } from "../model/link-model";

interface DownloadMessage {
	Linkref:    string
	Finished:   boolean
	InError:    boolean
	ErrorMsg:   string
	Total:      number
	Downloaded: number
}
class WebSocketStore {
    public connect() {
        const ws = new WebSocket("ws://localhost:8080/ws")
        ws.addEventListener("open", () => console.log("Websocket connected"))
        ws.addEventListener("error", () => console.error("Websocket Error"))
        ws.addEventListener("message", (message: any) => {
            console.log(message)
            const content = JSON.parse(message.data) as DownloadMessage
            linkStore.updateProgress(content.Linkref, {
                Downloaded: content.Downloaded,
                Finished: content.Finished,
                InError: content.InError,
            } as UpdateProgress)
        })
    }
}

export const webSocketStore = new WebSocketStore()