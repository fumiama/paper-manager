export enum MessageTypeEnum {
  MessageNormal = 0,
  MessageRegister = 1,
  MessageUserAdded = 2,
  MessageContactChange = 3,
  MessagePasswordChange = 4,
  MessageResetPassword = 5,
}

export interface MessageItem {
  id: number
  avatar: string
  date: string
  text: string
  type: MessageTypeEnum
}
