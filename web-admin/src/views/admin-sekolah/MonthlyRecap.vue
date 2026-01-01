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
  Select,
  SelectOption,
  Progress,
  Statistic,
  Empty,
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
import { schoolService, exportService } from '@/services'
import type { Class } from '@/types/school'
import type { MonthlyRecapResponse, StudentRecapSummary } from '@/types/export'
import { MONTH_OPTIONS, getYearOptions, getMonthName } from '@/types/export'

const { Title, Text } = Typography

// State
const loading = ref(false)
const exporting = ref(false)
const recapData = ref<MonthlyRecapResponse | null>(null)

// Filter state
const selectedMonth = ref<number>(new Date().getMonth() + 1)
const selectedYear = ref<number>(new Date().getFullYear())
const filterClassId = ref<number | undefined>(undefined)

// Classes for filter
const classes = ref<Class[]>([])
const loadingClasses = ref(false)

// Year options
const yearOptions = getYearOptions()

// Table columns
// Requirements: 2.1 - Display summary per student including total days present, late, very late, and absent
const columns: TableProps['columns'] = [
  {
    title: 'No',
    key: 'index',
    width: 60,
    align: 'center',
  },
  {
    title: 'NIS',
    dataIndex: 'studentNis',
    key: 'studentNis',
    width: 120,
  },
  {
    title: 'Nama Siswa',
    dataIndex: 'studentName',
    key: 'studentName',
    ellipsis: true,
  },
  {
    title: 'Kelas',
    dataIndex: 'className',
    key: 'className',
    width: 100,
  },
  {
    title: 'Hadir',
    dataIndex: 'totalPresent',
    key: 'totalPresent',
    width: 80,
    align: 'center',
  },
  {
    title: 'Terlambat',
    dataIndex: 'totalLate',
    key: 'totalLate',
    width: 100,
    align: 'center',
  },
  {
    title: 'Sangat Terlambat',
    dataIndex: 'totalVeryLate',
    key: 'totalVeryLate',
    width: 130,
    align: 'center',
  },
  {
    title: 'Tidak Hadir',
    dataIndex: 'totalAbsent',
    key: 'totalAbsent',
    width: 100,
    align: 'center',
  },
  {
    title: 'Persentase',
    key: 'percentage',
    width: 150,
    align: 'center',
  },
]

// Computed summary statistics
const summaryStats = computed(() => {
  if (!recapData.value || !recapData.value.studentRecaps.length) {
    return {
      totalStudents: 0,
      avgPresent: 0,
      avgLate: 0,
      avgAbsent: 0,
      avgPercentage: 0,
    }
  }
  
  const students = recapData.value.studentRecaps
  const totalStudents = students.length
  const totalPresent = students.reduce((sum, s) => sum + s.totalPresent, 0)
  const totalLate = students.reduce((sum, s) => sum + s.totalLate, 0)
  const totalAbsent = students.reduce((sum, s) => sum + s.totalAbsent, 0)
  const avgPercentage = students.reduce((sum, s) => sum + s.attendancePercent, 0) / totalStudents
  
  return {
    totalStudents,
    avgPresent: Math.round(totalPresent / totalStudents),
    avgLate: Math.round(totalLate / totalStudents),
    avgAbsent: Math.round(totalAbsent / totalStudents),
    avgPercentage: Math.round(avgPercentage * 100) / 100,
  }
})

// Get percentage color
const getPercentageColor = (percentage: number): string => {
  if (percentage >= 95) return '#22c55e'
  if (percentage >= 85) return '#f97316'
  return '#ef4444'
}

// Load monthly recap data
// Requirements: 2.1 - Display summary per student
const loadRecap = async () => {
  loading.value = true
  try {
    const response = await exportService.getMonthlyRecap({
      month: selectedMonth.value,
      year: selectedYear.value,
      classId: filterClassId.value,
    })
    recapData.value = response
  } catch (err) {
    console.error('Failed to load monthly recap:', err)
    message.error('Gagal memuat data rekap bulanan')
    recapData.value = null
  } finally {
    loading.value = false
  }
}

// Load classes for filter
const loadClasses = async () => {
  loadingClasses.value = true
  try {
    const response = await schoolService.getClasses({ page_size: 100 })
    classes.value = response.classes
  } catch (err) {
    console.error('Failed to load classes:', err)
    classes.value = []
  } finally {
    loadingClasses.value = false
  }
}

// Handle filter change
const handleFilterChange = () => {
  loadRecap()
}

