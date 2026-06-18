import { test, expect } from '@playwright/test'

const API = 'http://localhost:8080'

const fakeProject = (id: string, title: string) => ({
  id,
  title,
  color: '#3b82f6',
  imageUrl: null,
  createdBy: 'user-1',
  createdAt: new Date().toISOString(),
})

const fakeTopic = (id: string, projectId: string, title: string, color = '#10b981') => ({
  id,
  projectId,
  index: 0,
  title,
  color,
  imageUrl: null,
  createdAt: new Date().toISOString(),
})

const fakeStatus = (id: string, name: string, color = '#3b82f6', position = 0) => ({
  id,
  projectId: 'p1',
  name,
  color,
  position,
  isTerminal: false,
})

const fakeTicketWithTopic = (
  id: string,
  topicId: string,
  title: string,
  statusId: string | null = null,
  topicColor = '#10b981',
  topicTitle = 'Backend'
) => ({
  id,
  topicId,
  statusId,
  title,
  body: '',
  urls: [],
  createdAt: new Date().toISOString(),
  topicColor,
  topicTitle,
})

function setupBaseRoutes(page: any, project: any, topics: any[], statuses: any[], tickets: any[]) {
  page.route(`${API}/projects`, (route: any) =>
    route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify([project]) })
  )
  page.route(`${API}/projects/p1`, (route: any) =>
    route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(project) })
  )
  page.route(`${API}/projects/p1/topics`, (route: any) =>
    route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(topics) })
  )
  page.route(`${API}/projects/p1/statuses`, (route: any) =>
    route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(statuses) })
  )
  page.route(`${API}/projects/p1/tickets`, (route: any) =>
    route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(tickets) })
  )
}

test.describe('Project kanban page', () => {
  test('renders kanban columns with tickets from all topics', async ({ page }) => {
    const project = fakeProject('p1', 'Alpha')
    const topics = [
      fakeTopic('t1', 'p1', 'Backend', '#10b981'),
      fakeTopic('t2', 'p1', 'Frontend', '#3b82f6'),
    ]
    const statuses = [fakeStatus('s1', 'Open'), fakeStatus('s2', 'Closed', '#ef4444', 1)]
    const tickets = [
      fakeTicketWithTopic('tk1', 't1', 'Fix bug', 's1', '#10b981', 'Backend'),
      fakeTicketWithTopic('tk2', 't2', 'Design UI', 's1', '#3b82f6', 'Frontend'),
    ]

    setupBaseRoutes(page, project, topics, statuses, tickets)

    await page.goto('/projects/p1')
    await expect(page.getByText('Fix bug')).toBeVisible()
    await expect(page.getByText('Design UI')).toBeVisible()
    await expect(page.getByText('Open')).toBeVisible()
  })

  test('topic filter chips filter tickets by topic', async ({ page }) => {
    const project = fakeProject('p1', 'Alpha')
    const topics = [
      fakeTopic('t1', 'p1', 'Backend', '#10b981'),
      fakeTopic('t2', 'p1', 'Frontend', '#3b82f6'),
    ]
    const statuses = [fakeStatus('s1', 'Open')]
    const tickets = [
      fakeTicketWithTopic('tk1', 't1', 'Fix bug', 's1', '#10b981', 'Backend'),
      fakeTicketWithTopic('tk2', 't2', 'Design UI', 's1', '#3b82f6', 'Frontend'),
    ]

    setupBaseRoutes(page, project, topics, statuses, tickets)

    await page.goto('/projects/p1')
    await expect(page.getByText('Fix bug')).toBeVisible()
    await expect(page.getByText('Design UI')).toBeVisible()

    await page.getByRole('button', { name: 'Backend' }).click()
    await expect(page.getByText('Fix bug')).toBeVisible()
    await expect(page.getByText('Design UI')).not.toBeVisible()
  })

  test('search filters tickets by title', async ({ page }) => {
    const project = fakeProject('p1', 'Alpha')
    const topics = [fakeTopic('t1', 'p1', 'Backend')]
    const statuses = [fakeStatus('s1', 'Open')]
    const tickets = [
      fakeTicketWithTopic('tk1', 't1', 'Fix bug', 's1'),
      fakeTicketWithTopic('tk2', 't1', 'Add feature', 's1'),
    ]

    setupBaseRoutes(page, project, topics, statuses, tickets)

    await page.goto('/projects/p1')
    await page.getByPlaceholder('Search tickets…').fill('fix')
    await expect(page.getByText('Fix bug')).toBeVisible()
    await expect(page.getByText('Add feature')).not.toBeVisible()
  })

  test('topic badge visible on ticket card', async ({ page }) => {
    const project = fakeProject('p1', 'Alpha')
    const topics = [fakeTopic('t1', 'p1', 'Backend')]
    const statuses = [fakeStatus('s1', 'Open')]
    const tickets = [fakeTicketWithTopic('tk1', 't1', 'Fix bug', 's1', '#10b981', 'Backend')]

    setupBaseRoutes(page, project, topics, statuses, tickets)

    await page.goto('/projects/p1')
    await expect(page.getByText('Backend').last()).toBeVisible()
  })

  test('Topics tab navigates to topics grid', async ({ page }) => {
    const project = fakeProject('p1', 'Alpha')
    const topics = [fakeTopic('t1', 'p1', 'Backend')]
    const statuses = [fakeStatus('s1', 'Open')]

    setupBaseRoutes(page, project, topics, statuses, [])

    await page.goto('/projects/p1')
    await page.getByRole('link', { name: 'Topics', exact: true }).first().click()
    await expect(page).toHaveURL(/\/projects\/p1\/topics/)
  })
})
