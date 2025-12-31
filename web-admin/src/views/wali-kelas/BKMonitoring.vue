<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import {
  Card,
  Row,
  Col,
  Typography,
  Tabs,
  TabPane,
  Table,
  Tag,
  Avatar,
  Empty,
  Select,
  SelectOption,
  Badge,
  Drawer,
  Descriptions,
  DescriptionsItem,
  Spin,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  WarningOutlined,
  TrophyOutlined,
  FileProtectOutlined,
  UserOutlined,
  EyeOutlined,
  LockOutlined,
} from '@ant-design/icons-vue'
import { bkService } from '@/services'
import { homeroomService } from '@/services'
import type { Violation, Achievement, Permit } from '@/types/bk'
import type { ClassStudent } from '@/types/homeroom'
import { ReadOnlyBanner } from '@/components'

const { Title, Text } = Typography

// State
const loading = ref(false)
const activeTab = ref('violations')
const selectedStudentId = ref<number | undefined>(undefined)
const className = ref('VII-A')

// Data
const students = ref<ClassStudent[]>([])
const violations = ref<Violation[]>([])
const achievements = ref<Achievement[]>([])
const permits = ref<Permit[]>([])

// Drawer state
const drawerVisible = ref(false)
const selectedItem = ref<Violation | Achievement | Permit | null>(null)
const selectedItemType = ref<'violation' | 'achievement' | 'permit'>('violation')

// Mock data for development
const mockStudents: ClassStudent[] = [
  { id: 1, nis: '2024001', nisn: '0012345678', name: 'Ahmad Fauzi', isActive: true },
  { id: 2, nis: '2024002', nisn: '0012345679', name: 'Budi Santoso', isActive: true },
  { id: 3, nis: '2024003', nisn: '0012345680', name: 'Citra Dewi', isActive: true },
  { id: 4, nis: '2024004', nisn: '0012345681', name: 'Dian Pratama', isActive: true },
  { id: 5, nis: '2024005', nisn: '0012345682', name: 'Eka Putri', isActive: true },
]

const mockViolations: Violation[] = [
  { id: 1, studentId: 1, studentName: 'Ahmad Fauzi', studentNis: '2024001', studentClass: 'VII-A', category: 'Keterlambatan', level: 'ringan', description: 'Terlambat 15 menit', createdBy: 1, createdByName: 'Ibu Guru BK', createdAt: new Date().toISOString() },
  { id: 2, studentId: 3, studentName: 'Citra Dewi', studentNis: '2024003', studentClass: 'VII-A', category: 'Bolos', level: 'sedang', description: 'Tidak masuk tanpa keterangan', createdBy: 1, createdByName: 'Ibu Guru BK', createdAt: new Date(Date.now() - 86400000).toISOString() },
  { id: 3, studentId: 2, studentName: 'Budi Santoso', studentNis: '2024002', studentClass: 'VII-A', category: 'Seragam', level: 'ringan', description: 'Tidak memakai dasi', createdBy: 1, createdByName: 'Ibu Guru BK', createdAt: new Date(Date.now() - 172800000).toISOString() },
]

const mockAchievements: Achievement[] = [
  { id: 1, studentId: 4, studentName: 'Dian Pratama', studentNis: '2024004', studentClass: 'VII-A', title: 'Juara 1 Olimpiade Matematika', point: 100, description: 'Tingkat Kabupaten', createdBy: 1, createdByName: 'Ibu Guru BK', createdAt: new Date().toISOString() },
  { id: 2, studentId: 5, studentName: 'Eka Putri', studentNis: '2024005', studentClass: 'VII-A', title: 'Juara 2 Lomba Pidato', point: 75, description: 'Tingkat Sekolah', createdBy: 1, createdByName: 'Ibu Guru BK', createdAt: new Date(Date.now() - 86400000).toISOString() },
  { id: 3, studentId: 1, studentName: 'Ahmad Fauzi', studentNis: '2024001', studentClass: 'VII-A', title: 'Siswa Teladan', point: 50, description: 'Bulan Oktober', createdBy: 1, createdByName: 'Ibu Guru BK', createdAt: new Date(Date.now() - 172800000).toISOString() },
]