// Handle export
// Requirements: 2.5 - Export monthly recap to Excel
const handleExport = async () => {
  exporting.value = true
  try {
    await exportService.exportMonthlyRecap({
      month: selectedMonth.value,
      year: selectedYear.value,
      classId: filterClassId.value,
    })
    message.success('Rekap bulanan berhasil diunduh')
  } catch (err) {
    console.error('Failed to export monthly recap:', err)
    message.error('Gagal mengunduh rekap bulanan')
  } finally {
    exporting.value = false
  }
}

// Format period for display
const formattedPeriod = computed(() => {
  return `${getMonthName(selectedMonth.value)} ${selectedYear.value}`
})

onMounted(() => {
  loadRecap()
  loadClasses()
})
</script>

<template>
  <div class="monthly-recap">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Rekap Bulanan Absensi</Title>
      <Text type="secondary">
        <CalendarOutlined /> {{ formattedPeriod }}
        <span v-if="recapData"> - {{ recapData.totalDays }} hari sekolah</span>
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
            title="Rata-rata Hadir"
            :value="summaryStats.avgPresent"
            suffix="hari"
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
            title="Rata-rata Terlambat"
            :value="summaryStats.avgLate"
            suffix="hari"
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
            title="Rata-rata Kehadiran"
            :value="summaryStats.avgPercentage"
            suffix="%"
            :value-style="{ color: getPercentageColor(summaryStats.avgPercentage) }"
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
            <Select
              v-model:value="selectedMonth"
              style="width: 150px"
              @change="handleFilterChange"
            >
              <SelectOption v-for="month in MONTH_OPTIONS" :key="month.value" :value="month.value">
                {{ month.label }}
              </SelectOption>
            </Select>
            <Select
              v-model:value="selectedYear"
              style="width: 120px"
              @change="handleFilterChange"
            >
              <SelectOption v-for="year in yearOptions" :key="year.value" :value="year.value">
                {{ year.label }}
              </SelectOption>
            </Select>
            <Select
              v-model:value="filterClassId"
              placeholder="Semua Kelas"
              allow-clear
              style="width: 150px"
              :loading="loadingClasses"
              @change="handleFilterChange"
            >
              <SelectOption v-for="cls in classes" :key="cls.id" :value="cls.id">
                {{ cls.name }}
              </SelectOption>
            </Select>
          </Space>
        </Col>
        <Col :xs="24" :sm="24" :md="8" class="toolbar-right">
          <Space>
            <Button @click="loadRecap">
              <template #icon><ReloadOutlined /></template>
              Refresh
            </Button>
            <Button type="primary" :loading="exporting" @click="handleExport">
              <template #icon><DownloadOutlined /></template>
              Export Excel
            </Button>
          </Space>
        </Col>
      </Row>

      <!-- Table -->
      <Table
        :columns="columns"
        :data-source="recapData?.studentRecaps || []"
        :loading="loading"
        :pagination="{ pageSize: 20, showSizeChanger: true, showTotal: (total: number) => `Total ${total} siswa` }"
        row-key="studentId"
        :scroll="{ x: 1000 }"
      >
        <template #bodyCell="{ column, record, index }">
          <template v-if="column.key === 'index'">
            {{ index + 1 }}
          </template>
          <template v-else-if="column.key === 'totalPresent'">
            <Tag color="success">{{ (record as StudentRecapSummary).totalPresent }}</Tag>
          </template>
          <template v-else-if="column.key === 'totalLate'">
            <Tag v-if="(record as StudentRecapSummary).totalLate > 0" color="warning">
              {{ (record as StudentRecapSummary).totalLate }}
            </Tag>
            <span v-else>0</span>
          </template>
          <template v-else-if="column.key === 'totalVeryLate'">
            <Tag v-if="(record as StudentRecapSummary).totalVeryLate > 0" color="orange">
              {{ (record as StudentRecapSummary).totalVeryLate }}
            </Tag>
            <span v-else>0</span>
          </template>
          <template v-else-if="column.key === 'totalAbsent'">
            <Tag v-if="(record as StudentRecapSummary).totalAbsent > 0" color="error">
              {{ (record as StudentRecapSummary).totalAbsent }}
            </Tag>
            <span v-else>0</span>
          </template>
          <template v-else-if="column.key === 'percentage'">
            <div class="percentage-cell">
              <Progress
                :percent="Math.round((record as StudentRecapSummary).attendancePercent)"
                :stroke-color="getPercentageColor((record as StudentRecapSummary).attendancePercent)"
                :show-info="true"
                size="small"
              />
            </div>
          </template>
        </template>

        <!-- Empty state -->
        <template #emptyText>
          <Empty description="Tidak ada data rekap untuk periode ini" />
        </template>
      </Table>
    </Card>
  </div>
</template>

<style scoped>
.monthly-recap {
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
  min-width: 100px;
}

@media (max-width: 768px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
}
</style>
