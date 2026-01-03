<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import {
  Table, Button, Input, Space, Tag, Card, Row, Col, Typography, Avatar, Statistic,
} from 'ant-design-vue'
import type { TableProps } from 'ant-design-vue'
import {
  SearchOutlined, ReloadOutlined, UserOutlined, WarningOutlined, TrophyOutlined, EyeOutlined,
} from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'
import { bkService, schoolService } from '@/services'
import type { StudentBKProfile } from '@/types/bk'
import type { Student } from '@/types/school'
import type { Class } from '@/types/school'

const { Title } = Typography
const router = useRouter()

const loading = ref(false)
const students = ref<Student[]>([])
const studentProfiles = ref<Map<number, StudentBKProfile>>(new Map())
const studentViolationPoints = ref<Map<number, number>>(new Map())
const classes = ref<Class[]>([])
const total = ref(0)
const pagination = reactive({ current: 1, pageSize: 10 })
const searchText = ref('')
const filterClassId = ref<number | undefined>(undefined)

const columns: TableProps['columns'] = [
  { title: 'Siswa', key: 'student', width: 250 },
  { title: 'Kelas', dataIndex: 'className', key: 'className', width: 100 },
  { title: 'NIS', dataIndex: 'nis', key: 'nis', width: 120 },
  { title: 'Pelanggaran', key: 'violations', width: 120, align: 'center' },
  { title: 'Prestasi', key: 'achievements', width: 120, align: 'center' },
  { title: 'Poin', key: 'points', width: 100, align: 'center' },
  { title: 'Aksi', key: 'action', width: 100, align: 'center' },
]

const filteredStudents = computed(() => {
  let result = students.value
  if (filterClassId.value) {
    result = result.filter(s => s.classId === filterClassId.value)
  }
  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(s => 
      s.name.toLowerCase().includes(search) || 
      s.nis.toLowerCase().includes(search)
    )
  }
  return result
})

const totalStats = computed(() => {
  let violations = 0, achievements = 0, achievementPoints = 0, violationPts = 0
  studentProfiles.value.forEach((p, studentId) => {
    violations += p.violationCount || 0
    achievements += p.achievementCount || 0
    achievementPoints += p.totalPoints || p.totalAchievementPoints || 0
    violationPts += studentViolationPoints.value.get(studentId) || 0
  })
  return { violations, achievements, achievementPoints, violationPoints: violationPts, totalPoints: achievementPoints + violationPts }
})

const getStudentProfile = (studentId: number): StudentBKProfile | undefined => {
  return studentProfiles.value.get(studentId)
}

const getStudentTotalPoints = (studentId: number): number => {
  const profile = studentProfiles.value.get(studentId)
  const achievementPts = profile?.totalPoints || profile?.totalAchievementPoints || 0
  const violationPts = studentViolationPoints.value.get(studentId) || 0
  return achievementPts + violationPts
}

const loadData = async () => {
  loading.value = true
  try {
    const [studentsRes, classesRes] = await Promise.all([
      schoolService.getStudents({ page: pagination.current, page_size: pagination.pageSize }),
      schoolService.getClasses({ page_size: 100 }),
    ])
    students.value = studentsRes.students || []
    total.value = studentsRes.pagination?.total || 0
    classes.value = classesRes.classes || []

    // Load BK profiles and violation points for each student
    const dataPromises = students.value.map(async (s) => {
      try {
        const [profile, violationPts] = await Promise.all([
          bkService.getStudentBKProfile(s.id),
          bkService.getStudentViolationPoints(s.id),
        ])
        return { id: s.id, profile, violationPts }
      } catch {
        return { id: s.id, profile: null, violationPts: 0 }
      }
    })
    const results = await Promise.all(dataPromises)
    const profileMap = new Map<number, StudentBKProfile>()
    const violationPtsMap = new Map<number, number>()
    results.forEach(r => {
      if (r.profile) profileMap.set(r.id, r.profile)
      violationPtsMap.set(r.id, r.violationPts)
    })
    studentProfiles.value = profileMap
    studentViolationPoints.value = violationPtsMap
  } catch (err) {
    console.error('Failed to load students:', err)
  } finally {
    loading.value = false
  }
}

const handleTableChange: TableProps['onChange'] = (pag) => {
  pagination.current = pag.current || 1
  pagination.pageSize = pag.pageSize || 10
  loadData()
}

const handleSearch = () => {
  pagination.current = 1
  loadData()
}

