<script setup>
import { onMounted, reactive, ref } from 'vue'
import { api, clearAuth, getToken, getUserEmail, setAuth } from './api'

const token = ref(getToken())
const userEmail = ref(getUserEmail())
const healthStatus = ref('checking')

const authMode = ref('login')
const authForm = reactive({ email: '', password: '' })
const authLoading = ref(false)
const authError = ref('')
const authSuccess = ref('')

const todos = ref([])
const meta = reactive({
  page: 1,
  limit: 10,
  total_data: 0,
  total_pages: 0,
})
const filters = reactive({
  done: '',
  sort: 'created_at',
})

const newTitle = ref('')
const editingId = ref(null)
const editingTitle = ref('')
const listLoading = ref(false)
const actionError = ref('')

function handleUnauthorized() {
  token.value = ''
  userEmail.value = ''
  todos.value = []
}

function formatDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
  })
}

async function checkHealth() {
  try {
    const res = await api.health()
    healthStatus.value = res?.database === 'up' ? 'up' : 'down'
  } catch {
    healthStatus.value = 'down'
  }
}

async function loadTodos() {
  if (!token.value) return
  listLoading.value = true
  actionError.value = ''
  try {
    const res = await api.listTodos(
      {
        page: meta.page,
        limit: meta.limit,
        done: filters.done,
        sort: filters.sort,
      },
      handleUnauthorized
    )
    todos.value = res.data || []
    Object.assign(meta, res.meta || {})
  } catch (err) {
    actionError.value = err.message
  } finally {
    listLoading.value = false
  }
}

async function submitAuth() {
  authLoading.value = true
  authError.value = ''
  authSuccess.value = ''
  try {
    if (authMode.value === 'register') {
      await api.register(authForm.email, authForm.password)
      authSuccess.value = 'Akun berhasil dibuat. Silakan login.'
      authMode.value = 'login'
      authForm.password = ''
    } else {
      const res = await api.login(authForm.email, authForm.password)
      setAuth(res.token, authForm.email)
      token.value = res.token
      userEmail.value = authForm.email
      authForm.password = ''
      meta.page = 1
      await loadTodos()
    }
  } catch (err) {
    authError.value = err.message
  } finally {
    authLoading.value = false
  }
}

function logout() {
  clearAuth()
  handleUnauthorized()
}

async function createTodo() {
  if (!newTitle.value.trim()) return
  actionError.value = ''
  try {
    await api.createTodo(newTitle.value.trim(), handleUnauthorized)
    newTitle.value = ''
    meta.page = 1
    await loadTodos()
  } catch (err) {
    actionError.value = err.message
  }
}

async function toggleDone(todo) {
  actionError.value = ''
  try {
    await api.updateTodo(
      todo.id,
      { title: todo.title, done: !todo.done },
      handleUnauthorized
    )
    await loadTodos()
  } catch (err) {
    actionError.value = err.message
  }
}

function startEdit(todo) {
  editingId.value = todo.id
  editingTitle.value = todo.title
}

function cancelEdit() {
  editingId.value = null
  editingTitle.value = ''
}

async function saveEdit(todo) {
  if (!editingTitle.value.trim()) return
  actionError.value = ''
  try {
    await api.updateTodo(
      todo.id,
      { title: editingTitle.value.trim(), done: todo.done },
      handleUnauthorized
    )
    cancelEdit()
    await loadTodos()
  } catch (err) {
    actionError.value = err.message
  }
}

async function removeTodo(id) {
  actionError.value = ''
  try {
    await api.deleteTodo(id, handleUnauthorized)
    if (todos.value.length === 1 && meta.page > 1) {
      meta.page -= 1
    }
    await loadTodos()
  } catch (err) {
    actionError.value = err.message
  }
}

function setStatusFilter(val) {
  filters.done = val
  meta.page = 1
  loadTodos()
}

