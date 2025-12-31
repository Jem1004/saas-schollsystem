<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import {
  Table,
  Button,
  Space,
  Tag,
  Card,
  Row,
  Col,
  Typography,
  DatePicker,
  Select,
  SelectOption,
  Progress,
  Statistic,
  message,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  ReloadOutlined,
  DownloadOutlined,
  CalendarOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  ClockCircleOutlined,
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import type { Dayjs } from 'dayjs'
import { schoolService } from '@/services'
import type { AttendanceSummary, Class } from '@/types/school'

const { Title, Text } = Typography

// State
const loading = ref(false)
const attendanceData = ref<AttendanceSummary[]>([])
const selectedDate = ref<Dayjs>(dayjs())
const filterClassId = ref<number | undefined>(undefined)

// Classes for filter
const classes = ref<Class[]>([])
const loadingClasses = ref(false)

// Export state
const exporting = ref(false)

// Mock data for development
const mockAttendance: AttendanceSummary[] = [
  { date: dayjs().format('YYYY-MM-DD'), classId: 1, className: 'VII-A', totalStudents: 30, present: 28, absent: 1, late: 1, excused: 0 },
  { date: dayjs().format('YYYY-MM-DD'), classId: 2, className: 'VII-B', totalStudents: 30, present: 29, absent: 0, late: 1, excused: 0 },
  { date: dayjs().format('YYYY-MM-DD'), classId: 3, className: 'VIII-A', totalStudents: 32, present: 30, absent: 1, late: 1, excused: 0 },
  { date: dayjs().format('YYYY-MM-DD'), classId: 4, className: 'VIII-B', totalStudents: 28, present: 27, absent: 0, late: 1, excused: 0 },
  { date: dayjs().format('YYYY-MM-DD'), classId: 5, className: 'IX-A', totalStudents: 30, present: 28, absent: 2, late: 0, excused: 0 },
  { date: dayjs().format('YYYY-MM-DD'), classId: 6, className: 'IX-B', totalStudents: 29, present: 27, absent: 1, late: 1, excused: 0 },
]

const mockClasses: Class[] = [
  { id: 1, schoolId: 1, name: 'VII-A', grade: 7, year: '2024/2025', createdAt: '2024-01-15', updatedAt: '2024-01-15' },
  { id: 2, schoolId: 1, name: 'VII-B', grade: 7, year: '2024/2025', createdAt: '2024-01-15', updatedAt: '2024-01-15' },
  { id: 3, schoolId: 1, name: 'VIII-A', grade: 8, year: '2024/2025', createdAt: '2024-01-15', updatedAt: '2024-01-15' },
  { id: 4, schoolId: 1, name: 'VIII-B', grade: 8, year: '2024/2025', createdAt: '2024-01-15', updatedAt: '2024-01-15' },
  { id: 5, schoolId: 1, name: 'IX-A', grade: 9, year: '2024/2025', createdAt: '2024-01-15', updatedAt: '2024-01-15' },
  { id: 6, schoolId: 1, name: 'IX-B', grade: 9, year: '2024/2025', createdAt: '2024-01-15', updatedAt: '2024-01-15' },
]

// Table columns
const columns: TableProps['columns'] = [
  {
    title: 'Kelas',
    dataIndex: 'className',
    key: 'className',
    width: 100,
  },
  {
    title: 'Total Siswa',
    dataIndex: 'totalStudents',
    key: 'totalStudents',
    width: 100,
    align: 'center',
  },
  {
    title: 'Hadir',
    dataIndex: 'present',
    key: 'present',
    width: 100,
    align: 'center',
  },
  {
    title: 'Terlambat',
    dataIndex: 'late',
    key: 'late',
    width: 100,
    align: 'center',
  },
  {
    title: 'Tidak Hadir',
    dataIndex: 'absent',
    key: 'absent',
    width: 100,
    align: 'center',
  },
  {
    title: 'Persentase Kehadiran',
    key: 'percentage',
    width: 200,
  },
]

// Computed summary statistics
const summaryStats = computed(() => {
  const data = filteredAttendance.value
  const totalStudents = data.reduce((sum, item) => sum + item.totalStudents, 0)
  const totalPresent = data.reduce((sum, item) => sum + item.present, 0)
  const totalLate = data.reduce((sum, item) => sum + item.late, 0)
  const totalAbsent = data.reduce((sum, item) => sum + item.absent, 0)
  
  return {
    totalStudents,
    totalPresent,
    totalLate,
    totalAbsent,
    percentage: totalStudents > 0 ? Math.round((totalPresent / totalStudents) * 100) : 0,
  }
})

// Computed filtered data
const filteredAttendance = computed(() => {
  if (!filterClassId.value) return attendanceData.value
  return attendanceData.value.filter(item => item.classId === filterClassId.value)
})

// Get attendance percentage color
const getPercentageColor = (present: number, total: number): string => {
  const percentage = (present / total) * 100
  if (percentage >= 95) return '#22c55e'
  if (percentage >= 85) return '#f97316'
  return '#ef4444'
}

// Load attendance data
const loadAttendance = async () => {
  loading.value = true
  try {
    const response = await schoolService.getAttendanceSummary({
      date: selectedDate.value.format('YYYY-MM-DD'),
      classId: filterClassId.value,
    })
    attendanceData.value = response.data
  } catch {
    // Use mock data on error
    attendanceData.value = mockAttendance
  } finally {
    loading.value = false
  }
}

// Load classes for filter
const loadClasses = async () => {
  loadingClasses.value = true
  try {
    const response = await schoolService.getClasses({ pageSize: 100 })
    classes.value = response.data
  } catch {
    classes.value = mockClasses
  } finally {
    loadingClasses.value = false
  }
}

// Handle date change
const handleDateChange = (date: string | Dayjs) => {
  if (date && typeof date !== 'string') {
    selectedDate.value = date
    loadAttendance()
  }
}

// Handle class filter change
const handleClassFilterChange = () => {
  loadAttendance()
}

// Handle export
const handleExport = async () => {
  exporting.value = true
  try {
    const startDate = selectedDate.value.startOf('month').format('YYYY-MM-DD')
    const endDate = selectedDate.value.endOf('month').format('YYYY-MM-DD')
    
    const blob = await schoolService.exportAttendance({
      startDate,
      endDate,
      classId: filterClassId.value,
    })
    
    // Create download link
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `attendance_${selectedDate.value.format('YYYY-MM')}.xlsx`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    message.success('Laporan berhasil diunduh')
  } catch {
    message.error('Gagal mengunduh laporan')
  } finally {
    exporting.value = false
  }
}

// Format date for display
const formattedDate = computed(() => {
  return selectedDate.value.format('dddd, D MMMM YYYY')
})

onMounted(() => {
  loadAttendance()
  loadClasses()
})
</script>

<template>
  <div class="attendance-monitoring">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Monitoring Absensi</Title>
      <Text type="secondary">
        <CalendarOutlined /> {{ formattedDate }}
      </Text>
    </div>

    <!-- Summary Cards -->
    <Row :gutter="[24, 24]" class="summary-row">
      <Col :xs="24" :sm="12" :lg="6">
        <Card class="stat-card">
          <Statistic
            title="Total Siswa"
            :value="summaryStats.totalStudents"
            :value-style="{ color: '#3b82f6' }"
          />
        </Card>
      </Col>
      <Col :xs="24" :sm="12" :lg="6">
        <Card class="stat-card">
          <Statistic
            title="Hadir"
            :value="summaryStats.totalPresent"
            :value-style="{ color: '#22c55e' }"
          >
            <template #prefix>
              <CheckCircleOutlined />
            </template>
          </Statistic>
        </Card>
      </Col>
      <Col :xs="24" :sm="12" :lg="6">
        <Card class="stat-card">
          <Statistic
            title="Terlambat"
            :value="summaryStats.totalLate"
            :value-style="{ color: '#f97316' }"
          >
            <template #prefix>
              <ClockCircleOutlined />
            </template>
          </Statistic>
        </Card>
      </Col>
      <Col :xs="24" :sm="12" :lg="6">
        <Card class="stat-card">
          <Statistic
            title="Tidak Hadir"
            :value="summaryStats.totalAbsent"
            :value-style="{ color: '#ef4444' }"
          >
            <template #prefix>
              <CloseCircleOutlined />
            </template>
          </Statistic>
        </Card>
      </Col>
    </Row>

    <Card style="margin-top: 24px">
      <!-- Toolbar -->
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col :xs="24" :sm="24" :md="16">
          <Space wrap>
            <DatePicker
              :value="selectedDate"
              format="DD MMMM YYYY"
              :allow-clear="false"
              style="width: 200px"
              @change="handleDateChange"
            />
            <Select
              v-model:value="filterClassId"
              placeholder="Semua Kelas"
              allow-clear
              style="width: 150px"
              :loading="loadingClasses"
              @change="handleClassFilterChange"
            >
              <SelectOption v-for="cls in classes" :key="cls.id" :value="cls.id">
                {{ cls.name }}
              </SelectOption>
            </Select>
          </Space>
        </Col>
        <Col :xs="24" :sm="24" :md="8" class="toolbar-right">
          <Space>
            <Button @click="loadAttendance">
              <template #icon><ReloadOutlined /></template>
              Refresh
            </Button>
            <Button type="primary" :loading="exporting" @click="handleExport">
              <template #icon><DownloadOutlined /></template>
              Export
            </Button>
          </Space>
        </Col>
      </Row>

      <!-- Table -->
      <Table
        :columns="columns"
        :data-source="filteredAttendance"
        :loading="loading"
        :pagination="false"
        row-key="classId"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'className'">
            <Tag color="blue">{{ (record as AttendanceSummary).className }}</Tag>
          </template>
          <template v-else-if="column.key === 'present'">
            <Tag color="success">{{ (record as AttendanceSummary).present }}</Tag>
          </template>
          <template v-else-if="column.key === 'late'">
            <Tag v-if="(record as AttendanceSummary).late > 0" color="warning">
              {{ (record as AttendanceSummary).late }}
            </Tag>
            <span v-else>0</span>
          </template>
          <template v-else-if="column.key === 'absent'">
            <Tag v-if="(record as AttendanceSummary).absent > 0" color="error">
              {{ (record as AttendanceSummary).absent }}
            </Tag>
            <span v-else>0</span>
          </template>
          <template v-else-if="column.key === 'percentage'">
            <div class="percentage-cell">
              <Progress
                :percent="Math.round(((record as AttendanceSummary).present / (record as AttendanceSummary).totalStudents) * 100)"
                :stroke-color="getPercentageColor((record as AttendanceSummary).present, (record as AttendanceSummary).totalStudents)"
                :show-info="true"
                size="small"
              />
            </div>
          </template>
        </template>
        
        <!-- Summary Footer -->
        <template #summary>
          <Table.Summary fixed>
            <Table.Summary.Row class="summary-row">
              <Table.Summary.Cell :index="0">
                <Text strong>Total</Text>
              </Table.Summary.Cell>
              <Table.Summary.Cell :index="1" align="center">
                <Text strong>{{ summaryStats.totalStudents }}</Text>
              </Table.Summary.Cell>
              <Table.Summary.Cell :index="2" align="center">
                <Tag color="success">{{ summaryStats.totalPresent }}</Tag>
              </Table.Summary.Cell>
              <Table.Summary.Cell :index="3" align="center">
                <Tag v-if="summaryStats.totalLate > 0" color="warning">{{ summaryStats.totalLate }}</Tag>
                <span v-else>0</span>
              </Table.Summary.Cell>
              <Table.Summary.Cell :index="4" align="center">
                <Tag v-if="summaryStats.totalAbsent > 0" color="error">{{ summaryStats.totalAbsent }}</Tag>
                <span v-else>0</span>
              </Table.Summary.Cell>
              <Table.Summary.Cell :index="5">
                <div class="percentage-cell">
                  <Progress
                    :percent="summaryStats.percentage"
                    :stroke-color="getPercentageColor(summaryStats.totalPresent, summaryStats.totalStudents)"
                    :show-info="true"
                    size="small"
                  />
                </div>
              </Table.Summary.Cell>
            </Table.Summary.Row>
          </Table.Summary>
        </template>
      </Table>
    </Card>
  </div>
</template>

<style scoped>
.attendance-monitoring {
  padding: 0;
}

.page-header {
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.summary-row {
  margin-bottom: 0;
}

.stat-card {
  height: 100%;
}

.stat-card :deep(.ant-statistic-title) {
  font-size: 14px;
  color: #8c8c8c;
}

.stat-card :deep(.ant-statistic-content-prefix) {
  margin-right: 8px;
}

.toolbar {
  margin-bottom: 16px;
}

.toolbar-right {
  display: flex;
  justify-content: flex-end;
}

.percentage-cell {
  min-width: 150px;
}

:deep(.summary-row) {
  background: #fafafa;
}

@media (max-width: 768px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
}
</style>
