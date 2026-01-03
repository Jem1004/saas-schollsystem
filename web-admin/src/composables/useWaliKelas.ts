/**
 * Composables untuk halaman Wali Kelas
 * Menghilangkan duplikasi kode antar halaman
 */

import { ref, onUnmounted } from 'vue'
import { homeroomService } from '@/services'
import type { ClassStudent } from '@/types/homeroom'

/**
 * Format tanggal ke format Indonesia
 */
export const useDateFormat = () => {
  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
    })
  }

  const formatShortDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'short',
    })
  }

  const formatFullDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('id-ID', {
      weekday: 'long',
      day: 'numeric',
      month: 'long',
      year: 'numeric',
    })
  }

  const formatTodayDate = () => {
    return new Date().toLocaleDateString('id-ID', {
      weekday: 'long',
      day: 'numeric',
      month: 'long',
      year: 'numeric',
    })
  }

  return {
    formatDate,
    formatShortDate,
    formatFullDate,
    formatTodayDate,
  }
}

/**
 * Load dan manage data siswa kelas
 */
export const useClassStudents = () => {
  const students = ref<ClassStudent[]>([])
  const loadingStudents = ref(false)
  const isMounted = ref(true)

  const loadStudents = async () => {
    loadingStudents.value = true
    try {
      const response = await homeroomService.getClassStudents({ pageSize: 100 })
      if (isMounted.value) {
        students.value = response.data || []
      }
    } catch (err) {
      console.error('Failed to load students:', err)
      if (isMounted.value) {
        students.value = []
      }
    } finally {
      if (isMounted.value) {
        loadingStudents.value = false
      }
    }
  }

  // Filter option untuk Select component
  const filterStudentOption = (input: string, option: unknown) => {
    const opt = option as { label?: string } | undefined
    return opt?.label?.toLowerCase().includes(input.toLowerCase()) ?? false
  }

  onUnmounted(() => {
    isMounted.value = false
  })

  return {
    students,
    loadingStudents,
    loadStudents,
    filterStudentOption,
  }
}

/**
 * Load dan manage info kelas
 */
export const useClassInfo = () => {
  const className = ref('')
  const classId = ref<number | null>(null)
  const loadingClass = ref(false)
  const isMounted = ref(true)

  const loadClassInfo = async () => {
    loadingClass.value = true
    try {
      const classInfo = await homeroomService.getMyClass()
      if (isMounted.value) {
        className.value = classInfo?.name || ''
        classId.value = classInfo?.id || null
      }
    } catch (err) {
      console.error('Failed to load class info:', err)
      if (isMounted.value) {
        className.value = ''
        classId.value = null
      }
    } finally {
      if (isMounted.value) {
        loadingClass.value = false
      }
    }
  }

  onUnmounted(() => {
    isMounted.value = false
  })

  return {
    className,
    classId,
    loadingClass,
    loadClassInfo,
  }
}

/**
 * Utility untuk warna nilai
 */
export const useScoreColor = () => {
  const getScoreColor = (score: number): string => {
    if (score >= 85) return '#22c55e'
    if (score >= 70) return '#f97316'
    return '#ef4444'
  }

  const getScoreTagColor = (score: number): string => {
    if (score >= 85) return 'success'
    if (score >= 70) return 'warning'
    return 'error'
  }

  return {
    getScoreColor,
    getScoreTagColor,
  }
}

/**
 * Helper untuk extract array dari response API
 */
export const extractArrayFromResponse = <T>(response: unknown): T[] => {
  if (!response) return []
  if (Array.isArray(response)) return response
  if (typeof response === 'object' && response !== null) {
    const obj = response as Record<string, unknown>
    // Try common response structures
    if (Array.isArray(obj.data)) return obj.data
    if (Array.isArray(obj.notes)) return obj.notes
    if (Array.isArray(obj.violations)) return obj.violations
    if (Array.isArray(obj.achievements)) return obj.achievements
    if (Array.isArray(obj.permits)) return obj.permits
    if (Array.isArray(obj.items)) return obj.items
    if (Array.isArray(obj.records)) return obj.records
  }
  return []
}

/**
 * Composable untuk component lifecycle dengan cleanup
 */
export const useMountedState = () => {
  const isMounted = ref(true)

  onUnmounted(() => {
    isMounted.value = false
  })

  return { isMounted }
}