const mockPermits: Permit[] = [
  { id: 1, studentId: 2, studentName: 'Budi Santoso', studentNis: '2024002', studentNisn: '0012345679', studentClass: 'VII-A', reason: 'Sakit perut, perlu ke klinik', exitTime: '10:30', returnTime: '11:45', responsibleTeacherId: 1, responsibleTeacherName: 'Ibu Guru BK', createdBy: 1, createdByName: 'Ibu Guru BK', createdAt: new Date().toISOString() },
  { id: 2, studentId: 3, studentName: 'Citra Dewi', studentNis: '2024003', studentNisn: '0012345680', studentClass: 'VII-A', reason: 'Dijemput orang tua untuk keperluan keluarga', exitTime: '12:00', responsibleTeacherId: 1, responsibleTeacherName: 'Ibu Guru BK', createdBy: 1, createdByName: 'Ibu Guru BK', createdAt: new Date(Date.now() - 86400000).toISOString() },
]

// Table columns for violations
const violationColumns: TableProps['columns'] = [
  { title: 'Siswa', dataIndex: 'studentName', key: 'studentName' },
  { title: 'Kategori', dataIndex: 'category', key: 'category' },
  { title: 'Level', dataIndex: 'level', key: 'level', width: 100, align: 'center' },
  { title: 'Tanggal', dataIndex: 'createdAt', key: 'createdAt', width: 150 },
  { title: 'Aksi', key: 'action', width: 80, align: 'center' },
]

// Table columns for achievements
const achievementColumns: TableProps['columns'] = [
  { title: 'Siswa', dataIndex: 'studentName', key: 'studentName' },
  { title: 'Prestasi', dataIndex: 'title', key: 'title' },
  { title: 'Poin', dataIndex: 'point', key: 'point', width: 80, align: 'center' },
  { title: 'Tanggal', dataIndex: 'createdAt', key: 'createdAt', width: 150 },
  { title: 'Aksi', key: 'action', width: 80, align: 'center' },
]

// Table columns for permits
const permitColumns: TableProps['columns'] = [
  { title: 'Siswa', dataIndex: 'studentName', key: 'studentName' },
  { title: 'Alasan', dataIndex: 'reason', key: 'reason', ellipsis: true },
  { title: 'Jam Keluar', dataIndex: 'exitTime', key: 'exitTime', width: 100, align: 'center' },
  { title: 'Jam Kembali', dataIndex: 'returnTime', key: 'returnTime', width: 100, align: 'center' },
  { title: 'Aksi', key: 'action', width: 80, align: 'center' },
]

// Get violation level color
const getViolationLevelColor = (level: string): string => {
  switch (level) {
    case 'ringan': return 'warning'
    case 'sedang': return 'orange'
    case 'berat': return 'error'
    default: return 'default'
  }
}

