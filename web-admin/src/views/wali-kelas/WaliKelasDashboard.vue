<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import {
  Row,
  Col,
  Card,
  Statistic,
  Spin,
  Alert,
  List,
  ListItem,
  ListItemMeta,
  Typography,
  Tag,
  Avatar,
  Button,
  Empty,
  Table,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  TeamOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  ClockCircleOutlined,
  BookOutlined,
  FormOutlined,
  UserOutlined,
  RightOutlined,
  CalendarOutlined,
} from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'
import { homeroomService } from '@/services'
import type { HomeroomStats, Grade, HomeroomNote, ClassStudent } from '@/types/homeroom'

const { Title, Text } = Typography
const router = useRouter()

const loading = ref(true)
const error = ref<string | null>(null)

const stats = ref<HomeroomStats>({
  classId: 0,
  className: '',
  totalStudents: 0,
  todayAttendance: {
    present: 0,
    absent: 0,
    late: 0,
    excused: 0,
  },
  recentGrades: [],
  recentNotes: [],
})

const students = ref<ClassStudent[]>([])

// Mock data for development
const loadMockData = () => {
  stats.value = {
    classId: 1,
    className: 'VII-A',
    totalStudents: 30,
    todayAttendance: {
      present: 26,
      absent: 2,
      late: 2,
      excused: 0,
    },
    recentGrades: [
      { id: 1, studentId: 1, studentName: 'Ahmad Fauzi', studentNis: '2024001', title: 'Ulangan Matematika', score: 85, createdBy: 1, createdAt: new Date().toISOString(), updatedAt: new Date().toISOString() },
      { id: 2, studentId: 2, studentName: 'Budi Santoso', studentNis: '2024002', title: 'Ulangan Matematika', score: 78, createdBy: 1, createdAt: new Date().toISOString(), updatedAt: new Date().toISOString() },
      { id: 3, studentId: 3, studentName: 'Citra Dewi', studentNis: '2024003', title: 'Tugas IPA', score: 90, createdBy: 1, createdAt: new Date(Date.now() - 86400000).toISOString(), updatedAt: new Date(Date.now() - 86400000).toISOString() },
    ],
    recentNotes: [
      { id: 1, studentId: 1, studentName: 'Ahmad Fauzi', studentNis: '2024001', teacherId: 1, content: 'Siswa menunjukkan peningkatan dalam partisipasi kelas.', createdAt: new Date().toISOString(), updatedAt: new Date().toISOString() },
      { id: 2, studentId: 4, studentName: 'Dian Pratama', studentNis: '2024004', teacherId: 1, content: 'Perlu perhatian lebih dalam mata pelajaran Bahasa Inggris.', createdAt: new Date(Date.now() - 86400000).toISOString(), updatedAt: new Date(Date.now() - 86400000).toISOString() },
    ],
  }

  students.value = [
    { id: 1, nis: '2024001', nisn: '0012345678', name: 'Ahmad Fauzi', isActive: true },
    { id: 2, nis: '2024002', nisn: '0012345679', name: 'Budi Santoso', isActive: true },
    { id: 3, nis: '2024003', nisn: '0012345680', name: 'Citra Dewi', isActive: true },
    { id: 4, nis: '2024004', nisn: '0012345681', name: 'Dian Pratama', isActive: true },
    { id: 5, nis: '2024005', nisn: '0012345682', name: 'Eka Putri', isActive: true },
  ]
}

const loadData = async () => {
  loading.value = true
  error.value = null

  try {
    const [statsResponse, studentsResponse] = await Promise.all([
      homeroomService.getStats(),
      homeroomService.getClassStudents({ pageSize: 10 }),
    ])
    stats.value = statsResponse
    students.value = studentsResponse.data
  } catch {
    loadMockData()
  } finally {
    loading.value = false
  }
}

// Format date
const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
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

// Get score color
const getScoreColor = (score: number): string => {
  if (score >= 85) return '#22c55e'
  if (score >= 70) return '#f97316'
  return '#ef4444'
}

// Student table columns
const studentColumns: TableProps['columns'] = [
  {
    title: 'NIS',
    dataIndex: 'nis',
    key: 'nis',
    width: 100,
  },
  {
    title: 'Nama',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: 'Status',
    key: 'status',
    width: 100,
    align: 'center',
  },
]

// Navigate to attendance page
const goToAttendance = () => {
  router.push('/homeroom/attendance')
}

// Navigate to grades page
const goToGrades = () => {
  router.push('/homeroom/grades')
}

// Navigate to notes page
const goToNotes = () => {
  router.push('/homeroom/notes')
}

