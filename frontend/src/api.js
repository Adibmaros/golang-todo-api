const TOKEN_KEY = 'todo_jwt_token'
const EMAIL_KEY = 'todo_user_email'

export function getToken() {
  return localStorage.getItem(TOKEN_KEY) || ''
}

export function getUserEmail() {
  return localStorage.getItem(EMAIL_KEY) || ''
}

export function setAuth(token, email) {
  localStorage.setItem(TOKEN_KEY, token)
  localStorage.setItem(EMAIL_KEY, email)
}

export function clearAuth() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(EMAIL_KEY)
}

async function request(path, options = {}, onUnauthorized) {
  const token = getToken()
  const headers = {
    'Content-Type': 'application/json',
    ...(options.headers || {}),
  }

  if (token) {
    headers.Authorization = `Bearer ${token}`
  }

  const res = await fetch(`/api${path}`, {
    ...options,
    headers,
  })

  if (res.status === 401) {
    clearAuth()
    if (onUnauthorized) onUnauthorized()
  }

  if (res.status === 204) {
    return null
  }

  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    const fieldMsg = data.fields ? Object.values(data.fields).join(', ') : ''
    throw new Error(fieldMsg || data.error || `Request gagal (${res.status})`)
  }

  return data
}

export const api = {
  health: () => request('/health'),
  register: (email, password) =>
    request('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),
  login: (email, password) =>
    request('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),
  listTodos: ({ page = 1, limit = 10, done = '', sort = 'created_at' }, onUnauthorized) => {
    const params = new URLSearchParams({
      page: String(page),
      limit: String(limit),
      sort,
    })
    if (done !== '') {
      params.set('done', done)
    }
    return request(`/todos?${params.toString()}`, {}, onUnauthorized)
  },
  createTodo: (title, onUnauthorized) =>
    request(
      '/todos',
      {
        method: 'POST',
        body: JSON.stringify({ title }),
      },
      onUnauthorized
    ),
  updateTodo: (id, payload, onUnauthorized) =>
    request(
      `/todos/${id}`,
      {
        method: 'PUT',
        body: JSON.stringify(payload),
      },
      onUnauthorized
    ),
  deleteTodo: (id, onUnauthorized) =>
    request(`/todos/${id}`, { method: 'DELETE' }, onUnauthorized),
}
