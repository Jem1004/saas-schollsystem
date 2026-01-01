<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { Select, SelectOption, Spin } from 'ant-design-vue'
import { schoolService } from '@/services'
import type { Student } from '@/types/school'

interface Props {
  modelValue?: number | number[]
  placeholder?: string
  mode?: 'default' | 'multiple' | 'tags'
  disabled?: boolean
  allowClear?: boolean
  classId?: number
  showClass?: boolean
  showNis?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: 'Pilih siswa',
  mode: 'default',
  disabled: false,
  allowClear: true,
  showClass: true,
  showNis: true,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: number | number[] | undefined): void
  (e: 'change', value: number | number[] | undefined, option: Student | Student[]): void
}>()

const loading = ref(false)
const students = ref<Student[]>([])
const searchValue = ref('')

const loadStudents = async (search?: string) => {
  loading.value = true
  try {
    const response = await schoolService.getStudents({
      page_size: 100,
      search,
      class_id: props.classId,
    })
    students.value = response.students
  } catch (err) {
    console.error('Failed to load students:', err)
    students.value = []
  } finally {
    loading.value = false
  }
}

const handleSearch = (value: string) => {
  searchValue.value = value
  loadStudents(value)
}

const handleChange = (value: unknown) => {
  const typedValue = value as number | number[] | undefined
  emit('update:modelValue', typedValue)
  
  if (typedValue === undefined) {
    emit('change', typedValue, [] as Student[])
    return
  }
  
  if (Array.isArray(typedValue)) {
    const selectedStudents = students.value.filter(s => typedValue.includes(s.id))
    emit('change', typedValue, selectedStudents)
  } else {
    const selectedStudent = students.value.find(s => s.id === typedValue)
    emit('change', typedValue, selectedStudent as Student)
  }
}

const formatLabel = (student: Student): string => {
  let label = student.name
  if (props.showNis) {
    label = `${student.nis} - ${label}`
  }
  if (props.showClass && student.className) {
    label = `${label} (${student.className})`
  }
  return label
}

// Watch for classId changes
watch(() => props.classId, () => {
  loadStudents(searchValue.value)
})

onMounted(() => {
  loadStudents()
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
    :filter-option="false"
    :not-found-content="loading ? undefined : 'Tidak ada data'"
    style="width: 100%"
    @search="handleSearch"
    @change="handleChange"
  >
    <template #notFoundContent>
      <Spin v-if="loading" size="small" />
      <span v-else>Tidak ada data</span>
    </template>
    <SelectOption
      v-for="student in students"
      :key="student.id"
      :value="student.id"
      :label="formatLabel(student)"
    >
      <div class="student-option">
        <span class="student-name">{{ student.name }}</span>
        <span v-if="showNis" class="student-nis">{{ student.nis }}</span>
        <span v-if="showClass && student.className" class="student-class">{{ student.className }}</span>
      </div>
    </SelectOption>
  </Select>
</template>

<style scoped>
.student-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.student-name {
  flex: 1;
}

.student-nis {
  color: #8c8c8c;
  font-size: 12px;
}

.student-class {
  color: #1890ff;
  font-size: 12px;
  background: #e6f7ff;
  padding: 0 6px;
  border-radius: 4px;
}
</style>
