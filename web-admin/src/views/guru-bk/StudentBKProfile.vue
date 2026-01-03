<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  Row,
  Col,
  Card,
  Tabs,
  TabPane,
  Spin,
  Alert,
  Typography,
  Tag,
  Avatar,
  Statistic,
  Timeline,
  TimelineItem,
  Button,
  Empty,
  Descriptions,
  DescriptionsItem,
  Breadcrumb,
  BreadcrumbItem,
} from 'ant-design-vue'
import {
  UserOutlined,
  WarningOutlined,
  TrophyOutlined,
  FileProtectOutlined,
  MessageOutlined,
  ArrowLeftOutlined,
  CalendarOutlined,
} from '@ant-design/icons-vue'
import { bkService, schoolService } from '@/services'
import type { StudentBKProfile, Violation, Achievement, Permit, CounselingNote } from '@/types/bk'
import type { Student } from '@/types/school'

const { Title, Text } = Typography

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const error = ref<string | null>(null)
const activeTab = ref('violations')

const studentId = computed(() => Number(route.params.id))

const student = ref<Student | null>(null)
const profile = ref<StudentBKProfile | null>(null)
const violations = ref<Violation[]>([])
const achievements = ref<Achievement[]>([])
const permits = ref<Permit[]>([])
const counselingNotes = ref<CounselingNote[]>([])
const violationPoints = ref(0)

const loadData = async () => {
  loading.value = true
  error.value = null

  try {
    const [studentRes, profileRes, violationsRes, achievementsRes, permitsRes, counselingRes, violationPointsRes] = await Promise.allSettled([
      schoolService.getStudent(studentId.value),
      bkService.getStudentBKProfile(studentId.value),
      bkService.getStudentViolations(studentId.value),
      bkService.getStudentAchievements(studentId.value),
      bkService.getStudentPermits(studentId.value),
      bkService.getStudentCounselingNotes(studentId.value),
      bkService.getStudentViolationPoints(studentId.value),
    ])

    if (studentRes.status === 'fulfilled') student.value = studentRes.value
    if (profileRes.status === 'fulfilled') profile.value = profileRes.value
    if (violationsRes.status === 'fulfilled') violations.value = violationsRes.value.data || []
    if (achievementsRes.status === 'fulfilled') achievements.value = achievementsRes.value.data || []
    if (permitsRes.status === 'fulfilled') permits.value = permitsRes.value.data || []
    if (counselingRes.status === 'fulfilled') counselingNotes.value = counselingRes.value.data || []
    if (violationPointsRes.status === 'fulfilled') violationPoints.value = violationPointsRes.value || 0

    // If student data failed, show error
    if (studentRes.status === 'rejected' && profileRes.status === 'rejected') {
      error.value = 'Gagal memuat data siswa. Silakan coba lagi.'
    }
  } catch (err) {
    console.error('Failed to load student BK profile:', err)
    error.value = 'Gagal memuat data siswa. Silakan coba lagi.'
  } finally {
    loading.value = false
  }
}

// Computed achievement points from actual achievements data
const achievementPoints = computed(() => {
  return achievements.value.reduce((sum, a) => sum + (a.point || 0), 0)
})

// Computed total net points (achievement points + violation points)
const totalNetPoints = computed(() => {
  return achievementPoints.value + violationPoints.value
})

// Format date
const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
}

// Format time
const formatTime = (dateStr: string) => {
  return new Date(dateStr).toLocaleTimeString('id-ID', {
    hour: '2-digit',
    minute: '2-digit',
  })
}

// Format datetime
const formatDateTime = (dateStr: string) => {
  return `${formatDate(dateStr)} ${formatTime(dateStr)}`
}

// Get violation level color
const getViolationLevelColor = (level: string) => {
  switch (level) {
    case 'ringan': return 'warning'
    case 'sedang': return 'orange'
    case 'berat': return 'error'
    default: return 'default'
  }
}

// Get violation level label
const getViolationLevelLabel = (level: string) => {
  switch (level) {
    case 'ringan': return 'Ringan'
    case 'sedang': return 'Sedang'
    case 'berat': return 'Berat'
    default: return level
  }
}

