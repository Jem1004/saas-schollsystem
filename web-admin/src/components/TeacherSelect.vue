<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Select, SelectOption, Spin, Tag } from 'ant-design-vue'
import { schoolService } from '@/services'
import type { SchoolUser } from '@/types/school'

interface Props {
  modelValue?: number | number[]
  placeholder?: string
  mode?: 'default' | 'multiple' | 'tags'
  disabled?: boolean
  allowClear?: boolean
  roleFilter?: ('guru' | 'wali_kelas' | 'guru_bk')[]
  showRole?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: 'Pilih guru',
  mode: 'default',
  disabled: false,
  allowClear: true,
  showRole: true,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: number | number[] | undefined): void
  (e: 'change', value: number | number[] | undefined, option: SchoolUser | SchoolUser[]): void
}>()

const loading = ref(false)
const teachers = ref<SchoolUser[]>([])

// Mock data for development
const mockTeachers: SchoolUser[] = [
  { id: 1, schoolId: 1, role: 'guru', username: 'guru1', name: 'Pak Ahmad', isActive: true, mustResetPwd: false, createdAt: '', updatedAt: '' },
  { id: 2, schoolId: 1, role: 'wali_kelas', username: 'walikelas1', name: 'Bu Siti', assignedClassName: 'VII-A', isActive: true, mustResetPwd: false, createdAt: '', updatedAt: '' },
  { id: 3, schoolId: 1, role: 'wali_kelas', username: 'walikelas2', name: 'Pak Budi', assignedClassName: 'VII-B', isActive: true, mustResetPwd: false, createdAt: '', updatedAt: '' },
  { id: 4, schoolId: 1, role: 'guru_bk', username: 'gurubk1', name: 'Bu Dewi', isActive: true, mustResetPwd: false, createdAt: '', updatedAt: '' },
  { id: 5, schoolId: 1, role: 'guru', username: 'guru2', name: 'Pak Eko', isActive: true, mustResetPwd: false, createdAt: '', updatedAt: '' },
]

const roleLabels: Record<string, string> = {
  guru: 'Guru',
  wali_kelas: 'Wali Kelas',
  guru_bk: 'Guru BK',
}

const roleColors: Record<string, string> = {
  guru: 'blue',
  wali_kelas: 'green',
  guru_bk: 'orange',
}

const loadTeachers = async () => {
  loading.value = true
  try {
    const response = await schoolService.getUsers({ pageSize: 100 })
    let data = response.data
    if (props.roleFilter && props.roleFilter.length > 0) {
      data = data.filter(u => props.roleFilter!.includes(u.role))
    }
    teachers.value = data
  } catch {
    // Use mock data on error
    let data = mockTeachers
    if (props.roleFilter && props.roleFilter.length > 0) {
      data = data.filter(u => props.roleFilter!.includes(u.role))
    }
    teachers.value = data
  } finally {
    loading.value = false
  }
}

const handleChange = (value: unknown) => {
  const typedValue = value as number | number[] | undefined
  emit('update:modelValue', typedValue)
  
  if (typedValue === undefined) {
    emit('change', typedValue, [] as SchoolUser[])
    return
  }
  
  if (Array.isArray(typedValue)) {
    const selectedTeachers = teachers.value.filter(t => typedValue.includes(t.id))
    emit('change', typedValue, selectedTeachers)
  } else {
    const selectedTeacher = teachers.value.find(t => t.id === typedValue)
    emit('change', typedValue, selectedTeacher as SchoolUser)
  }
}

const formatLabel = (teacher: SchoolUser): string => {
  let label = teacher.name || teacher.username
  if (props.showRole) {
    label = `${label} (${roleLabels[teacher.role]})`
  }
  return label
}

const filterOption = (input: string, option: { label?: string }) => {
  if (!option.label) return false
  return option.label.toLowerCase().includes(input.toLowerCase())
}

onMounted(() => {
  loadTeachers()
})
</script>

<template>
  <Select
    :value="modelValue"
    :placeholder="placeholder"
    :mode="mode === 'default' ? undefined : mode"
    :disabled="disabled"
    :allow-clear="allowClear"
    :loading="loading"
    show-search
    :filter-option="filterOption"
    :not-found-content="loading ? undefined : 'Tidak ada data'"
    style="width: 100%"
    @change="handleChange"
  >
    <template #notFoundContent>
      <Spin v-if="loading" size="small" />
      <span v-else>Tidak ada data</span>
    </template>
    <SelectOption
      v-for="teacher in teachers"
      :key="teacher.id"
      :value="teacher.id"
      :label="formatLabel(teacher)"
    >
      <div class="teacher-option">
        <span class="teacher-name">{{ teacher.name || teacher.username }}</span>
        <Tag v-if="showRole" :color="roleColors[teacher.role]" size="small">
          {{ roleLabels[teacher.role] }}
        </Tag>
        <span v-if="teacher.assignedClassName" class="teacher-class">
          {{ teacher.assignedClassName }}
        </span>
      </div>
    </SelectOption>
  </Select>
</template>

<style scoped>
.teacher-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.teacher-name {
  flex: 1;
}

.teacher-class {
  color: #8c8c8c;
  font-size: 12px;
}
</style>
