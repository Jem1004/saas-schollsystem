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
  Tag,
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
                  :size="120"
                >
                  <template #format="percent">
                    <span class="progress-text">{{ percent }}%</span>
                  </template>
                </Progress>
                <Text type="secondary" style="margin-top: 8px">Kehadiran</Text>
              </div>
              <div class="attendance-details">
                <div class="attendance-item">
                  <CheckCircleOutlined class="icon success" />
                  <div class="item-content">
                    <Text strong>{{ stats.todayAttendance.present }}</Text>
                    <Text type="secondary">Hadir</Text>
                  </div>
                </div>
                <div class="attendance-item">
                  <ClockCircleOutlined class="icon warning" />
                  <div class="item-content">
                    <Text strong>{{ stats.todayAttendance.late }}</Text>
                    <Text type="secondary">Terlambat</Text>
                  </div>
                </div>
                <div class="attendance-item">
                  <CloseCircleOutlined class="icon danger" />
                  <div class="item-content">
                    <Text strong>{{ stats.todayAttendance.absent }}</Text>
                    <Text type="secondary">Tidak Hadir</Text>
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
              size="small"
            >
              <template #renderItem="{ item }">
                <ListItem>
                  <ListItemMeta :title="item.className">
                    <template #description>
                      <div class="class-attendance-stats">
                        <Tag color="success">{{ item.present }} hadir</Tag>
                        <Tag v-if="item.late > 0" color="warning">{{ item.late }} terlambat</Tag>
                        <Tag v-if="item.absent > 0" color="error">{{ item.absent }} absen</Tag>
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
                  <Card size="small" class="class-item-card">
                    <div class="class-info">
                      <BookOutlined class="class-icon" />
                      <div class="class-details">
                        <Text strong>{{ item.name }}</Text>
                        <Text type="secondary" class="class-meta">
                          {{ item.studentCount || 0 }} siswa
                        </Text>
                      </div>
                    </div>
                    <div class="class-teacher" v-if="item.homeroomTeacherName">
                      <UserOutlined />
                      <Text type="secondary">{{ item.homeroomTeacherName }}</Text>
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
}

.stat-card :deep(.ant-statistic-title) {
  font-size: 14px;
  color: #8c8c8c;
}

.stat-card :deep(.ant-statistic-content-prefix) {
  margin-right: 8px;
}

.attendance-card,
.class-attendance-card,
.classes-card {
  height: 100%;
}

.attendance-overview {
  display: flex;
  align-items: center;
  gap: 32px;
}

.attendance-progress {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.progress-text {
  font-size: 24px;
  font-weight: 600;
}

.attendance-details {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.attendance-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.attendance-item .icon {
  font-size: 24px;
}

.attendance-item .icon.success {
  color: #22c55e;
}

.attendance-item .icon.warning {
  color: #f97316;
}

.attendance-item .icon.danger {
  color: #ef4444;
}

.item-content {
  display: flex;
  flex-direction: column;
}

.class-attendance-stats {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.class-item-card {
  height: 100%;
}

.class-info {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.class-icon {
  font-size: 24px;
  color: #f97316;
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
  font-size: 12px;
  color: #8c8c8c;
  padding-top: 8px;
  border-top: 1px solid #f0f0f0;
}

@media (max-width: 768px) {
  .attendance-overview {
    flex-direction: column;
    gap: 24px;
  }

  .attendance-details {
    flex-direction: row;
    justify-content: space-around;
    width: 100%;
  }

  .attendance-item {
    flex-direction: column;
    text-align: center;
  }
}
</style>
