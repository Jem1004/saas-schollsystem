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
  Modal,
  Tabs,
  TabPane,
  Empty,
  message,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  ReloadOutlined,
  CalendarOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  ClockCircleOutlined,
  DownloadOutlined,
  BarChartOutlined,
  UnorderedListOutlined,
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import type { Dayjs } from 'dayjs'
import { schoolService, exportService } from '@/services'
import type { AttendanceSummary, Class } from '@/types/school'
import type { MonthlyRecapResponse, StudentRecapSummary } from '@/types/export'
import { MONTH_OPTIONS, getYearOptions, getMonthName } from '@/types/export'

const { Title, Text } = Typography
const { RangePicker } = DatePicker

// ==================== SHARED STATE ====================
const activeTab = ref<'daily' | 'monthly'>('daily')
const classes = ref<Class[]>([])
const loadingClasses = ref(false)
const filterClassId = ref<number | undefined>(undefined)

// Export state
const exporting = ref(false)
const exportModalVisible = ref(false)
const exportDateRange = ref<[Dayjs, Dayjs]>([dayjs().startOf('month'), dayjs()])
const exportClassId = ref<number | undefined>(undefined)

// ==================== DAILY TAB STATE ====================
const dailyLoading = ref(false)
const attendanceData = ref<AttendanceSummary[]>([])
const selectedDate = ref<Dayjs>(dayjs())

// ==================== MONTHLY TAB STATE ====================
const monthlyLoading = ref(false)
const recapData = ref<MonthlyRecapResponse | null>(null)
const selectedMonth = ref<number>(new Date().getMonth() + 1)
const selectedYear = ref<number>(new Date().getFullYear())

// Year options for monthly
const yearOptions = getYearOptions()

// ==================== DAILY TABLE COLUMNS ====================
const dailyColumns: TableProps['columns'] = [
  { title: 'Kelas', dataIndex: 'className', key: 'className', width: 100 },
  { title: 'Total Siswa', dataIndex: 'totalStudents', key: 'totalStudents', width: 100, align: 'center' },
  { title: 'Hadir', dataIndex: 'present', key: 'present', width: 100, align: 'center' },
  { title: 'Terlambat', dataIndex: 'late', key: 'late', width: 100, align: 'center' },
  { title: 'Tidak Hadir', dataIndex: 'absent', key: 'absent', width: 100, align: 'center' },
  { title: 'Persentase Kehadiran', key: 'percentage', width: 200 },
]

// ==================== MONTHLY TABLE COLUMNS ====================
const monthlyColumns: TableProps['columns'] = [
  { title: 'No', key: 'index', width: 60, align: 'center' },
  { title: 'NIS', dataIndex: 'student_nis', key: 'student_nis', width: 120 },
  { title: 'Nama Siswa', dataIndex: 'student_name', key: 'student_name', ellipsis: true },
  { title: 'Kelas', dataIndex: 'class_name', key: 'class_name', width: 100 },
  { title: 'Hadir', dataIndex: 'total_present', key: 'total_present', width: 80, align: 'center' },
  { title: 'Terlambat', dataIndex: 'total_late', key: 'total_late', width: 100, align: 'center' },
  { title: 'Sangat Terlambat', dataIndex: 'total_very_late', key: 'total_very_late', width: 130, align: 'center' },
  { title: 'Tidak Hadir', dataIndex: 'total_absent', key: 'total_absent', width: 100, align: 'center' },
  { title: 'Persentase', key: 'percentage', width: 150, align: 'center' },
]

// ==================== COMPUTED: DAILY STATS ====================
const filteredDailyAttendance = computed(() => {
  if (!filterClassId.value) return attendanceData.value
  return attendanceData.value.filter(item => item.classId === filterClassId.value)
})

