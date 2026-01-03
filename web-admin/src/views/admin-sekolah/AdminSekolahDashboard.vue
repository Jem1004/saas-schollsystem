<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import {
  Row,
  Col,
  Card,
  Statistic,
  Spin,
  Alert,
  Progress,
  List,
  ListItem,
  ListItemMeta,
  Typography,
} from 'ant-design-vue'
import {
  TeamOutlined,
  BookOutlined,
  UserOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  ClockCircleOutlined,
  CalendarOutlined,
} from '@ant-design/icons-vue'
import { schoolService } from '@/services'
import type { SchoolStats, AttendanceSummary, Class } from '@/types/school'

const { Title, Text } = Typography

const loading = ref(true)
const error = ref<string | null>(null)

const stats = ref<SchoolStats>({
  totalStudents: 0,
  totalClasses: 0,
  totalTeachers: 0,
  totalParents: 0,
  todayAttendance: {
    present: 0,
    absent: 0,
    late: 0,
    total: 0,
  },
})

const recentClasses = ref<Class[]>([])
const todayAttendanceByClass = ref<AttendanceSummary[]>([])

// Computed attendance percentage
const attendancePercentage = computed(() => {
  const { present, total } = stats.value.todayAttendance
  if (total === 0) return 0
  return Math.round((present / total) * 100)
})

const loadData = async () => {
  loading.value = true
  error.value = null

  try {
    const [statsRes, classesRes, attendanceRes] = await Promise.allSettled([
      schoolService.getStats(),
      schoolService.getClasses({ page: 1, page_size: 5 }),
      schoolService.getAttendanceSummary({ date: new Date().toISOString().split('T')[0] }),
    ])

    if (statsRes.status === 'fulfilled') {
      stats.value = statsRes.value
    } else {
      error.value = 'Gagal memuat statistik sekolah'
    }
    
    if (classesRes.status === 'fulfilled') {
      recentClasses.value = classesRes.value.classes
    }
    
    if (attendanceRes.status === 'fulfilled') {
      todayAttendanceByClass.value = attendanceRes.value.data
    } else if (classesRes.status === 'fulfilled') {
      // Build attendance by class from classes data if attendance API fails
      todayAttendanceByClass.value = classesRes.value.classes.map((cls: Class) => ({
        date: new Date().toISOString().split('T')[0],
        classId: cls.id,
        className: cls.name,
        totalStudents: cls.studentCount || 0,
        present: 0,
        absent: cls.studentCount || 0,
        late: 0,
        sick: 0,
        excused: 0,
      }))
    }
  } catch (err) {
    error.value = 'Gagal memuat data dashboard. Silakan coba lagi.'
    console.error('Dashboard load error:', err)
  } finally {
    loading.value = false
  }
}

// Get attendance status color
const getAttendanceStatusColor = (summary: AttendanceSummary): string => {
  const percentage = (summary.present / summary.totalStudents) * 100
  if (percentage >= 95) return '#22c55e'
  if (percentage >= 85) return '#f97316'
  return '#ef4444'
}

