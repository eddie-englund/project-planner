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

test.describe('Projects page', () => {
  test('renders project cards from API', async ({ page }) => {
    await page.route(`${API}/projects`, (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([fakeProject('p1', 'Alpha'), fakeProject('p2', 'Beta')]),
      })
    })

    await page.goto('/')
    await expect(page.locator('h1')).toHaveText('Projects')
    await expect(page.getByText('Alpha')).toBeVisible()
    await expect(page.getByText('Beta')).toBeVisible()
  })

  test('shows empty state when no projects', async ({ page }) => {
    await page.route(`${API}/projects`, (route) => {
      route.fulfill({ status: 200, contentType: 'application/json', body: '[]' })
    })

    await page.goto('/')
    await expect(page.getByText('No projects yet')).toBeVisible()
  })

  test('search input filters visible cards', async ({ page }) => {
    await page.route(`${API}/projects`, (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([fakeProject('p1', 'Alpha Project'), fakeProject('p2', 'Beta Project')]),
      })
    })

    await page.goto('/')
    await page.getByPlaceholder('Filter projects…').fill('alpha')
    await expect(page.getByText('Alpha Project')).toBeVisible()
    await expect(page.getByText('Beta Project')).not.toBeVisible()
  })

  test('navigates to project detail on card click', async ({ page }) => {
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

    await page.goto('/')
    await page.getByText('Alpha').first().click()
    await expect(page).toHaveURL(/\/projects\/p1/)
  })
})
