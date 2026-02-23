type MessageHandler = (data: any) => void

export class WebSocketService {
  private ws: WebSocket | null = null
  private url: string
  private onMessageCallback: MessageHandler | null = null

  constructor(url: string) {
    this.url = url
  }

  connect(conversationId: string) {
    if (this.ws) {
      this.ws.close()
    }

    const wsUrl = `${this.url}/ws/chat/${conversationId}`
    console.log('Connecting to WebSocket:', wsUrl)
    this.ws = new WebSocket(wsUrl)

    this.ws.onopen = () => {
      console.log('WebSocket connected')
    }

    this.ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        if (this.onMessageCallback) {
          this.onMessageCallback(data)
        }
      } catch (e) {
        console.error('Error parsing WebSocket message:', e)
      }
    }

    this.ws.onclose = () => {
      console.log('WebSocket disconnected')
    }

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }
  }

  sendMessage(message: any) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message))
    } else {
      console.error('WebSocket is not connected')
    }
  }

  onMessage(callback: MessageHandler) {
    this.onMessageCallback = callback
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }
}
