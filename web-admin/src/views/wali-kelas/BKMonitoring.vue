<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import {
  Card, Row, Col, Typography, Table, Tag, Select, SelectOption, Spin, Segmented, Alert, Empty,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  WarningOutlined, TrophyOutlined, FileProtectOutlined, EyeOutlined,
} from '@ant-design/icons-vue'
import { bkService } from '@/services'
import { useClassInfo, useClassStudents, useDateFormat, extractArrayFromResponse } from '@/composables/useWaliKelas'
import type { Violation, Achievement, Permit } from '@/types/bk'

const { Title, Text } = Typography

// Composables
const { className, loadClassInfo } = useClassInfo()
const { students, loadStudents } = useClassStudents()
const { formatShortDate } = useDateFormat()

// Mounted state
const isMounted = ref(true)

// State
const loading = ref(false)
const error = ref<string | null>(null)
const activeTab = ref<string>('violations')
const selectedStudentId = ref<number | undefined>(undefined)

// Data
const violations = ref<Violation[]>([])
const achievements = ref<Achievement[]>([])
const permits = ref<Permit[]>([])

// Tab options
const tabOptions = [
  { value: 'violations', label: 'Pelanggaran' },
  { value: 'achievements', label: 'Prestasi' },
  { value: 'permits', label: 'Izin Keluar' },
]

// Table columns
const violationColumns: TableProps['columns'] = [
  { title: 'Siswa', dataIndex: 'studentName', key: 'studentName' },
  { title: 'Kategori', dataIndex: 'category', key: 'category' },
  { title: 'Level', dataIndex: 'level', key: 'level', width: 100 },
  { title: 'Tanggal', dataIndex: 'createdAt', key: 'createdAt', width: 120 },
]

const achievementColumns: TableProps['columns'] = [
  { title: 'Siswa', dataIndex: 'studentName', key: 'studentName' },
  { title: 'Prestasi', dataIndex: 'title', key: 'title' },
  { title: 'Poin', dataIndex: 'point', key: 'point', width: 80 },
  { title: 'Tanggal', dataIndex: 'createdAt', key: 'createdAt', width: 120 },
]

const permitColumns: TableProps['columns'] = [
  { title: 'Siswa', dataIndex: 'studentName', key: 'studentName' },
  { title: 'Alasan', dataIndex: 'reason', key: 'reason', ellipsis: true },
  { title: 'Keluar', dataIndex: 'exitTime', key: 'exitTime', width: 80 },
  { title: 'Kembali', dataIndex: 'returnTime', key: 'returnTime', width: 80 },
]

// Helpers
const getLevelColor = (level: string) => {
  const colors: Record<string, string> = { ringan: 'warning', sedang: 'orange', berat: 'error' }
  return colors[level] || 'default'
}

// Get student IDs from class
const classStudentIds = computed(() => students.value.map(s => s.id))

// Filter data by class students
const filterByClassStudents = <T extends { studentId: number }>(data: T[]): T[] => {
  if (classStudentIds.value.length === 0) return data
  return data.filter(item => classStudentIds.value.includes(item.studentId))
}

// Filtered data
const filteredViolations = computed(() => {
  let data = filterByClassStudents(violations.value || [])
  return selectedStudentId.value ? data.filter(v => v.studentId === selectedStudentId.value) : data
})

const filteredAchievements = computed(() => {
  let data = filterByClassStudents(achievements.value || [])
  return selectedStudentId.value ? data.filter(a => a.studentId === selectedStudentId.value) : data
})

const filteredPermits = computed(() => {
  let data = filterByClassStudents(permits.value || [])
  return selectedStudentId.value ? data.filter(p => p.studentId === selectedStudentId.value) : data
})

// Summary
const summary = computed(() => ({
  violations: filteredViolations.value.length,
  achievements: filteredAchievements.value.length,
  permits: filteredPermits.value.length,
}))

// Load data
const loadData = async () => {
  if (!isMounted.value) return
  
  loading.value = true
  error.value = null
  violations.value = []
  achievements.value = []
  permits.value = []
  
  try {
    // Load students and class info first
    await Promise.all([loadStudents(), loadClassInfo()])
    
    if (!isMounted.value) return

    // Load BK data
    const [violationsRes, achievementsRes, permitsRes] = await Promise.allSettled([
      bkService.getViolations({ pageSize: 100 }),
      bkService.getAchievements({ pageSize: 100 }),
      bkService.getPermits({ pageSize: 100 }),
    ])
    
    if (!isMounted.value) return
    
    violations.value = violationsRes.status === 'fulfilled' ? extractArrayFromResponse<Violation>(violationsRes.value) : []
    achievements.value = achievementsRes.status === 'fulfilled' ? extractArrayFromResponse<Achievement>(achievementsRes.value) : []
    permits.value = permitsRes.status === 'fulfilled' ? extractArrayFromResponse<Permit>(permitsRes.value) : []
  } catch (err) {
    console.error('Failed to load BK monitoring data:', err)
    if (isMounted.value) error.value = 'Gagal memuat data monitoring BK'
  } finally {
    if (isMounted.value) loading.value = false
  }
}

