<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { http } from '../../services/http'
import { useAuthStore } from '../../store/auth'
import {
  Table,
  Button,
  Modal,
  Input,
  Switch,
  Typography,
  Select,
  Toast
} from '@kousum/semi-ui-vue'
import type { ColumnProps } from '@kousum/semi-ui-vue/dist/table/interface'

const { Title, Text } = Typography
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

// Load all users for admin management table
const fetchUsers = async () => {
  loading.value = true
  try {
    const res = await http.get<UserRow[]>('/users')
    users.value = res.data
  } catch (e) {
    console.error('Failed to fetch users', e)
    Toast.error('Failed to load users')
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
    Toast.success('User updated')
    showEditModal.value = false
  } catch (e) {
    console.error('Failed to update user', e)
    Toast.error('Failed to update user')
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
    Toast.warning('Username and passwords are required')
    return
  }
  if (password !== confirm) {
    Toast.warning('Passwords do not match')
    return
  }
  try {
    await http.post('/users', {
      username,
      password,
      description: newDescription.value,
      type: newType.value
    })
    Toast.success('User created')
    showAddModal.value = false
    await fetchUsers()
  } catch (e) {
    console.error('Failed to create user', e)
    Toast.error('Failed to create user')
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
    Toast.warning('New password and confirmation are required')
    return
  }
  if (passwordNew.value !== passwordConfirm.value) {
    Toast.warning('Passwords do not match')
    return
  }
  try {
    await http.patch(`/users/${passwordUser.value.id}/password`, {
      new_password: passwordNew.value,
      new_password_confirm: passwordConfirm.value,
      old_password: passwordOld.value || undefined
    })
    Toast.success('Password updated')
    showPasswordModal.value = false
  } catch (e) {
    console.error('Failed to change password', e)
    Toast.error('Failed to change password')
  }
}

const columns = computed<ColumnProps<UserRow>[]>(() => {
  const cols: ColumnProps<UserRow>[] = [
    { title: 'Username', dataIndex: 'username', key: 'username' },
    {
      title: 'Role',
      dataIndex: 'is_admin',
      key: 'role',
      render: (_text, record) => (record.is_admin ? 'Admin' : 'User')
    },
    { title: 'Description', dataIndex: 'description', key: 'description' }
  ]

  if (auth.isAdmin) {
    cols.push({
      title: 'Enabled',
      dataIndex: 'enabled',
      key: 'enabled',
      render: (_text, record) => record.enabled ? 'Yes' : 'No'
    })
  }

  cols.push({
    title: 'Actions',
    key: 'actions',
    render: (_text, record) => null
  })

  return cols
})

onMounted(() => {
  fetchUsers()
})
</script>

<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between mb-2">
      <div>
        <Title :heading="5">User Management</Title>
        <Text type="secondary" class="text-sm">
          Manage users and their roles.
        </Text>
      </div>
      <Button v-if="auth.isAdmin" type="primary" theme="solid" @click="openAdd">
        Add User
      </Button>
    </div>

    <Table
      :columns="columns"
      :dataSource="users"
      rowKey="id"
      :loading="loading"
      :pagination="false"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'actions'">
          <div class="flex gap-2 justify-end">
            <Button size="small" theme="borderless" @click="openEdit(record)">
              Edit
            </Button>
            <Button size="small" theme="borderless" @click="openChangePassword(record)">
              Change Password
            </Button>
          </div>
        </template>
        <template v-else-if="column.key === 'enabled'">
          <Switch
            v-if="auth.isAdmin"
            :checked="record.enabled"
            disabled
          />
        </template>
        <template v-else>
          {{ (record as any)[column.dataIndex as string] }}
        </template>
      </template>
    </Table>

    <!-- Edit User Modal -->
    <Modal
      :visible="showEditModal"
      title="Edit User"
      :okText="'Save'"
      @ok="saveEdit"
      @cancel="showEditModal = false"
    >
      <div class="space-y-3">
        <div v-if="auth.isAdmin">
          <label class="block text-sm font-medium mb-1">Username</label>
          <Input :value="editUsername" @change="val => (editUsername = val)" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Description</label>
          <Input :value="editDescription" @change="val => (editDescription = val)" />
        </div>
        <div v-if="auth.isAdmin">
          <label class="block text-sm font-medium mb-1">Enabled</label>
          <Switch v-model:checked="editEnabled" />
        </div>
      </div>
    </Modal>

    <!-- Add User Modal -->
    <Modal
      :visible="showAddModal"
      title="Add User"
      :okText="'Create'"
      @ok="createUser"
      @cancel="showAddModal = false"
    >
      <div class="space-y-3">
        <div>
          <label class="block text-sm font-medium mb-1">Type</label>
          <Select v-model="newType" style="width: 100%">
            <Select.Option value="user">User</Select.Option>
            <Select.Option value="admin">Admin</Select.Option>
          </Select>
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Username</label>
          <Input :value="newUsername" @change="val => (newUsername = val)" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Password</label>
          <Input :value="newPassword" type="password" @change="val => (newPassword = val)" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Confirm Password</label>
          <Input
            :value="newPasswordConfirm"
            type="password"
            @change="val => (newPasswordConfirm = val)"
          />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Description</label>
          <Input :value="newDescription" @change="val => (newDescription = val)" />
        </div>
      </div>
    </Modal>

    <!-- Change Password Modal -->
    <Modal
      :visible="showPasswordModal"
      title="Change Password"
      :okText="'Update'"
      @ok="changePassword"
      @cancel="showPasswordModal = false"
    >
      <div class="space-y-3">
        <div>
          <label class="block text-sm font-medium mb-1">Old Password (for self)</label>
          <Input :value="passwordOld" type="password" @change="val => (passwordOld = val)" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">New Password</label>
          <Input :value="passwordNew" type="password" @change="val => (passwordNew = val)" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Confirm New Password</label>
          <Input
            :value="passwordConfirm"
            type="password"
            @change="val => (passwordConfirm = val)"
          />
        </div>
      </div>
    </Modal>
  </div>
</template>
