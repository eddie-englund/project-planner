export interface Project {
  id: string
  title: string
  color: string
  imageUrl: string | null
  createdBy: string
  createdAt: string
}

export interface CreateProjectPayload {
  title: string
  color: string
}

export interface UpdateProjectPayload {
  title: string
  color: string
  image_url?: string | null
}
