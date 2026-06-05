import type { Item, User, Area, Notification, LoginResponse } from '@/types'

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? ''
const BASE = `${API_URL}/api/v1`
const AUTH = `${API_URL}/auth`

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(url, {
    credentials: 'include',
    headers: { 'Content-Type': 'application/json', ...options?.headers },
    ...options,
  })

  if (res.status === 401) {
    // Try refresh
    const refreshed = await fetch(`${AUTH}/refresh`, { method: 'POST', credentials: 'include' })
    if (!refreshed.ok) {
      window.location.href = '/login'
      throw new Error('session expired')
    }
    // Retry original
    const retry = await fetch(url, { credentials: 'include', ...options })
    if (!retry.ok) throw new Error(await retry.text())
    return retry.json()
  }

  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

// Auth
export const api = {
  login: (email: string, password: string) =>
    request<LoginResponse>(`${AUTH}/login`, { method: 'POST', body: JSON.stringify({ email, password }) }),

  logout: () => request(`${AUTH}/logout`, { method: 'POST' }),

  me: () => request<User>(`${BASE}/auth/me`),

  // Items
  items: {
    list: () => request<{ items: Item[]; total: number }>(`${BASE}/items`),
    get: (id: string) => request<Item>(`${BASE}/items/${id}`),
    create: (data: Partial<Item>) => request<Item>(`${BASE}/items`, { method: 'POST', body: JSON.stringify(data) }),
    update: (id: string, data: Partial<Item>) => request<Item>(`${BASE}/items/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
    delete: (id: string) => request(`${BASE}/items/${id}`, { method: 'DELETE' }),
    comments: {
      list: (itemId: string) => request<{ comments: Comment[] }>(`${BASE}/items/${itemId}/comments`),
      create: (itemId: string, content: string) =>
        request(`${BASE}/items/${itemId}/comments`, { method: 'POST', body: JSON.stringify({ content }) }),
    },
    interactions: {
      list: (itemId: string) => request(`${BASE}/items/${itemId}/interactions`),
      create: (itemId: string, data: { addressed_to: string; message: string }) =>
        request(`${BASE}/items/${itemId}/interactions`, { method: 'POST', body: JSON.stringify(data) }),
    },
  },

  // Users
  users: {
    list: () => request<{ users: User[]; total: number }>(`${BASE}/users`),
    get: (id: string) => request<User>(`${BASE}/users/${id}`),
  },

  areas: () => request<{ areas: Area[] }>(`${BASE}/areas`),

  notifications: () => request<{ notifications: Notification[] }>(`${BASE}/notifications`),
}