// Navigate to student profile
const goToStudentProfile = (studentId: number) => {
  router.push(`/bk/students/${studentId}`)
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="wali-kelas-dashboard">
    <div class="page-header">
      <div>
        <Title :level="2" style="margin: 0">Dashboard Wali Kelas</Title>
        <Text type="secondary">
          Kelas {{ stats.className }}
        </Text>
      </div>
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
              :value-style="{ color: '#3b82f6' }"
            >
              <template #prefix>
                <TeamOutlined />
              </template>
            </Statistic>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :lg="6">
          <Card class="stat-card" hoverable @click="goToAttendance">
            <Statistic
              title="Hadir Hari Ini"
              :value="stats.todayAttendance.present"
              :suffix="`/ ${stats.totalStudents}`"
              :value-style="{ color: '#22c55e' }"
            >
              <template #prefix>
                <CheckCircleOutlined />
              </template>
            </Statistic>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :lg="6">
          <Card class="stat-card" hoverable @click="goToAttendance">
            <Statistic
              title="Terlambat"
              :value="stats.todayAttendance.late"
              :value-style="{ color: '#f97316' }"
            >
              <template #prefix>
                <ClockCircleOutlined />
              </template>
            </Statistic>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :lg="6">
          <Card class="stat-card" hoverable @click="goToAttendance">
            <Statistic
              title="Tidak Hadir"
              :value="stats.todayAttendance.absent"
              :value-style="{ color: '#ef4444' }"
            >
              <template #prefix>
                <CloseCircleOutlined />
              </template>
            </Statistic>
          </Card>
        </Col>
      </Row>

      <!-- Students List & Recent Activities -->
      <Row :gutter="[24, 24]" style="margin-top: 24px">
        <Col :xs="24" :lg="12">
          <Card title="Daftar Siswa" class="students-card">
            <template #extra>
              <Text type="secondary">{{ stats.className }}</Text>
            </template>
            <Table
              :columns="studentColumns"
              :data-source="students"
              :pagination="false"
              :loading="loading"
              row-key="id"
              size="small"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'status'">
                  <Tag v-if="(record as ClassStudent).isActive" color="success">Aktif</Tag>
                  <Tag v-else color="default">Tidak Aktif</Tag>
                </template>
              </template>
            </Table>
            <div v-if="students.length > 5" class="view-all">
              <Button type="link" size="small">
                Lihat Semua Siswa <RightOutlined />
              </Button>
            </div>
          </Card>
        </Col>

        <Col :xs="24" :lg="12">
          <Card title="Nilai Terbaru" class="recent-card">
            <template #extra>
              <Button type="link" size="small" @click="goToGrades">
                Lihat Semua <RightOutlined />
              </Button>
            </template>
            <List
              v-if="stats.recentGrades.length > 0"
              :data-source="stats.recentGrades"
              :loading="loading"
              size="small"
            >
              <template #renderItem="{ item }">
                <ListItem class="recent-item" @click="goToStudentProfile((item as Grade).studentId)">
                  <ListItemMeta>
                    <template #avatar>
                      <Avatar :style="{ backgroundColor: getScoreColor((item as Grade).score) }" size="small">
                        {{ (item as Grade).score }}
                      </Avatar>
                    </template>
                    <template #title>
                      <span>{{ (item as Grade).studentName }}</span>
                    </template>
                    <template #description>
                      <div>{{ (item as Grade).title }}</div>
                      <Text type="secondary" class="date-text">{{ formatDate((item as Grade).createdAt) }}</Text>
                    </template>
                  </ListItemMeta>
                </ListItem>
              </template>
            </List>
            <Empty v-else description="Belum ada nilai tercatat" />
          </Card>
        </Col>
      </Row>

      <!-- Recent Notes -->
      <Row :gutter="[24, 24]" style="margin-top: 24px">
        <Col :xs="24">
          <Card title="Catatan Wali Kelas Terbaru" class="notes-card">
            <template #extra>
              <Button type="link" size="small" @click="goToNotes">
                Lihat Semua <RightOutlined />
              </Button>
            </template>
            <List
              v-if="stats.recentNotes.length > 0"
              :data-source="stats.recentNotes"
              :loading="loading"
            >
              <template #renderItem="{ item }">
                <ListItem class="note-item" @click="goToStudentProfile((item as HomeroomNote).studentId)">
                  <ListItemMeta>
                    <template #avatar>
                      <Avatar :style="{ backgroundColor: '#8b5cf6' }">
                        <template #icon><UserOutlined /></template>
                      </Avatar>
                    </template>
                    <template #title>
                      <span class="student-name">{{ (item as HomeroomNote).studentName }}</span>
                      <Text type="secondary" style="margin-left: 8px; font-size: 12px">
                        {{ formatDate((item as HomeroomNote).createdAt) }}
                      </Text>
                    </template>
                    <template #description>
                      <div class="note-content">{{ (item as HomeroomNote).content }}</div>
                    </template>
                  </ListItemMeta>
                </ListItem>
              </template>
            </List>
            <Empty v-else description="Belum ada catatan wali kelas" />
          </Card>
        </Col>
      </Row>

      <!-- Quick Actions -->
      <Row :gutter="[24, 24]" style="margin-top: 24px">
        <Col :xs="24">
          <Card title="Aksi Cepat">
            <Row :gutter="[16, 16]">
              <Col :xs="24" :sm="8">
                <Button type="primary" block size="large" @click="goToAttendance">
                  <template #icon><CheckCircleOutlined /></template>
                  Input Absensi Manual
                </Button>
              </Col>
              <Col :xs="24" :sm="8">
                <Button block size="large" @click="goToGrades">
                  <template #icon><BookOutlined /></template>
                  Input Nilai
                </Button>
              </Col>
              <Col :xs="24" :sm="8">
                <Button block size="large" @click="goToNotes">
                  <template #icon><FormOutlined /></template>
                  Tulis Catatan
                </Button>
              </Col>
            </Row>
          </Card>
        </Col>
      </Row>
    </Spin>
  </div>
</template>

<style scoped>
.wali-kelas-dashboard {
  padding: 0;
}

.page-header {
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  flex-wrap: wrap;
  gap: 8px;
}

.stat-card {
  height: 100%;
  cursor: pointer;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-card :deep(.ant-statistic-title) {
  font-size: 14px;
  color: #8c8c8c;
}

.stat-card :deep(.ant-statistic-content-prefix) {
  margin-right: 8px;
}

.students-card,
.recent-card,
.notes-card {
  height: 100%;
}

.view-all {
  text-align: center;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

.recent-item,
.note-item {
  cursor: pointer;
  transition: background-color 0.2s;
}

.recent-item:hover,
.note-item:hover {
  background-color: #fafafa;
}

.student-name {
  font-weight: 500;
}

.note-content {
  color: #595959;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.date-text {
  font-size: 12px;
  display: block;
  margin-top: 4px;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
