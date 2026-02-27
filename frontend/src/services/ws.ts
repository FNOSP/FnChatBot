import { useAuthStore } from '../store/auth'

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

    let wsUrl = `${this.url}/ws/chat/${conversationId}`
    // Attach JWT token as query parameter when available.
    const auth = useAuthStore()
    if (auth?.token) {
      const sep = wsUrl.includes('?') ? '&' : '?'
      wsUrl = `${wsUrl}${sep}token=${encodeURIComponent(auth.token)}`
    }

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