// Format today's date
const todayFormatted = computed(() => {
  return new Date().toLocaleDateString('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
})

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="admin-sekolah-dashboard">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Dashboard Admin Sekolah</Title>
      <Text type="secondary">
        <CalendarOutlined /> {{ todayFormatted }}
      </Text>
    </div>

    <Spin :spinning="loading">
      <Alert
        v-if="error"
        type="error"
        :message="error"
        show-icon
        closable
        style="margin-bottom: 24px"
      />

      <!-- Statistics Cards -->
      <Row :gutter="[24, 24]" class="stats-row">
        <Col :xs="24" :sm="12" :lg="6">
          <Card class="stat-card">
            <Statistic
              title="Total Siswa"
              :value="stats.totalStudents"
              :value-style="{ color: '#f97316' }"
            >
              <template #prefix>
                <TeamOutlined />
              </template>
            </Statistic>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :lg="6">
          <Card class="stat-card">
            <Statistic
              title="Total Kelas"
              :value="stats.totalClasses"
              :value-style="{ color: '#3b82f6' }"
            >
              <template #prefix>
                <BookOutlined />
              </template>
            </Statistic>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :lg="6">
          <Card class="stat-card">
            <Statistic
              title="Total Guru"
              :value="stats.totalTeachers"
              :value-style="{ color: '#8b5cf6' }"
            >
              <template #prefix>
                <UserOutlined />
              </template>
            </Statistic>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :lg="6">
          <Card class="stat-card">
            <Statistic
              title="Total Orang Tua"
              :value="stats.totalParents"
              :value-style="{ color: '#22c55e' }"
            >
              <template #prefix>
                <UserOutlined />
              </template>
            </Statistic>
          </Card>
        </Col>
      </Row>

      <!-- Today's Attendance Overview -->
      <Row :gutter="[24, 24]" style="margin-top: 24px">
        <Col :xs="24" :lg="12">
          <Card title="Absensi Hari Ini" class="attendance-card">
            <div class="attendance-overview">
              <div class="attendance-progress">
                <Progress
                  type="circle"
                  :percent="attendancePercentage"
                  :stroke-color="attendancePercentage >= 95 ? '#22c55e' : attendancePercentage >= 85 ? '#f97316' : '#ef4444'"
                  :size="140"
                  :stroke-width="8"
                >
                  <template #format="percent">
                    <div class="progress-content">
                      <span class="progress-text">{{ percent }}%</span>
                      <span class="progress-label">hadir</span>
                    </div>
                  </template>
                </Progress>
              </div>
              <div class="attendance-details">
                <div class="attendance-item">
                  <div class="icon-wrapper success">
                    <CheckCircleOutlined />
                  </div>
                  <div class="item-content">
                    <Text strong class="stat-value">{{ stats.todayAttendance.present }}</Text>
                    <Text type="secondary" class="stat-label">Hadir</Text>
                  </div>
                </div>
                <div class="attendance-item">
                  <div class="icon-wrapper warning">
                    <ClockCircleOutlined />
                  </div>
                  <div class="item-content">
                    <Text strong class="stat-value">{{ stats.todayAttendance.late }}</Text>
                    <Text type="secondary" class="stat-label">Terlambat</Text>
                  </div>
                </div>
                <div class="attendance-item">
                  <div class="icon-wrapper danger">
                    <CloseCircleOutlined />
                  </div>
                  <div class="item-content">
                    <Text strong class="stat-value">{{ stats.todayAttendance.absent }}</Text>
                    <Text type="secondary" class="stat-label">Absen</Text>
                  </div>
                </div>
              </div>
            </div>
          </Card>
        </Col>
        <Col :xs="24" :lg="12">
          <Card title="Absensi Per Kelas" class="class-attendance-card">
            <List
              :data-source="todayAttendanceByClass"
              :loading="loading"
              size="large"
              item-layout="horizontal"
            >
              <template #renderItem="{ item }">
                <ListItem class="attendance-list-item">
                  <ListItemMeta>
                    <template #title>
                      <span class="class-name">{{ item.className }}</span>
                    </template>
                    <template #description>
                      <div class="class-attendance-stats">
                        <span class="stat-pill success">{{ item.present }} hadir</span>
                        <span v-if="item.late > 0" class="stat-pill warning">{{ item.late }} terlambat</span>
                        <span v-if="item.absent > 0" class="stat-pill error">{{ item.absent }} absen</span>
                      </div>
                    </template>
                  </ListItemMeta>
                  <template #actions>
                    <Progress
                      :percent="Math.round((item.present / item.totalStudents) * 100)"
                      :stroke-color="getAttendanceStatusColor(item)"
                      :show-info="true"
                      size="small"
                      style="width: 100px"
                    />
                  </template>
                </ListItem>
              </template>
            </List>
          </Card>
        </Col>
      </Row>

      <!-- Recent Classes -->
      <Row :gutter="[24, 24]" style="margin-top: 24px">
        <Col :xs="24">
          <Card title="Daftar Kelas" class="classes-card">
            <List
              :data-source="recentClasses"
              :loading="loading"
              :grid="{ gutter: 16, xs: 1, sm: 2, md: 3, lg: 4, xl: 5 }"
            >
              <template #renderItem="{ item }">
                <ListItem>
                  <Card size="small" class="class-item-card" :bordered="false">
                    <div class="class-info">
                      <div class="class-icon-wrapper">
                         <BookOutlined />
                      </div>
                      <div class="class-details">
                        <Text strong style="font-size: 16px;">{{ item.name }}</Text>
                        <Text type="secondary" class="class-meta">
                          {{ item.studentCount || 0 }} siswa
                        </Text>
                      </div>
                    </div>
                    <div class="class-teacher" v-if="item.homeroomTeacherName">
                      <UserOutlined style="color: #64748b;" />
                      <Text type="secondary" style="font-size: 13px;">{{ item.homeroomTeacherName }}</Text>
                    </div>
                    <div class="class-teacher empty" v-else>
                      <Text type="secondary" style="font-size: 13px; font-style: italic;">Belum ada wali kelas</Text>
                    </div>
                  </Card>
                </ListItem>
              </template>
            </List>
          </Card>
        </Col>
      </Row>
    </Spin>
  </div>
</template>

<style scoped>
.admin-sekolah-dashboard {
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

.stat-card {
  height: 100%;
  border: 1px solid #f1f5f9;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
}

.stat-card :deep(.ant-statistic-title) {
  font-size: 13px;
  color: #64748b;
  margin-bottom: 8px;
}

.stat-card :deep(.ant-statistic-content) {
  font-weight: 600;
}

.stat-card :deep(.ant-statistic-content-prefix) {
  margin-right: 12px;
}

/* Attendance Card */
.attendance-card,
.class-attendance-card,
.classes-card {
  height: 100%;
  border: 1px solid #f1f5f9;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
}

.attendance-overview {
  display: flex;
  align-items: center;
  gap: 48px;
  padding: 16px 0;
}

.attendance-progress {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.progress-content {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
}

.progress-text {
  font-size: 28px;
  font-weight: 700;
  color: #0f172a;
}

.progress-label {
  font-size: 12px;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.attendance-details {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.attendance-item {
  display: flex;
  align-items: center;
  gap: 16px;
}

.icon-wrapper {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
}

.icon-wrapper.success { background: #ecfdf5; color: #10b981; }
.icon-wrapper.warning { background: #fffbeb; color: #f59e0b; }
.icon-wrapper.danger { background: #fef2f2; color: #ef4444; }

.item-content {
  display: flex;
  flex-direction: column;
}

.stat-value { font-size: 18px; line-height: 1.2; }
.stat-label { font-size: 13px; }

/* Class Attendance List */
.attendance-list-item {
  padding: 16px;
  border-bottom: 1px solid #f8fafc;
}

.class-name {
  font-weight: 600;
  color: #334155;
}

.class-attendance-stats {
  display: flex;
  gap: 8px;
  margin-top: 4px;
}

.stat-pill {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}
.stat-pill.success { background: #f0fdf4; color: #166534; }
.stat-pill.warning { background: #fff7ed; color: #9a3412; }
.stat-pill.error { background: #fef2f2; color: #991b1b; }

/* Class Grid Items */
.class-item-card {
  height: 100%;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  transition: all 0.2s;
}

.class-item-card:hover {
  border-color: #cbd5e1;
  background: #ffffff;
  transform: translateY(-2px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.class-info {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.class-icon-wrapper {
  width: 40px;
  height: 40px;
  background: #fff7ed;
  color: #f97316;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.class-details {
  display: flex;
  flex-direction: column;
}

.class-meta {
  font-size: 12px;
}

.class-teacher {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-top: 12px;
  border-top: 1px solid #e2e8f0;
}

@media (max-width: 768px) {
  .attendance-overview {
    flex-direction: column;
    gap: 32px;
  }

  .attendance-details {
    flex-direction: row;
    justify-content: space-between;
    width: 100%;
  }

  .attendance-item {
    flex-direction: column;
    text-align: center;
    gap: 8px;
  }
}
</style>
