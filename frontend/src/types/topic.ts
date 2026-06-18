export interface Topic {
  id: string
  projectId: string
  index: number
  title: string
  color: string
  imageUrl: string | null
  createdAt: string
}

export interface CreateTopicPayload {
  title: string
  color: string
  index: number
}

export interface UpdateTopicPayload {
  title: string
  color: string
  image_url?: string | null
}
