export enum MessageType {
  alert = 'alert',
  error = 'error'
}

export interface Message {
  type: MessageType
  message: string
}
