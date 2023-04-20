export enum MessageTypeEnum {
  MessageNormal = 0,
  MessageRegister = 1,
  MessageUserAdded = 2,
  MessageContactChange = 3,
  MessagePasswordChange = 4,
  MessageResetPassword = 5,
  MessageOperator = 6,
}

export interface MessageItem {
  id: number
  avatar: string
  date: string
  text: string
  type: MessageTypeEnum
}

export interface UserRegex {
  ID: number
  Title: string
  Class: string
  OpenCl: string
  Date: string
  Time: string
  Rate: string
  Major: string
  Sub: string
}