// Get violation level label
const getViolationLevelLabel = (level: string): string => {
  switch (level) {
    case 'ringan': return 'Ringan'
    case 'sedang': return 'Sedang'
    case 'berat': return 'Berat'
    default: return level
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

// Format full date
const formatFullDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
}

// Filtered data based on selected student
const filteredViolations = computed(() => {
  if (!selectedStudentId.value) return violations.value
  return violations.value.filter(v => v.studentId === selectedStudentId.value)
})

const filteredAchievements = computed(() => {
  if (!selectedStudentId.value) return achievements.value
  return achievements.value.filter(a => a.studentId === selectedStudentId.value)
})

const filteredPermits = computed(() => {
  if (!selectedStudentId.value) return permits.value
  return permits.value.filter(p => p.studentId === selectedStudentId.value)
})

// Summary counts
const summaryCounts = computed(() => ({
  violations: filteredViolations.value.length,
  achievements: filteredAchievements.value.length,
  permits: filteredPermits.value.length,
}))

// Load data
const loadData = async () => {
  loading.value = true
  try {
    const [studentsRes, violationsRes, achievementsRes, permitsRes] = await Promise.all([
      homeroomService.getClassStudents({ pageSize: 100 }),
      bkService.getViolations({ pageSize: 100 }),
      bkService.getAchievements({ pageSize: 100 }),
      bkService.getPermits({ pageSize: 100 }),
    ])
    students.value = studentsRes.data
    violations.value = violationsRes.data
    achievements.value = achievementsRes.data
    permits.value = permitsRes.data
  } catch {
    students.value = mockStudents
    violations.value = mockViolations
    achievements.value = mockAchievements
    permits.value = mockPermits
  } finally {
    loading.value = false
  }
}

// View detail
const viewDetail = (item: Violation | Achievement | Permit, type: 'violation' | 'achievement' | 'permit') => {
  selectedItem.value = item
  selectedItemType.value = type
  drawerVisible.value = true
}

// Close drawer
const closeDrawer = () => {
  drawerVisible.value = false
  selectedItem.value = null
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="bk-monitoring">
    <div class="page-header">
      <div>
        <Title :level="2" style="margin: 0">
          <EyeOutlined style="margin-right: 8px; color: #8c8c8c" />
          Monitoring BK
        </Title>
        <div class="page-subtitle">
          <Text type="secondary">Kelas {{ className }}</Text>
          <Tag color="blue" class="readonly-tag">
            <LockOutlined /> Baca Saja
          </Tag>
        </div>
      </div>
    </div>

    <!-- Read-only Alert -->
    <ReadOnlyBanner
      title="Mode Baca Saja"
      description="Sebagai Wali Kelas, Anda dapat melihat data BK siswa di kelas Anda tetapi tidak dapat mengedit. Untuk perubahan data, silakan hubungi Guru BK."
      type="info"
      tooltip-text="Hanya Guru BK yang dapat menambah atau mengubah data BK"
    />

    <!-- Filter -->
    <Card style="margin-bottom: 24px">
      <Row :gutter="16" align="middle">
        <Col :xs="24" :sm="8">
          <Text strong>Filter Siswa:</Text>
        </Col>
        <Col :xs="24" :sm="16">
          <Select
            v-model:value="selectedStudentId"
            placeholder="Semua Siswa"
            allow-clear
            style="width: 100%; max-width: 300px"
            :loading="loading"
          >
            <SelectOption v-for="student in students" :key="student.id" :value="student.id">
              {{ student.nis }} - {{ student.name }}
            </SelectOption>
          </Select>
        </Col>
      </Row>
    </Card>

    <!-- Summary Cards -->
    <Row :gutter="[24, 24]" style="margin-bottom: 24px">
      <Col :xs="24" :sm="8">
        <Card class="summary-card" :class="{ active: activeTab === 'violations' }" @click="activeTab = 'violations'">
          <div class="summary-content">
            <Avatar :style="{ backgroundColor: '#ef4444' }" size="large">
              <template #icon><WarningOutlined /></template>
            </Avatar>
            <div class="summary-text">
              <Text type="secondary">Pelanggaran</Text>
              <Title :level="3" style="margin: 0">{{ summaryCounts.violations }}</Title>
            </div>
          </div>
        </Card>
      </Col>
      <Col :xs="24" :sm="8">
        <Card class="summary-card" :class="{ active: activeTab === 'achievements' }" @click="activeTab = 'achievements'">
          <div class="summary-content">
            <Avatar :style="{ backgroundColor: '#22c55e' }" size="large">
              <template #icon><TrophyOutlined /></template>
            </Avatar>
            <div class="summary-text">
              <Text type="secondary">Prestasi</Text>
              <Title :level="3" style="margin: 0">{{ summaryCounts.achievements }}</Title>
            </div>
          </div>
        </Card>
      </Col>
      <Col :xs="24" :sm="8">
        <Card class="summary-card" :class="{ active: activeTab === 'permits' }" @click="activeTab = 'permits'">
          <div class="summary-content">
            <Avatar :style="{ backgroundColor: '#3b82f6' }" size="large">
              <template #icon><FileProtectOutlined /></template>
            </Avatar>
            <div class="summary-text">
              <Text type="secondary">Izin Keluar</Text>
              <Title :level="3" style="margin: 0">{{ summaryCounts.permits }}</Title>
            </div>
          </div>
        </Card>
      </Col>
    </Row>

    <!-- Data Tables -->
    <Card>
      <Spin :spinning="loading">
        <Tabs v-model:activeKey="activeTab">
          <!-- Violations Tab -->
          <TabPane key="violations">
            <template #tab>
              <Badge :count="summaryCounts.violations" :offset="[10, 0]">
                <span><WarningOutlined /> Pelanggaran</span>
              </Badge>
            </template>
            <Table
              :columns="violationColumns"
              :data-source="filteredViolations"
              :pagination="{ pageSize: 10 }"
              row-key="id"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'level'">
                  <Tag :color="getViolationLevelColor((record as Violation).level)">
                    {{ getViolationLevelLabel((record as Violation).level) }}
                  </Tag>
                </template>
                <template v-else-if="column.key === 'createdAt'">
                  {{ formatDate((record as Violation).createdAt) }}
                </template>
                <template v-else-if="column.key === 'action'">
                  <a @click="viewDetail(record as Violation, 'violation')">
                    <EyeOutlined /> Lihat
                  </a>
                </template>
              </template>
            </Table>
            <Empty v-if="filteredViolations.length === 0 && !loading" description="Tidak ada data pelanggaran" />
          </TabPane>

          <!-- Achievements Tab -->
          <TabPane key="achievements">
            <template #tab>
              <Badge :count="summaryCounts.achievements" :offset="[10, 0]">
                <span><TrophyOutlined /> Prestasi</span>
              </Badge>
            </template>
            <Table
              :columns="achievementColumns"
              :data-source="filteredAchievements"
              :pagination="{ pageSize: 10 }"
              row-key="id"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'point'">
                  <Tag color="success">+{{ (record as Achievement).point }}</Tag>
                </template>
                <template v-else-if="column.key === 'createdAt'">
                  {{ formatDate((record as Achievement).createdAt) }}
                </template>
                <template v-else-if="column.key === 'action'">
                  <a @click="viewDetail(record as Achievement, 'achievement')">
                    <EyeOutlined /> Lihat
                  </a>
                </template>
              </template>
            </Table>
            <Empty v-if="filteredAchievements.length === 0 && !loading" description="Tidak ada data prestasi" />
          </TabPane>

          <!-- Permits Tab -->
          <TabPane key="permits">
            <template #tab>
              <Badge :count="summaryCounts.permits" :offset="[10, 0]">
                <span><FileProtectOutlined /> Izin Keluar</span>
              </Badge>
            </template>
            <Table
              :columns="permitColumns"
              :data-source="filteredPermits"
              :pagination="{ pageSize: 10 }"
              row-key="id"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'returnTime'">
                  <span v-if="(record as Permit).returnTime">{{ (record as Permit).returnTime }}</span>
                  <Tag v-else color="warning">Belum Kembali</Tag>
                </template>
                <template v-else-if="column.key === 'action'">
                  <a @click="viewDetail(record as Permit, 'permit')">
                    <EyeOutlined /> Lihat
                  </a>
                </template>
              </template>
            </Table>
            <Empty v-if="filteredPermits.length === 0 && !loading" description="Tidak ada data izin keluar" />
          </TabPane>
        </Tabs>
      </Spin>
    </Card>

    <!-- Detail Drawer -->
    <Drawer
      v-model:open="drawerVisible"
      :title="selectedItemType === 'violation' ? 'Detail Pelanggaran' : selectedItemType === 'achievement' ? 'Detail Prestasi' : 'Detail Izin Keluar'"
      width="500"
      @close="closeDrawer"
    >
      <div v-if="selectedItem">
        <!-- Violation Detail -->
        <template v-if="selectedItemType === 'violation'">
          <Descriptions :column="1" bordered size="small">
            <DescriptionsItem label="Siswa">
              <div class="student-info">
                <Avatar :style="{ backgroundColor: '#f97316' }" size="small">
                  <template #icon><UserOutlined /></template>
                </Avatar>
                <span style="margin-left: 8px">{{ (selectedItem as Violation).studentName }}</span>
              </div>
            </DescriptionsItem>
            <DescriptionsItem label="NIS">{{ (selectedItem as Violation).studentNis }}</DescriptionsItem>
            <DescriptionsItem label="Kategori">{{ (selectedItem as Violation).category }}</DescriptionsItem>
            <DescriptionsItem label="Level">
              <Tag :color="getViolationLevelColor((selectedItem as Violation).level)">
                {{ getViolationLevelLabel((selectedItem as Violation).level) }}
              </Tag>
            </DescriptionsItem>
            <DescriptionsItem label="Deskripsi">{{ (selectedItem as Violation).description }}</DescriptionsItem>
            <DescriptionsItem label="Dicatat Oleh">{{ (selectedItem as Violation).createdByName }}</DescriptionsItem>
            <DescriptionsItem label="Tanggal">{{ formatFullDate((selectedItem as Violation).createdAt) }}</DescriptionsItem>
          </Descriptions>
        </template>

        <!-- Achievement Detail -->
        <template v-else-if="selectedItemType === 'achievement'">
          <Descriptions :column="1" bordered size="small">
            <DescriptionsItem label="Siswa">
              <div class="student-info">
                <Avatar :style="{ backgroundColor: '#f97316' }" size="small">
                  <template #icon><UserOutlined /></template>
                </Avatar>
                <span style="margin-left: 8px">{{ (selectedItem as Achievement).studentName }}</span>
              </div>
            </DescriptionsItem>
            <DescriptionsItem label="NIS">{{ (selectedItem as Achievement).studentNis }}</DescriptionsItem>
            <DescriptionsItem label="Prestasi">{{ (selectedItem as Achievement).title }}</DescriptionsItem>
            <DescriptionsItem label="Poin">
              <Tag color="success">+{{ (selectedItem as Achievement).point }}</Tag>
            </DescriptionsItem>
            <DescriptionsItem v-if="(selectedItem as Achievement).description" label="Deskripsi">
              {{ (selectedItem as Achievement).description }}
            </DescriptionsItem>
            <DescriptionsItem label="Dicatat Oleh">{{ (selectedItem as Achievement).createdByName }}</DescriptionsItem>
            <DescriptionsItem label="Tanggal">{{ formatFullDate((selectedItem as Achievement).createdAt) }}</DescriptionsItem>
          </Descriptions>
        </template>

        <!-- Permit Detail -->
        <template v-else-if="selectedItemType === 'permit'">
          <Descriptions :column="1" bordered size="small">
            <DescriptionsItem label="Siswa">
              <div class="student-info">
                <Avatar :style="{ backgroundColor: '#f97316' }" size="small">
                  <template #icon><UserOutlined /></template>
                </Avatar>
                <span style="margin-left: 8px">{{ (selectedItem as Permit).studentName }}</span>
              </div>
            </DescriptionsItem>
            <DescriptionsItem label="NIS">{{ (selectedItem as Permit).studentNis }}</DescriptionsItem>
            <DescriptionsItem label="NISN">{{ (selectedItem as Permit).studentNisn }}</DescriptionsItem>
            <DescriptionsItem label="Alasan">{{ (selectedItem as Permit).reason }}</DescriptionsItem>
            <DescriptionsItem label="Jam Keluar">{{ (selectedItem as Permit).exitTime }}</DescriptionsItem>
            <DescriptionsItem label="Jam Kembali">
              <span v-if="(selectedItem as Permit).returnTime">{{ (selectedItem as Permit).returnTime }}</span>
              <Tag v-else color="warning">Belum Kembali</Tag>
            </DescriptionsItem>
            <DescriptionsItem label="Guru Penanggung Jawab">{{ (selectedItem as Permit).responsibleTeacherName }}</DescriptionsItem>
            <DescriptionsItem label="Tanggal">{{ formatFullDate((selectedItem as Permit).createdAt) }}</DescriptionsItem>
          </Descriptions>
        </template>
      </div>
    </Drawer>
  </div>
</template>

<style scoped>
.bk-monitoring {
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

.page-subtitle {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 4px;
}

.readonly-tag {
  font-size: 11px;
  padding: 0 6px;
  line-height: 18px;
}

.summary-card {
  cursor: pointer;
  transition: all 0.3s;
  border: 2px solid transparent;
}

.summary-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.summary-card.active {
  border-color: #f97316;
}

.summary-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.summary-text {
  flex: 1;
}

.student-info {
  display: flex;
  align-items: center;
}

:deep(.ant-badge-count) {
  background-color: #f97316;
}
</style>
