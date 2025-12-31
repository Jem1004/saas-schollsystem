<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Select, SelectOption, Spin } from 'ant-design-vue'
import { schoolService } from '@/services'
import type { Class } from '@/types/school'

interface Props {
  modelValue?: number | number[]
  placeholder?: string
  mode?: 'default' | 'multiple' | 'tags'
  disabled?: boolean
  allowClear?: boolean
  showGrade?: boolean
  showYear?: boolean
  gradeFilter?: number
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: 'Pilih kelas',
  mode: 'default',
  disabled: false,
  allowClear: true,
  showGrade: false,
  showYear: false,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: number | number[] | undefined): void
  (e: 'change', value: number | number[] | undefined, option: Class | Class[]): void
}>()

const loading = ref(false)
const classes = ref<Class[]>([])

// Mock data for development
const mockClasses: Class[] = [
  { id: 1, schoolId: 1, name: 'VII-A', grade: 7, year: '2024/2025', studentCount: 30, createdAt: '', updatedAt: '' },
  { id: 2, schoolId: 1, name: 'VII-B', grade: 7, year: '2024/2025', studentCount: 30, createdAt: '', updatedAt: '' },
  { id: 3, schoolId: 1, name: 'VIII-A', grade: 8, year: '2024/2025', studentCount: 32, createdAt: '', updatedAt: '' },
  { id: 4, schoolId: 1, name: 'VIII-B', grade: 8, year: '2024/2025', studentCount: 28, createdAt: '', updatedAt: '' },
  { id: 5, schoolId: 1, name: 'IX-A', grade: 9, year: '2024/2025', studentCount: 30, createdAt: '', updatedAt: '' },
  { id: 6, schoolId: 1, name: 'IX-B', grade: 9, year: '2024/2025', studentCount: 29, createdAt: '', updatedAt: '' },
]

const loadClasses = async () => {
  loading.value = true
  try {
    const response = await schoolService.getClasses({ pageSize: 100 })
    let data = response.data
    if (props.gradeFilter) {
      data = data.filter(c => c.grade === props.gradeFilter)
    }
    classes.value = data
  } catch {
    // Use mock data on error
    let data = mockClasses
    if (props.gradeFilter) {
      data = data.filter(c => c.grade === props.gradeFilter)
    }
    classes.value = data
  } finally {
    loading.value = false
  }
}

const handleChange = (value: unknown) => {
  const typedValue = value as number | number[] | undefined
  emit('update:modelValue', typedValue)
  
  if (typedValue === undefined) {
    emit('change', typedValue, [] as Class[])
    return
  }
  
  if (Array.isArray(typedValue)) {
    const selectedClasses = classes.value.filter(c => typedValue.includes(c.id))
    emit('change', typedValue, selectedClasses)
  } else {
    const selectedClass = classes.value.find(c => c.id === typedValue)
    emit('change', typedValue, selectedClass as Class)
  }
}

const formatLabel = (cls: Class): string => {
  let label = cls.name
  if (props.showGrade) {
    label = `Kelas ${cls.grade} - ${label}`
  }
  if (props.showYear) {
    label = `${label} (${cls.year})`
  }
  return label
}

const filterOption = (input: string, option: { label?: string }) => {
  if (!option.label) return false
  return option.label.toLowerCase().includes(input.toLowerCase())
}

onMounted(() => {
  loadClasses()
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
      v-for="cls in classes"
      :key="cls.id"
      :value="cls.id"
      :label="formatLabel(cls)"
    >
      <div class="class-option">
        <span class="class-name">{{ cls.name }}</span>
        <span v-if="showGrade" class="class-grade">Kelas {{ cls.grade }}</span>
        <span v-if="showYear" class="class-year">{{ cls.year }}</span>
        <span v-if="cls.studentCount !== undefined" class="class-count">
          {{ cls.studentCount }} siswa
        </span>
      </div>
    </SelectOption>
  </Select>
</template>

<style scoped>
.class-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.class-name {
  font-weight: 500;
}

.class-grade {
  color: #8c8c8c;
  font-size: 12px;
}

.class-year {
  color: #52c41a;
  font-size: 12px;
}

.class-count {
  color: #8c8c8c;
  font-size: 12px;
  margin-left: auto;
}
</style>