const viewStudentProfile = (studentId: number) => {
  router.push(`/bk/students/${studentId}`)
}

onMounted(() => { loadData() })
</script>

<template>
  <div class="bk-student-list">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Profil Siswa BK</Title>
    </div>

    <Row :gutter="24" style="margin-bottom: 24px">
      <Col :xs="24" :sm="8">
        <Card>
          <Statistic title="Total Pelanggaran" :value="totalStats.violations" :value-style="{ color: '#ef4444' }">
            <template #prefix><WarningOutlined /></template>
          </Statistic>
        </Card>
      </Col>
      <Col :xs="24" :sm="8">
        <Card>
          <Statistic title="Total Prestasi" :value="totalStats.achievements" :value-style="{ color: '#22c55e' }">
            <template #prefix><TrophyOutlined /></template>
          </Statistic>
        </Card>
      </Col>
      <Col :xs="24" :sm="8">
        <Card>
          <Statistic title="Total Poin" :value="totalStats.totalPoints" :value-style="{ color: totalStats.totalPoints >= 0 ? '#22c55e' : '#ef4444' }" suffix="poin" />
        </Card>
      </Col>
    </Row>

    <Card>
      <Row :gutter="16" class="toolbar" justify="space-between" align="middle">
        <Col :xs="24" :md="16">
          <Space wrap>
            <Input v-model:value="searchText" placeholder="Cari nama atau NIS..." allow-clear style="width: 250px" @press-enter="handleSearch">
              <template #prefix><SearchOutlined /></template>
            </Input>
            <a-select v-model:value="filterClassId" placeholder="Filter Kelas" allow-clear style="width: 150px" @change="handleSearch">
              <a-select-option v-for="c in classes" :key="c.id" :value="c.id">{{ c.name }}</a-select-option>
            </a-select>
          </Space>
        </Col>
        <Col :xs="24" :md="8" class="toolbar-right">
          <Button @click="loadData"><template #icon><ReloadOutlined /></template></Button>
        </Col>
      </Row>

      <Table :columns="columns" :data-source="filteredStudents" :loading="loading"
        :pagination="{ current: pagination.current, pageSize: pagination.pageSize, total, showSizeChanger: true, showTotal: (t: number) => `Total ${t} siswa` }"
        row-key="id" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'student'">
            <div class="student-cell">
              <Avatar :style="{ backgroundColor: '#f97316' }"><template #icon><UserOutlined /></template></Avatar>
              <div class="student-info">
                <a @click="viewStudentProfile((record as Student).id)">{{ (record as Student).name }}</a>
                <span class="student-nisn">NISN: {{ (record as Student).nisn }}</span>
              </div>
            </div>
          </template>
          <template v-else-if="column.key === 'className'">
            <Tag color="blue">{{ (record as Student).className }}</Tag>
          </template>
          <template v-else-if="column.key === 'violations'">
            <Tag color="error">{{ getStudentProfile((record as Student).id)?.violationCount || 0 }}</Tag>
          </template>
          <template v-else-if="column.key === 'achievements'">
            <Tag color="success">{{ getStudentProfile((record as Student).id)?.achievementCount || 0 }}</Tag>
          </template>
          <template v-else-if="column.key === 'points'">
            <span :class="['points', getStudentTotalPoints((record as Student).id) >= 0 ? 'positive' : 'negative']">
              {{ getStudentTotalPoints((record as Student).id) }}
            </span>
          </template>
          <template v-else-if="column.key === 'action'">
            <Button type="primary" size="small" @click="viewStudentProfile((record as Student).id)">
              <template #icon><EyeOutlined /></template> Detail
            </Button>
          </template>
        </template>
      </Table>
    </Card>
  </div>
</template>

<style scoped>
.bk-student-list { padding: 0; }
.page-header { margin-bottom: 24px; }
.toolbar { margin-bottom: 16px; }
.toolbar-right { display: flex; justify-content: flex-end; }
.student-cell { display: flex; align-items: center; gap: 12px; }
.student-info { display: flex; flex-direction: column; }
.student-info a { font-weight: 500; }
.student-nisn { font-size: 12px; color: #8c8c8c; }
.points { font-weight: 600; }
.points.positive { color: #22c55e; }
.points.negative { color: #ef4444; }
@media (max-width: 768px) { .toolbar-right { margin-top: 16px; justify-content: flex-start; } }
</style>
