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

const fakeTopic = (id: string, projectId: string, title: string, index = 0) => ({
  id,
  projectId,
  index,
  title,
  color: '#10b981',
  imageUrl: null,
  createdAt: new Date().toISOString(),
})

test.describe('Project topics page', () => {
  test('renders topics list for a project', async ({ page }) => {
    const project = fakeProject('p1', 'Alpha')
    const topics = [fakeTopic('t1', 'p1', 'Frontend'), fakeTopic('t2', 'p1', 'Backend', 1)]

    await page.route(`${API}/projects`, (route) => {
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify([project]) })
    })
    await page.route(`${API}/projects/p1`, (route) => {
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(project) })
    })
    await page.route(`${API}/projects/p1/topics`, (route) => {
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(topics) })
    })

    await page.goto('/projects/p1/topics')
    await expect(page.getByText('Alpha')).toBeVisible()
    await expect(page.getByText('Frontend')).toBeVisible()
    await expect(page.getByText('Backend')).toBeVisible()
  })

  test('shows empty state when no topics', async ({ page }) => {
    const project = fakeProject('p1', 'Alpha')

    await page.route(`${API}/projects`, (route) => {
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify([project]) })
    })
    await page.route(`${API}/projects/p1`, (route) => {
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(project) })
    })
    await page.route(`${API}/projects/p1/topics`, (route) => {
      route.fulfill({ status: 200, contentType: 'application/json', body: '[]' })
    })

    await page.goto('/projects/p1/topics')
    await expect(page.getByText('No topics yet')).toBeVisible()
  })

  test('search filters visible topics', async ({ page }) => {
    const project = fakeProject('p1', 'Alpha')
    const topics = [fakeTopic('t1', 'p1', 'Frontend'), fakeTopic('t2', 'p1', 'Backend', 1)]

    await page.route(`${API}/projects`, (route) => {
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify([project]) })
    })
    await page.route(`${API}/projects/p1`, (route) => {
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(project) })
    })
    await page.route(`${API}/projects/p1/topics`, (route) => {
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(topics) })
    })

    await page.goto('/projects/p1/topics')
    await page.getByPlaceholder('Filter topics…').fill('front')
    await expect(page.getByText('Frontend')).toBeVisible()
    await expect(page.getByText('Backend')).not.toBeVisible()
  })
})