const dailyStats = computed(() => {
  const data = filteredDailyAttendance.value
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

// ==================== COMPUTED: MONTHLY STATS ====================
const monthlyStats = computed(() => {
  if (!recapData.value || !recapData.value.student_recaps.length) {
    return { totalStudents: 0, avgPresent: 0, avgLate: 0, avgAbsent: 0, avgPercentage: 0 }
  }
  
  const students = recapData.value.student_recaps
  const totalStudents = students.length
  const totalPresent = students.reduce((sum, s) => sum + s.total_present, 0)
  const totalLate = students.reduce((sum, s) => sum + s.total_late, 0)
  const totalAbsent = students.reduce((sum, s) => sum + s.total_absent, 0)
  const avgPercentage = students.reduce((sum, s) => sum + s.attendance_percent, 0) / totalStudents
  
  return {
    totalStudents,
    avgPresent: Math.round(totalPresent / totalStudents),
    avgLate: Math.round(totalLate / totalStudents),
    avgAbsent: Math.round(totalAbsent / totalStudents),
    avgPercentage: Math.round(avgPercentage * 100) / 100,
  }
})

// ==================== COMPUTED DISPLAY ====================
const formattedDailyDate = computed(() => selectedDate.value.format('dddd, D MMMM YYYY'))
const formattedMonthlyPeriod = computed(() => `${getMonthName(selectedMonth.value)} ${selectedYear.value}`)

// ==================== SHARED FUNCTIONS ====================
const getPercentageColor = (percentage: number): string => {
  if (percentage >= 95) return '#22c55e'
  if (percentage >= 85) return '#f97316'
  return '#ef4444'
}

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

// ==================== DAILY FUNCTIONS ====================
const loadDailyAttendance = async () => {
  dailyLoading.value = true
  try {
    const response = await schoolService.getAttendanceSummary({
      date: selectedDate.value.format('YYYY-MM-DD'),
      class_id: filterClassId.value,
    })
    attendanceData.value = response.data
  } catch (err) {
    console.error('Failed to load attendance:', err)
    message.error('Gagal memuat data absensi')
    attendanceData.value = []
  } finally {
    dailyLoading.value = false
  }
}

const handleDateChange = (date: string | Dayjs) => {
  if (date && typeof date !== 'string') {
    selectedDate.value = date
    loadDailyAttendance()
  }
}

// ==================== MONTHLY FUNCTIONS ====================
const loadMonthlyRecap = async () => {
  monthlyLoading.value = true
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
    monthlyLoading.value = false
  }
}

const handleMonthYearChange = () => {
  loadMonthlyRecap()
}

// ==================== EXPORT FUNCTIONS ====================
const openExportModal = () => {
  if (activeTab.value === 'daily') {
    exportDateRange.value = [dayjs().startOf('month'), dayjs()]
  } else {
    const startOfMonth = dayjs().year(selectedYear.value).month(selectedMonth.value - 1).startOf('month')
    const endOfMonth = startOfMonth.endOf('month')
    exportDateRange.value = [startOfMonth, endOfMonth]
  }
  exportClassId.value = filterClassId.value
  exportModalVisible.value = true
}

const handleExport = async () => {
  if (!exportDateRange.value || exportDateRange.value.length !== 2) {
    message.error('Pilih rentang tanggal terlebih dahulu')
    return
  }

  exporting.value = true
  try {
    if (activeTab.value === 'daily') {
      await exportService.exportAttendance({
        startDate: exportDateRange.value[0].format('YYYY-MM-DD'),
        endDate: exportDateRange.value[1].format('YYYY-MM-DD'),
        classId: exportClassId.value,
      })
    } else {
      await exportService.exportMonthlyRecap({
        month: selectedMonth.value,
        year: selectedYear.value,
        classId: exportClassId.value,
      })
    }
    
    message.success('Laporan berhasil diunduh')
    exportModalVisible.value = false
  } catch (err) {
    console.error('Failed to export:', err)
    message.error('Gagal mengunduh laporan')
  } finally {
    exporting.value = false
  }
}

// ==================== FILTER CHANGE ====================
const handleClassFilterChange = () => {
  if (activeTab.value === 'daily') {
    loadDailyAttendance()
  } else {
    loadMonthlyRecap()
  }
}

const handleRefresh = () => {
  if (activeTab.value === 'daily') {
    loadDailyAttendance()
  } else {
    loadMonthlyRecap()
  }
}

// ==================== TAB CHANGE ====================
const handleTabChange = (key: string | number) => {
  activeTab.value = key as 'daily' | 'monthly'
  if (key === 'daily' && attendanceData.value.length === 0) {
    loadDailyAttendance()
  } else if (key === 'monthly' && !recapData.value) {
    loadMonthlyRecap()
  }
}

// ==================== LIFECYCLE ====================
onMounted(() => {
  loadClasses()
  loadDailyAttendance()
})
</script>

<template>
  <div class="attendance-report">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Laporan Absensi</Title>
      <Text type="secondary">
        <CalendarOutlined />
        {{ activeTab === 'daily' ? formattedDailyDate : formattedMonthlyPeriod }}
        <span v-if="activeTab === 'monthly' && recapData"> - {{ recapData.total_days }} hari sekolah</span>
      </Text>
    </div>

    <Tabs v-model:activeKey="activeTab" @change="handleTabChange">
      <!-- ==================== TAB HARIAN ==================== -->
      <TabPane key="daily">
        <template #tab>
          <span><UnorderedListOutlined /> Harian</span>
        </template>

        <!-- Daily Summary Cards -->
        <Row :gutter="[24, 24]" class="summary-row">
          <Col :xs="24" :sm="12" :lg="6">
            <Card class="stat-card">
              <Statistic title="Total Siswa" :value="dailyStats.totalStudents" :value-style="{ color: '#3b82f6' }" />
            </Card>
          </Col>
          <Col :xs="24" :sm="12" :lg="6">
            <Card class="stat-card">
              <Statistic title="Hadir" :value="dailyStats.totalPresent" :value-style="{ color: '#22c55e' }">
                <template #prefix><CheckCircleOutlined /></template>
              </Statistic>
            </Card>
          </Col>
          <Col :xs="24" :sm="12" :lg="6">
            <Card class="stat-card">
              <Statistic title="Terlambat" :value="dailyStats.totalLate" :value-style="{ color: '#f97316' }">
                <template #prefix><ClockCircleOutlined /></template>
              </Statistic>
            </Card>
          </Col>
          <Col :xs="24" :sm="12" :lg="6">
            <Card class="stat-card">
              <Statistic title="Tidak Hadir" :value="dailyStats.totalAbsent" :value-style="{ color: '#ef4444' }">
                <template #prefix><CloseCircleOutlined /></template>
              </Statistic>
            </Card>
          </Col>
        </Row>

        <Card style="margin-top: 24px">
          <!-- Daily Toolbar -->
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
                <Button @click="handleRefresh">
                  <template #icon><ReloadOutlined /></template>
                  Refresh
                </Button>
                <Button type="primary" @click="openExportModal">
                  <template #icon><DownloadOutlined /></template>
                  Export
                </Button>
              </Space>
            </Col>
          </Row>

          <!-- Daily Table -->
          <Table
            :columns="dailyColumns"
            :data-source="filteredDailyAttendance"
            :loading="dailyLoading"
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
                <Tag v-if="(record as AttendanceSummary).late > 0" color="warning">{{ (record as AttendanceSummary).late }}</Tag>
                <span v-else>0</span>
              </template>
              <template v-else-if="column.key === 'absent'">
                <Tag v-if="(record as AttendanceSummary).absent > 0" color="error">{{ (record as AttendanceSummary).absent }}</Tag>
                <span v-else>0</span>
              </template>
              <template v-else-if="column.key === 'percentage'">
                <div class="percentage-cell">
                  <Progress
                    :percent="Math.round(((record as AttendanceSummary).present / (record as AttendanceSummary).totalStudents) * 100)"
                    :stroke-color="getPercentageColor(((record as AttendanceSummary).present / (record as AttendanceSummary).totalStudents) * 100)"
                    :show-info="true"
                    size="small"
                  />
                </div>
              </template>
            </template>
            
            <!-- Summary Footer -->
            <template #summary>
              <Table.Summary fixed>
                <Table.Summary.Row class="summary-footer">
                  <Table.Summary.Cell :index="0"><Text strong>Total</Text></Table.Summary.Cell>
                  <Table.Summary.Cell :index="1" align="center"><Text strong>{{ dailyStats.totalStudents }}</Text></Table.Summary.Cell>
                  <Table.Summary.Cell :index="2" align="center"><Tag color="success">{{ dailyStats.totalPresent }}</Tag></Table.Summary.Cell>
                  <Table.Summary.Cell :index="3" align="center">
                    <Tag v-if="dailyStats.totalLate > 0" color="warning">{{ dailyStats.totalLate }}</Tag>
                    <span v-else>0</span>
                  </Table.Summary.Cell>
                  <Table.Summary.Cell :index="4" align="center">
                    <Tag v-if="dailyStats.totalAbsent > 0" color="error">{{ dailyStats.totalAbsent }}</Tag>
                    <span v-else>0</span>
                  </Table.Summary.Cell>
                  <Table.Summary.Cell :index="5">
                    <div class="percentage-cell">
                      <Progress :percent="dailyStats.percentage" :stroke-color="getPercentageColor(dailyStats.percentage)" :show-info="true" size="small" />
                    </div>
                  </Table.Summary.Cell>
                </Table.Summary.Row>
              </Table.Summary>
            </template>
          </Table>
        </Card>
      </TabPane>

      <!-- ==================== TAB BULANAN ==================== -->
      <TabPane key="monthly">
        <template #tab>
          <span><BarChartOutlined /> Rekap Bulanan</span>
        </template>

        <!-- Monthly Summary Cards -->
        <Row :gutter="[24, 24]" class="summary-row">
          <Col :xs="24" :sm="12" :lg="6">
            <Card class="stat-card">
              <Statistic title="Total Siswa" :value="monthlyStats.totalStudents" :value-style="{ color: '#3b82f6' }" />
            </Card>
          </Col>
          <Col :xs="24" :sm="12" :lg="6">
            <Card class="stat-card">
              <Statistic title="Rata-rata Hadir" :value="monthlyStats.avgPresent" suffix="hari" :value-style="{ color: '#22c55e' }">
                <template #prefix><CheckCircleOutlined /></template>
              </Statistic>
            </Card>
          </Col>
          <Col :xs="24" :sm="12" :lg="6">
            <Card class="stat-card">
              <Statistic title="Rata-rata Terlambat" :value="monthlyStats.avgLate" suffix="hari" :value-style="{ color: '#f97316' }">
                <template #prefix><ClockCircleOutlined /></template>
              </Statistic>
            </Card>
          </Col>
          <Col :xs="24" :sm="12" :lg="6">
            <Card class="stat-card">
              <Statistic title="Rata-rata Kehadiran" :value="monthlyStats.avgPercentage" suffix="%" :value-style="{ color: getPercentageColor(monthlyStats.avgPercentage) }">
                <template #prefix><CloseCircleOutlined /></template>
              </Statistic>
            </Card>
          </Col>
        </Row>

        <Card style="margin-top: 24px">
          <!-- Monthly Toolbar -->
          <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
            <Col :xs="24" :sm="24" :md="16">
              <Space wrap>
                <Select v-model:value="selectedMonth" style="width: 150px" @change="handleMonthYearChange">
                  <SelectOption v-for="month in MONTH_OPTIONS" :key="month.value" :value="month.value">
                    {{ month.label }}
                  </SelectOption>
                </Select>
                <Select v-model:value="selectedYear" style="width: 120px" @change="handleMonthYearChange">
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
                <Button @click="handleRefresh">
                  <template #icon><ReloadOutlined /></template>
                  Refresh
                </Button>
                <Button type="primary" :loading="exporting" @click="openExportModal">
                  <template #icon><DownloadOutlined /></template>
                  Export
                </Button>
              </Space>
            </Col>
          </Row>

          <!-- Monthly Table -->
          <Table
            :columns="monthlyColumns"
            :data-source="recapData?.student_recaps || []"
            :loading="monthlyLoading"
            :pagination="{ pageSize: 20, showSizeChanger: true, showTotal: (total: number) => `Total ${total} siswa` }"
            row-key="student_id"
            :scroll="{ x: 1000 }"
          >
            <template #bodyCell="{ column, record, index }">
              <template v-if="column.key === 'index'">{{ index + 1 }}</template>
              <template v-else-if="column.key === 'total_present'">
                <Tag color="success">{{ (record as StudentRecapSummary).total_present }}</Tag>
              </template>
              <template v-else-if="column.key === 'total_late'">
                <Tag v-if="(record as StudentRecapSummary).total_late > 0" color="warning">{{ (record as StudentRecapSummary).total_late }}</Tag>
                <span v-else>0</span>
              </template>
              <template v-else-if="column.key === 'total_very_late'">
                <Tag v-if="(record as StudentRecapSummary).total_very_late > 0" color="orange">{{ (record as StudentRecapSummary).total_very_late }}</Tag>
                <span v-else>0</span>
              </template>
              <template v-else-if="column.key === 'total_absent'">
                <Tag v-if="(record as StudentRecapSummary).total_absent > 0" color="error">{{ (record as StudentRecapSummary).total_absent }}</Tag>
                <span v-else>0</span>
              </template>
              <template v-else-if="column.key === 'percentage'">
                <div class="percentage-cell">
                  <Progress
                    :percent="Math.round((record as StudentRecapSummary).attendance_percent)"
                    :stroke-color="getPercentageColor((record as StudentRecapSummary).attendance_percent)"
                    :show-info="true"
                    size="small"
                  />
                </div>
              </template>
            </template>

            <template #emptyText>
              <Empty description="Tidak ada data rekap untuk periode ini" />
            </template>
          </Table>
        </Card>
      </TabPane>
    </Tabs>

    <!-- Export Modal -->
    <Modal
      v-model:open="exportModalVisible"
      :title="activeTab === 'daily' ? 'Export Absensi Harian' : 'Export Rekap Bulanan'"
      :confirm-loading="exporting"
      ok-text="Export"
      cancel-text="Batal"
      @ok="handleExport"
    >
      <div class="export-form">
        <div v-if="activeTab === 'daily'" class="form-item">
          <label>Rentang Tanggal <span class="required">*</span></label>
          <RangePicker
            v-model:value="exportDateRange"
            format="DD MMMM YYYY"
            :placeholder="['Tanggal Mulai', 'Tanggal Akhir']"
            style="width: 100%"
          />
        </div>
        <div v-else class="form-item">
          <label>Periode</label>
          <Text>{{ formattedMonthlyPeriod }}</Text>
        </div>
        <div class="form-item">
          <label>Kelas (Opsional)</label>
          <Select
            v-model:value="exportClassId"
            placeholder="Semua Kelas"
            allow-clear
            style="width: 100%"
            :loading="loadingClasses"
          >
            <SelectOption v-for="cls in classes" :key="cls.id" :value="cls.id">
              {{ cls.name }}
            </SelectOption>
          </Select>
        </div>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
.attendance-report {
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

:deep(.summary-footer) {
  background: #fafafa;
}

.export-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.export-form .form-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.export-form .form-item label {
  font-weight: 500;
  color: #333;
}

.export-form .form-item .required {
  color: #ff4d4f;
}

@media (max-width: 768px) {
  .toolbar-right {
    margin-top: 16px;
    justify-content: flex-start;
  }
}
</style>
