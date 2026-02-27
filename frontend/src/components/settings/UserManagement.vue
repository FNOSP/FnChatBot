<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { http } from '../../services/http'
import { useAuthStore } from '../../store/auth'
import { MessagePlugin } from 'tdesign-vue-next'

const auth = useAuthStore()

interface UserRow {
  id: number
  username: string
  description?: string
  is_admin: boolean
  enabled: boolean
  must_change_password: boolean
}

const users = ref<UserRow[]>([])
const loading = ref(false)

const showEditModal = ref(false)
const editingUser = ref<UserRow | null>(null)
const editUsername = ref('')
const editDescription = ref('')
const editEnabled = ref(true)

const showAddModal = ref(false)
const newUsername = ref('')
const newDescription = ref('')
const newPassword = ref('')
const newPasswordConfirm = ref('')
const newType = ref<'admin' | 'user'>('user')

const showPasswordModal = ref(false)
const passwordUser = ref<UserRow | null>(null)
const passwordNew = ref('')
const passwordConfirm = ref('')
const passwordOld = ref('')

const typeOptions = [
  { label: 'User', value: 'user' },
  { label: 'Admin', value: 'admin' }
]

// Load all users for admin management table
const fetchUsers = async () => {
  loading.value = true
  try {
    const res = await http.get<UserRow[]>('/users')
    users.value = res.data
  } catch (e) {
    console.error('Failed to fetch users', e)
    MessagePlugin.error('Failed to load users')
  } finally {
    loading.value = false
  }
}

// Open edit modal and populate editable fields
const openEdit = (user: UserRow) => {
  editingUser.value = { ...user }
  editUsername.value = user.username
  editDescription.value = user.description || ''
  editEnabled.value = user.enabled
  showEditModal.value = true
}

// Persist edited user fields to backend
const saveEdit = async () => {
  if (!editingUser.value) return
  try {
    const payload: any = {
      description: editDescription.value
    }
    if (auth.isAdmin) {
      payload.username = editUsername.value
      payload.enabled = editEnabled.value
    }
    const res = await http.patch<UserRow>(`/users/${editingUser.value.id}`, payload)
    const idx = users.value.findIndex(u => u.id === res.data.id)
    if (idx >= 0) {
      users.value[idx] = res.data
    }
    MessagePlugin.success('User updated')
    showEditModal.value = false
  } catch (e) {
    console.error('Failed to update user', e)
    MessagePlugin.error('Failed to update user')
  }
}

// Open add-user modal and reset form state
const openAdd = () => {
  newUsername.value = ''
  newDescription.value = ''
  newPassword.value = ''
  newPasswordConfirm.value = ''
  newType.value = 'user'
  showAddModal.value = true
}

// Create new user after validating form input
const createUser = async () => {
  const username = newUsername.value.trim()
  const password = newPassword.value
  const confirm = newPasswordConfirm.value

  if (!username || !password || !confirm) {
    MessagePlugin.warning('Username and passwords are required')
    return
  }
  if (password !== confirm) {
    MessagePlugin.warning('Passwords do not match')
    return
  }
  try {
    await http.post('/users', {
      username,
      password,
      description: newDescription.value,
      type: newType.value
    })
    MessagePlugin.success('User created')
    showAddModal.value = false
    await fetchUsers()
  } catch (e) {
    console.error('Failed to create user', e)
    MessagePlugin.error('Failed to create user')
  }
}

// Open change-password modal for selected user
const openChangePassword = (user: UserRow) => {
  passwordUser.value = user
  passwordNew.value = ''
  passwordConfirm.value = ''
  passwordOld.value = ''
  showPasswordModal.value = true
}

// Change user password with basic client-side validation
const changePassword = async () => {
  if (!passwordUser.value) return
  if (!passwordNew.value || !passwordConfirm.value) {
    MessagePlugin.warning('New password and confirmation are required')
    return
  }
  if (passwordNew.value !== passwordConfirm.value) {
    MessagePlugin.warning('Passwords do not match')
    return
  }
  try {
    await http.patch(`/users/${passwordUser.value.id}/password`, {
      new_password: passwordNew.value,
      new_password_confirm: passwordConfirm.value,
      old_password: passwordOld.value || undefined
    })
    MessagePlugin.success('Password updated')
    showPasswordModal.value = false
  } catch (e) {
    console.error('Failed to change password', e)
    MessagePlugin.error('Failed to change password')
  }
}