function changeSort() {
  meta.page = 1
  loadTodos()
}

function goToPage(nextPage) {
  meta.page = nextPage
  loadTodos()
}

onMounted(() => {
  checkHealth()
  if (token.value) {
    loadTodos()
  }
})
</script>

<template>
  <!-- Auth Layout: Centered, Clean & Focused -->
  <div v-if="!token" class="auth-shell">
    <div class="auth-status">
      <a
        href="https://github.com/Adibmaros/golang-todo-api"
        target="_blank"
        rel="noopener noreferrer"
        class="github-pill"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
          <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12z" />
        </svg>
        GitHub
      </a>
      <span class="status-pill">
        <span class="dot" :class="healthStatus"></span>
        API {{ healthStatus === 'up' ? 'Online' : healthStatus === 'down' ? 'Offline' : '...' }}
      </span>
    </div>

    <div class="auth-wrapper">
      <div class="auth-brand">
        <div class="logo-mark" aria-hidden="true">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="20 6 9 17 4 12" />
          </svg>
        </div>
        <h1>Todo</h1>
        <p>Kelola tugas harian Anda dengan rapi dan fokus.</p>
      </div>

      <div class="card auth-card">
        <div class="tabs">
          <button
            type="button"
            class="tab-btn"
            :class="{ active: authMode === 'login' }"
            @click="authMode = 'login'; authError = ''"
          >
            Masuk
          </button>
          <button
            type="button"
            class="tab-btn"
            :class="{ active: authMode === 'register' }"
            @click="authMode = 'register'; authError = ''"
          >
            Daftar
          </button>
        </div>

        <div v-if="authError" class="alert alert-error">{{ authError }}</div>
        <div v-if="authSuccess" class="alert alert-success">{{ authSuccess }}</div>

        <form @submit.prevent="submitAuth">
          <div class="form-group">
            <label for="email">Email</label>
            <input
              id="email"
              v-model="authForm.email"
              type="email"
              placeholder="nama@email.com"
              required
            />
          </div>
          <div class="form-group">
            <label for="password">Password</label>
            <input
              id="password"
              v-model="authForm.password"
              type="password"
              placeholder="Minimal 6 karakter"
              required
            />
          </div>
          <button class="btn btn-primary full-width" :disabled="authLoading">
            {{ authLoading ? 'Memproses...' : authMode === 'login' ? 'Masuk' : 'Buat Akun' }}
          </button>
        </form>
      </div>
    </div>
  </div>

  <!-- Dashboard Layout: Top Navbar + Clean Workspace -->
  <div v-else class="app-shell">
    <nav class="navbar">
      <div class="navbar-inner">
        <div class="nav-brand">
          <div class="logo-mark sm" aria-hidden="true">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="20 6 9 17 4 12" />
            </svg>
          </div>
          <span>Todo</span>
        </div>

        <div class="nav-right">
          <a
            href="https://github.com/Adibmaros/golang-todo-api"
            target="_blank"
            rel="noopener noreferrer"
            class="github-pill"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
              <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12z" />
            </svg>
            GitHub
          </a>
          <span class="status-pill">
            <span class="dot" :class="healthStatus"></span>
            API {{ healthStatus === 'up' ? 'Online' : healthStatus === 'down' ? 'Offline' : '...' }}
          </span>
          <span class="user-chip">{{ userEmail || 'User' }}</span>
          <button class="btn btn-ghost btn-sm" @click="logout">Keluar</button>
        </div>
      </div>
    </nav>

    <main class="workspace">
      <div class="workspace-head">
        <div>
          <h2>Daftar Tugas</h2>
          <p class="subtitle">Total {{ meta.total_data }} tugas tersimpan</p>
        </div>
      </div>

      <form class="composer" @submit.prevent="createTodo">
        <input
          v-model="newTitle"
          type="text"
          placeholder="Tulis tugas baru (min. 3 karakter)..."
          minlength="3"
          maxlength="100"
          required
        />
        <button class="btn btn-primary" type="submit">+ Tambah</button>
      </form>

      <div class="toolbar">
        <div class="filter-pills" role="tablist">
          <button
            type="button"
            class="pill-btn"
            :class="{ active: filters.done === '' }"
            @click="setStatusFilter('')"
          >
            Semua
          </button>
          <button
            type="button"
            class="pill-btn"
            :class="{ active: filters.done === 'false' }"
            @click="setStatusFilter('false')"
          >
            Aktif
          </button>
          <button
            type="button"
            class="pill-btn"
            :class="{ active: filters.done === 'true' }"
            @click="setStatusFilter('true')"
          >
            Selesai
          </button>
        </div>

        <div class="sort-wrap">
          <label for="sort">Urutkan</label>
          <select id="sort" v-model="filters.sort" @change="changeSort">
            <option value="created_at">Terbaru</option>
            <option value="created_asc">Terlama</option>
            <option value="title">Judul (A-Z)</option>
            <option value="done">Status</option>
          </select>
        </div>
      </div>

      <div v-if="actionError" class="alert alert-error">{{ actionError }}</div>

      <div class="card list-card">
        <div v-if="listLoading" class="empty">Memuat daftar tugas...</div>

        <div v-else-if="todos.length === 0" class="empty">
          <div class="empty-icon" aria-hidden="true">
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3" y="4" width="18" height="18" rx="3" />
              <path d="M9 12l2 2 4-4" />
            </svg>
          </div>
          <p class="empty-title">Belum ada tugas</p>
          <p class="empty-desc">Tambahkan tugas pertama Anda di kolom atas.</p>
        </div>

        <ul v-else class="todo-list">
          <li v-for="todo in todos" :key="todo.id" class="todo-row">
            <div class="todo-left">
              <button
                type="button"
                class="check-btn"
                :class="{ checked: todo.done }"
                :aria-label="todo.done ? 'Tandai belum selesai' : 'Tandai selesai'"
                @click="toggleDone(todo)"
              >
                <svg v-if="todo.done" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="20 6 9 17 4 12" />
                </svg>
              </button>

              <template v-if="editingId === todo.id">
                <input
                  v-model="editingTitle"
                  type="text"
                  class="inline-edit"
                  @keyup.enter="saveEdit(todo)"
                  @keyup.esc="cancelEdit"
                />
              </template>
              <template v-else>
                <span class="todo-title" :class="{ done: todo.done }">{{ todo.title }}</span>
              </template>
            </div>

            <div class="todo-right">
              <span v-if="editingId !== todo.id && todo.created_at" class="todo-date">
                {{ formatDate(todo.created_at) }}
              </span>
              <div class="todo-actions">
                <template v-if="editingId === todo.id">
                  <button class="btn btn-primary btn-sm" @click="saveEdit(todo)">Simpan</button>
                  <button class="btn btn-ghost btn-sm" @click="cancelEdit">Batal</button>
                </template>
                <template v-else>
                  <button class="btn btn-ghost btn-sm" @click="startEdit(todo)">Edit</button>
                  <button class="btn btn-danger-ghost btn-sm" @click="removeTodo(todo.id)">Hapus</button>
                </template>
              </div>
            </div>
          </li>
        </ul>

        <div v-if="meta.total_data > 0" class="pagination">
          <span>
            Halaman {{ meta.page }} dari {{ meta.total_pages || 1 }}
          </span>
          <div class="pagination-controls">
            <button
              class="btn btn-secondary btn-sm"
              :disabled="meta.page <= 1"
              @click="goToPage(meta.page - 1)"
            >
              Sebelumnya
            </button>
            <button
              class="btn btn-secondary btn-sm"
              :disabled="meta.page >= meta.total_pages"
              @click="goToPage(meta.page + 1)"
            >
              Berikutnya
            </button>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