onMounted(loadData)
onUnmounted(() => { isMounted.value = false })
</script>

<template>
  <div class="wali-kelas-page">
    <div class="page-header">
      <div>
        <Title :level="2" style="margin: 0"><EyeOutlined /> Monitoring BK</Title>
        <Text type="secondary">Kelas {{ className }} · Mode Baca Saja</Text>
      </div>
      <Select
        v-model:value="selectedStudentId"
        placeholder="Filter siswa"
        allow-clear
        style="width: 200px"
        :loading="loading"
        :get-popup-container="(triggerNode: HTMLElement) => triggerNode.parentNode as HTMLElement"
      >
        <SelectOption v-for="s in students" :key="s.id" :value="s.id">{{ s.name }}</SelectOption>
      </Select>
    </div>

    <Alert v-if="error" type="error" :message="error" show-icon closable style="margin-bottom: 16px" @close="error = null" />

    <!-- Summary -->
    <Row :gutter="16" class="summary-row">
      <Col :span="8">
        <Card size="small" :class="{ active: activeTab === 'violations' }" @click="activeTab = 'violations'">
          <div class="summary-item">
            <WarningOutlined class="icon-violation" />
            <div>
              <div class="count">{{ summary.violations }}</div>
              <Text type="secondary">Pelanggaran</Text>
            </div>
          </div>
        </Card>
      </Col>
      <Col :span="8">
        <Card size="small" :class="{ active: activeTab === 'achievements' }" @click="activeTab = 'achievements'">
          <div class="summary-item">
            <TrophyOutlined class="icon-achievement" />
            <div>
              <div class="count">{{ summary.achievements }}</div>
              <Text type="secondary">Prestasi</Text>
            </div>
          </div>
        </Card>
      </Col>
      <Col :span="8">
        <Card size="small" :class="{ active: activeTab === 'permits' }" @click="activeTab = 'permits'">
          <div class="summary-item">
            <FileProtectOutlined class="icon-permit" />
            <div>
              <div class="count">{{ summary.permits }}</div>
              <Text type="secondary">Izin Keluar</Text>
            </div>
          </div>
        </Card>
      </Col>
    </Row>

    <!-- Content -->
    <Card size="small">
      <Segmented v-model:value="activeTab" :options="tabOptions" block style="margin-bottom: 16px" />
      
      <Spin :spinning="loading">
        <template v-if="activeTab === 'violations'">
          <Table v-if="filteredViolations.length > 0" :columns="violationColumns" :data-source="filteredViolations" :pagination="{ pageSize: 10, size: 'small' }" row-key="id" size="small">
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'level'">
                <Tag :color="getLevelColor((record as Violation).level)" size="small">{{ (record as Violation).level }}</Tag>
              </template>
              <template v-else-if="column.key === 'createdAt'">{{ formatShortDate((record as Violation).createdAt) }}</template>
            </template>
          </Table>
          <Empty v-else description="Tidak ada data pelanggaran" />
        </template>

        <template v-else-if="activeTab === 'achievements'">
          <Table v-if="filteredAchievements.length > 0" :columns="achievementColumns" :data-source="filteredAchievements" :pagination="{ pageSize: 10, size: 'small' }" row-key="id" size="small">
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'point'">
                <Tag color="success" size="small">+{{ (record as Achievement).point }}</Tag>
              </template>
              <template v-else-if="column.key === 'createdAt'">{{ formatShortDate((record as Achievement).createdAt) }}</template>
            </template>
          </Table>
          <Empty v-else description="Tidak ada data prestasi" />
        </template>

        <template v-else-if="activeTab === 'permits'">
          <Table v-if="filteredPermits.length > 0" :columns="permitColumns" :data-source="filteredPermits" :pagination="{ pageSize: 10, size: 'small' }" row-key="id" size="small">
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'returnTime'">
                <span v-if="(record as Permit).returnTime">{{ (record as Permit).returnTime }}</span>
                <Tag v-else color="warning" size="small">-</Tag>
              </template>
            </template>
          </Table>
          <Empty v-else description="Tidak ada data izin keluar" />
        </template>
      </Spin>
    </Card>
  </div>
</template>

<style scoped>
.wali-kelas-page { padding: 0; }
.page-header { margin-bottom: 24px; display: flex; justify-content: space-between; align-items: flex-start; flex-wrap: wrap; gap: 8px; }
.summary-row { margin-bottom: 16px; }
.summary-row .ant-card { cursor: pointer; transition: all 0.2s; border: 2px solid transparent; }
.summary-row .ant-card:hover { border-color: #d9d9d9; }
.summary-row .ant-card.active { border-color: #1890ff; }
.summary-item { display: flex; align-items: center; gap: 12px; }
.summary-item .count { font-size: 20px; font-weight: 600; line-height: 1; }
.icon-violation { font-size: 24px; color: #ef4444; }
.icon-achievement { font-size: 24px; color: #22c55e; }
.icon-permit { font-size: 24px; color: #3b82f6; }
</style>