const columns = computed(() => {
  const cols = [
    { title: 'Username', colKey: 'username' },
    {
      title: 'Role',
      colKey: 'role',
      cell: (h: any, { row }: any) => row.is_admin ? 'Admin' : 'User'
    },
    { title: 'Description', colKey: 'description' }
  ]

  if (auth.isAdmin) {
    cols.push({
      title: 'Enabled',
      colKey: 'enabled',
      cell: (h: any, { row }: any) => row.enabled ? 'Yes' : 'No'
    })
  }

  cols.push({
    title: 'Actions',
    colKey: 'actions',
    width: 200,
    fixed: 'right'
  })

  return cols
})

onMounted(() => {
  fetchUsers()
})
</script>

<template>
  <div class="space-y-4 h-full p-4">
    <div class="flex items-center justify-between mb-2">
      <div>
        <h3 class="text-lg font-bold">User Management</h3>
        <p class="text-sm text-text-secondary">
          Manage users and their roles.
        </p>
      </div>
      <t-button v-if="auth.isAdmin" theme="primary" @click="openAdd">
        Add User
      </t-button>
    </div>

    <t-table
      :columns="columns"
      :data="users"
      row-key="id"
      :loading="loading"
      :pagination="null"
      class="border border-border rounded-lg"
    >
      <template #enabled="{ row }">
        <t-switch
          v-if="auth.isAdmin"
          :value="row.enabled"
          disabled
        />
      </template>
      <template #actions="{ row }">
        <div class="flex gap-2 justify-end">
          <t-button size="small" variant="text" @click="openEdit(row)">
            Edit
          </t-button>
          <t-button size="small" variant="text" @click="openChangePassword(row)">
            Change Password
          </t-button>
        </div>
      </template>
    </t-table>

    <!-- Edit User Modal -->
    <t-dialog
      v-model:visible="showEditModal"
      header="Edit User"
      confirm-btn="Save"
      @confirm="saveEdit"
      @cancel="showEditModal = false"
    >
      <div class="space-y-3">
        <div v-if="auth.isAdmin">
          <label class="block text-sm font-medium mb-1">Username</label>
          <t-input v-model="editUsername" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Description</label>
          <t-input v-model="editDescription" />
        </div>
        <div v-if="auth.isAdmin">
          <label class="block text-sm font-medium mb-1">Enabled</label>
          <t-switch v-model="editEnabled" />
        </div>
      </div>
    </t-dialog>

    <!-- Add User Modal -->
    <t-dialog
      v-model:visible="showAddModal"
      header="Add User"
      confirm-btn="Create"
      @confirm="createUser"
      @cancel="showAddModal = false"
    >
      <div class="space-y-3">
        <div>
          <label class="block text-sm font-medium mb-1">Type</label>
          <t-select v-model="newType" :options="typeOptions" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Username</label>
          <t-input v-model="newUsername" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Password</label>
          <t-input v-model="newPassword" type="password" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Confirm Password</label>
          <t-input v-model="newPasswordConfirm" type="password" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Description</label>
          <t-input v-model="newDescription" />
        </div>
      </div>
    </t-dialog>

    <!-- Change Password Modal -->
    <t-dialog
      v-model:visible="showPasswordModal"
      header="Change Password"
      confirm-btn="Update"
      @confirm="changePassword"
      @cancel="showPasswordModal = false"
    >
      <div class="space-y-3">
        <div>
          <label class="block text-sm font-medium mb-1">Old Password (for self)</label>
          <t-input v-model="passwordOld" type="password" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">New Password</label>
          <t-input v-model="passwordNew" type="password" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Confirm New Password</label>
          <t-input v-model="passwordConfirm" type="password" />
        </div>
      </div>
    </t-dialog>
  </div>
</template>
