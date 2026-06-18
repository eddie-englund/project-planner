export interface Status {
  id: string
  projectId: string
  name: string
  color: string
  position: number
  isTerminal: boolean
}

export interface Ticket {
  id: string
  topicId: string
  statusId: string | null
  title: string
  body: string
  urls: string[]
  createdAt: string
}

export interface TicketWithTopic extends Ticket {
  topicColor: string
  topicTitle: string
}

export interface CreateTicketPayload {
  title: string
  body: string
  urls: string[]
  statusId?: string
}