// Go back
const goBack = () => {
  router.back()
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="student-bk-profile">
    <!-- Breadcrumb -->
    <Breadcrumb style="margin-bottom: 16px">
      <BreadcrumbItem>
        <router-link to="/bk">Dashboard BK</router-link>
      </BreadcrumbItem>
      <BreadcrumbItem>Profil Siswa</BreadcrumbItem>
    </Breadcrumb>

    <Button type="link" @click="goBack" style="padding: 0; margin-bottom: 16px">
      <ArrowLeftOutlined /> Kembali
    </Button>

    <Spin :spinning="loading">
      <Alert
        v-if="error"
        type="error"
        :message="error"
        show-icon
        closable
        style="margin-bottom: 24px"
      />

      <!-- Student Info Card -->
      <Card class="student-info-card" v-if="student">
        <Row :gutter="24" align="middle">
          <Col :xs="24" :sm="6" :md="4">
            <Avatar :size="100" :style="{ backgroundColor: '#f97316' }">
              <template #icon><UserOutlined /></template>
            </Avatar>
          </Col>
          <Col :xs="24" :sm="18" :md="12">
            <Title :level="3" style="margin: 0">{{ student.name }}</Title>
            <div class="student-meta">
              <Tag color="blue">{{ student.className }}</Tag>
              <Text type="secondary">NIS: {{ student.nis }}</Text>
              <Text type="secondary">NISN: {{ student.nisn }}</Text>
            </div>
          </Col>
          <Col :xs="24" :md="8">
            <Row :gutter="16">
              <Col :span="8">
                <Statistic
                  title="Poin Prestasi"
                  :value="achievementPoints"
                  :value-style="{ color: '#22c55e' }"
                >
                  <template #prefix><TrophyOutlined /></template>
                </Statistic>
              </Col>
              <Col :span="8">
                <Statistic
                  title="Poin Pelanggaran"
                  :value="violationPoints"
                  :value-style="{ color: '#ef4444' }"
                >
                  <template #prefix><WarningOutlined /></template>
                </Statistic>
              </Col>
              <Col :span="8">
                <Statistic
                  title="Total Poin"
                  :value="totalNetPoints"
                  :value-style="{ color: totalNetPoints >= 0 ? '#22c55e' : '#ef4444' }"
                />
              </Col>
            </Row>
          </Col>
        </Row>
      </Card>

      <!-- Tabs for BK Data -->
      <Card style="margin-top: 24px">
        <Tabs v-model:activeKey="activeTab">
          <!-- Pelanggaran Tab -->
          <TabPane key="violations">
            <template #tab>
              <span><WarningOutlined /> Pelanggaran ({{ violations.length }})</span>
            </template>
            <Timeline v-if="violations.length > 0">
              <TimelineItem
                v-for="violation in violations"
                :key="violation.id"
                :color="violation.level === 'berat' ? 'red' : violation.level === 'sedang' ? 'orange' : 'gold'"
              >
                <div class="timeline-content">
                  <div class="timeline-header">
                    <Text strong>{{ violation.category }}</Text>
                    <Tag :color="getViolationLevelColor(violation.level)">
                      {{ getViolationLevelLabel(violation.level) }}
                    </Tag>
                    <Tag color="error">{{ violation.point }} poin</Tag>
                  </div>
                  <Text>{{ violation.description }}</Text>
                  <div class="timeline-meta">
                    <Text type="secondary">
                      <CalendarOutlined /> {{ formatDate(violation.createdAt) }}
                    </Text>
                    <Text type="secondary" v-if="violation.createdByName">
                      oleh {{ violation.createdByName }}
                    </Text>
                  </div>
                </div>
              </TimelineItem>
            </Timeline>
            <Empty v-else description="Tidak ada catatan pelanggaran" />
          </TabPane>

          <!-- Prestasi Tab -->
          <TabPane key="achievements">
            <template #tab>
              <span><TrophyOutlined /> Prestasi ({{ achievements.length }})</span>
            </template>
            <Timeline v-if="achievements.length > 0">
              <TimelineItem
                v-for="achievement in achievements"
                :key="achievement.id"
                color="green"
              >
                <div class="timeline-content">
                  <div class="timeline-header">
                    <Text strong>{{ achievement.title }}</Text>
                    <Tag color="success">+{{ achievement.point }} poin</Tag>
                  </div>
                  <Text v-if="achievement.description">{{ achievement.description }}</Text>
                  <div class="timeline-meta">
                    <Text type="secondary">
                      <CalendarOutlined /> {{ formatDate(achievement.createdAt) }}
                    </Text>
                    <Text type="secondary" v-if="achievement.createdByName">
                      oleh {{ achievement.createdByName }}
                    </Text>
                  </div>
                </div>
              </TimelineItem>
            </Timeline>
            <Empty v-else description="Tidak ada catatan prestasi" />
          </TabPane>

          <!-- Izin Keluar Tab -->
          <TabPane key="permits">
            <template #tab>
              <span><FileProtectOutlined /> Izin Keluar ({{ permits.length }})</span>
            </template>
            <Timeline v-if="permits.length > 0">
              <TimelineItem
                v-for="permit in permits"
                :key="permit.id"
                color="blue"
              >
                <div class="timeline-content">
                  <div class="timeline-header">
                    <Text strong>{{ permit.reason }}</Text>
                    <Tag :color="permit.returnTime ? 'success' : 'processing'">
                      {{ permit.returnTime ? 'Sudah Kembali' : 'Belum Kembali' }}
                    </Tag>
                  </div>
                  <Descriptions :column="1" size="small" style="margin-top: 8px">
                    <DescriptionsItem label="Waktu Keluar">
                      {{ formatDateTime(permit.exitTime) }}
                    </DescriptionsItem>
                    <DescriptionsItem label="Waktu Kembali" v-if="permit.returnTime">
                      {{ formatDateTime(permit.returnTime) }}
                    </DescriptionsItem>
                    <DescriptionsItem label="Guru Penanggung Jawab">
                      {{ permit.responsibleTeacherName }}
                    </DescriptionsItem>
                  </Descriptions>
                  <div class="timeline-meta">
                    <Text type="secondary">
                      <CalendarOutlined /> {{ formatDate(permit.createdAt) }}
                    </Text>
                  </div>
                </div>
              </TimelineItem>
            </Timeline>
            <Empty v-else description="Tidak ada catatan izin keluar" />
          </TabPane>

          <!-- Konseling Tab -->
          <TabPane key="counseling">
            <template #tab>
              <span><MessageOutlined /> Konseling ({{ counselingNotes.length }})</span>
            </template>
            <Timeline v-if="counselingNotes.length > 0">
              <TimelineItem
                v-for="note in counselingNotes"
                :key="note.id"
                color="purple"
              >
                <div class="timeline-content">
                  <Card size="small" class="counseling-card">
                    <div class="counseling-section">
                      <Text strong type="secondary">Catatan Internal (Rahasia)</Text>
                      <div class="internal-note">
                        {{ note.internalNote }}
                      </div>
                    </div>
                    <div class="counseling-section" v-if="note.parentSummary">
                      <Text strong type="secondary">Ringkasan untuk Orang Tua</Text>
                      <div class="parent-summary">
                        {{ note.parentSummary }}
                      </div>
                    </div>
                  </Card>
                  <div class="timeline-meta">
                    <Text type="secondary">
                      <CalendarOutlined /> {{ formatDate(note.createdAt) }}
                    </Text>
                    <Text type="secondary" v-if="note.createdByName">
                      oleh {{ note.createdByName }}
                    </Text>
                  </div>
                </div>
              </TimelineItem>
            </Timeline>
            <Empty v-else description="Tidak ada catatan konseling" />
          </TabPane>
        </Tabs>
      </Card>
    </Spin>
  </div>
</template>

<style scoped>
.student-bk-profile {
  padding: 0;
}

.student-info-card {
  margin-bottom: 24px;
}

.student-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 8px;
  flex-wrap: wrap;
}

.timeline-content {
  padding-bottom: 8px;
}

.timeline-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
  flex-wrap: wrap;
}

.timeline-meta {
  display: flex;
  gap: 16px;
  margin-top: 8px;
  font-size: 12px;
}

.counseling-card {
  margin-bottom: 8px;
}

.counseling-section {
  margin-bottom: 12px;
}

.counseling-section:last-child {
  margin-bottom: 0;
}

.internal-note {
  background-color: #fff7e6;
  border: 1px solid #ffd591;
  border-radius: 4px;
  padding: 8px 12px;
  margin-top: 4px;
}

.parent-summary {
  background-color: #f6ffed;
  border: 1px solid #b7eb8f;
  border-radius: 4px;
  padding: 8px 12px;
  margin-top: 4px;
}

@media (max-width: 768px) {
  .student-info-card :deep(.ant-avatar) {
    margin-bottom: 16px;
  }
}
</style>
