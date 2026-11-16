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
      authSuccess.value = 'Registrasi berhasil! Silakan login.'
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

function changeFilter() {
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
  <div class="container">
    <header class="header">
      <div class="brand">
        <h1>Todo App</h1>
        <p>Golang + Gin + MySQL + Vue.js</p>
      </div>
      <div class="header-right">
        <span class="badge">
          <span class="dot" :class="healthStatus"></span>
          API: {{ healthStatus === 'up' ? 'Online' : healthStatus === 'down' ? 'Offline' : '...' }}
        </span>
        <template v-if="token">
          <span class="badge">{{ userEmail || 'Authenticated' }}</span>
          <button class="btn btn-secondary btn-sm" @click="logout">Logout</button>
        </template>
      </div>
    </header>

    <section v-if="!token" class="card auth-card">
      <div class="tabs">
        <button
          class="tab-btn"
          :class="{ active: authMode === 'login' }"
          @click="authMode = 'login'; authError = ''"
        >
          Login
        </button>
        <button
          class="tab-btn"
          :class="{ active: authMode === 'register' }"
          @click="authMode = 'register'; authError = ''"
        >
          Register
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
          {{ authLoading ? 'Memproses...' : authMode === 'login' ? 'Masuk' : 'Daftar' }}
        </button>
      </form>
    </section>

    <section v-else class="card">
      <form class="add-form" @submit.prevent="createTodo">
        <input
          v-model="newTitle"
          type="text"
          placeholder="Tambah todo baru (minimal 3 karakter)..."
          minlength="3"
          maxlength="100"
          required
        />
        <button class="btn btn-primary" type="submit">Tambah</button>
      </form>

      <div class="toolbar">
        <select v-model="filters.done" @change="changeFilter">
          <option value="">Semua Status</option>
          <option value="false">Belum Selesai</option>
          <option value="true">Selesai</option>
        </select>

        <select v-model="filters.sort" @change="changeFilter">
          <option value="created_at">Terbaru</option>
          <option value="created_asc">Terlama</option>
          <option value="title">Judul (A-Z)</option>
          <option value="done">Status</option>
        </select>
      </div>

      <div v-if="actionError" class="alert alert-error">{{ actionError }}</div>

      <div v-if="listLoading" class="empty">Memuat data...</div>
      <div v-else-if="todos.length === 0" class="empty">
        Belum ada todo untuk ditampilkan.
      </div>
      <ul v-else class="todo-list">
        <li v-for="todo in todos" :key="todo.id" class="todo-item">
          <div class="todo-left">
            <input
              type="checkbox"
              :checked="todo.done"
              @change="toggleDone(todo)"
            />
            <template v-if="editingId === todo.id">
              <input v-model="editingTitle" type="text" @keyup.enter="saveEdit(todo)" />
            </template>
            <template v-else>
              <span class="todo-title" :class="{ done: todo.done }">{{ todo.title }}</span>
            </template>
          </div>

          <div class="todo-actions">
            <template v-if="editingId === todo.id">
              <button class="btn btn-primary btn-sm" @click="saveEdit(todo)">Simpan</button>
              <button class="btn btn-secondary btn-sm" @click="cancelEdit">Batal</button>
            </template>
            <template v-else>
              <button class="btn btn-secondary btn-sm" @click="startEdit(todo)">Edit</button>
              <button class="btn btn-danger btn-sm" @click="removeTodo(todo.id)">Hapus</button>
            </template>
          </div>
        </li>
      </ul>

      <div class="pagination">
        <span>
          Halaman {{ meta.page }} dari {{ meta.total_pages || 1 }} ({{ meta.total_data }} total)
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
    </section>
  </div>
</template>
