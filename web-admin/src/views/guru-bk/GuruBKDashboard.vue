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
  Avatar,
  Button,
  Empty,
} from 'ant-design-vue'
import {
  WarningOutlined,
  TrophyOutlined,
  FileProtectOutlined,
  MessageOutlined,
  UserOutlined,
  RightOutlined,
  CalendarOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'
import { bkService } from '@/services'
import type { BKStats, Violation, Achievement, StudentAttentionItem } from '@/types/bk'

const { Title, Text } = Typography
const router = useRouter()

const loading = ref(true)
const error = ref<string | null>(null)

const stats = ref<BKStats>({
  totalViolations: 0,
  totalAchievements: 0,
  totalPermits: 0,
  totalCounselingNotes: 0,
  recentViolations: [],
  recentAchievements: [],
  studentsNeedingAttention: [],
})

const loadData = async () => {
  loading.value = true
  error.value = null

  try {
    const response = await bkService.getStats()
    stats.value = response
  } catch (err) {
    console.error('Failed to load BK stats:', err)
    error.value = 'Gagal memuat data statistik BK. Silakan coba lagi.'
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

// Get violation level color class
const getViolationLevelColor = (level: string) => {
  switch (level) {
    case 'ringan': return 'warning'
    case 'sedang': return 'orange'
    case 'berat': return 'error'
    default: return 'default'
  }
}

// Navigate to student profile
const goToStudentProfile = (studentId: number) => {
  router.push(`/bk/students/${studentId}`)
}

// Navigate to violations page
const goToViolations = () => {
  router.push('/bk/violations')
}

// Navigate to achievements page
const goToAchievements = () => {
  router.push('/bk/achievements')
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="guru-bk-dashboard">
    <div class="page-header">
      <Title :level="2" style="margin: 0">Dashboard Guru BK</Title>
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
          <Card class="stat-card" hoverable @click="goToViolations">
            <Statistic
              title="Total Pelanggaran"
              :value="stats.totalViolations"
              :value-style="{ color: '#ef4444' }"
            >
              <template #prefix>
                <WarningOutlined />
              </template>
            </Statistic>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :lg="6">
          <Card class="stat-card" hoverable @click="goToAchievements">
            <Statistic
              title="Total Prestasi"
              :value="stats.totalAchievements"
              :value-style="{ color: '#22c55e' }"
            >
              <template #prefix>
                <TrophyOutlined />
              </template>
            </Statistic>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :lg="6">
          <Card class="stat-card" hoverable @click="() => router.push('/bk/permits')">
            <Statistic
              title="Izin Keluar"
              :value="stats.totalPermits"
              :value-style="{ color: '#3b82f6' }"
            >
              <template #prefix>
                <FileProtectOutlined />
              </template>
            </Statistic>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :lg="6">
          <Card class="stat-card" hoverable @click="() => router.push('/bk/counseling')">
            <Statistic
              title="Catatan Konseling"
              :value="stats.totalCounselingNotes"
              :value-style="{ color: '#8b5cf6' }"
            >
              <template #prefix>
                <MessageOutlined />
              </template>
            </Statistic>
          </Card>
        </Col>
      </Row>

      <!-- Students Requiring Attention -->
      <Row :gutter="[24, 24]" style="margin-top: 24px">
        <Col :xs="24">
          <Card title="Siswa Perlu Perhatian" class="attention-card">
            <template #extra>
              <ExclamationCircleOutlined style="color: #f97316; font-size: 18px" />
            </template>
            <List
              v-if="stats.studentsNeedingAttention && stats.studentsNeedingAttention.length > 0"
              :data-source="stats.studentsNeedingAttention"
              :loading="loading"
            >
              <template #renderItem="{ item }">
                <ListItem class="attention-item" @click="goToStudentProfile((item as StudentAttentionItem).studentId)">
                  <ListItemMeta>
                    <template #avatar>
                      <Avatar :style="{ backgroundColor: '#f97316' }">
                        <template #icon><UserOutlined /></template>
                      </Avatar>
                    </template>
                    <template #title>
                      <span class="student-name">{{ (item as StudentAttentionItem).studentName }}</span>
                      <span class="class-badge" style="margin-left: 8px">{{ (item as StudentAttentionItem).className }}</span>
                    </template>
                    <template #description>
                      <div class="student-stats">
                        <span class="status-badge error">{{ (item as StudentAttentionItem).violationCount }} pelanggaran</span>
                        <span style="margin-left: 8px; color: #8c8c8c">{{ (item as StudentAttentionItem).reason }}</span>
                      </div>
                    </template>
                  </ListItemMeta>
                  <template #actions>
                    <Button type="link" size="small">
                      Lihat Detail <RightOutlined />
                    </Button>
                  </template>
                </ListItem>
              </template>
            </List>
            <Empty v-else description="Tidak ada siswa yang memerlukan perhatian khusus" />
          </Card>
        </Col>
      </Row>

      <!-- Recent Activities -->
      <Row :gutter="[24, 24]" style="margin-top: 24px">
        <Col :xs="24" :lg="12">
          <Card title="Pelanggaran Terbaru" class="recent-card">
            <template #extra>
              <Button type="link" size="small" @click="goToViolations">
                Lihat Semua <RightOutlined />
              </Button>
            </template>
            <List
              v-if="stats.recentViolations.length > 0"
              :data-source="stats.recentViolations"
              :loading="loading"
              size="small"
            >
              <template #renderItem="{ item }">
                <ListItem class="recent-item" @click="goToStudentProfile((item as Violation).studentId)">
                  <ListItemMeta>
                    <template #avatar>
                      <Avatar :style="{ backgroundColor: '#ef4444' }" size="small">
                        <template #icon><WarningOutlined /></template>
                      </Avatar>
                    </template>
                    <template #title>
                      <span>{{ (item as Violation).studentName }}</span>
                      <span class="status-badge" :class="getViolationLevelColor((item as Violation).level)" style="margin-left: 8px">
                        {{ (item as Violation).level }}
                      </span>
                    </template>
                    <template #description>
                      <div>{{ (item as Violation).category }} - {{ (item as Violation).description }}</div>
                      <Text type="secondary" class="date-text">{{ formatDate((item as Violation).createdAt) }}</Text>
                    </template>
                  </ListItemMeta>
                </ListItem>
              </template>
            </List>
            <Empty v-else description="Belum ada pelanggaran tercatat" />
          </Card>
        </Col>
        <Col :xs="24" :lg="12">
          <Card title="Prestasi Terbaru" class="recent-card">
            <template #extra>
              <Button type="link" size="small" @click="goToAchievements">
                Lihat Semua <RightOutlined />
              </Button>
            </template>
            <List
              v-if="stats.recentAchievements.length > 0"
              :data-source="stats.recentAchievements"
              :loading="loading"
              size="small"
            >
              <template #renderItem="{ item }">
                <ListItem class="recent-item" @click="goToStudentProfile((item as Achievement).studentId)">
                  <ListItemMeta>
                    <template #avatar>
                      <Avatar :style="{ backgroundColor: '#22c55e' }" size="small">
                        <template #icon><TrophyOutlined /></template>
                      </Avatar>
                    </template>
                    <template #title>
                      <span>{{ (item as Achievement).studentName }}</span>
                      <span class="status-badge success" style="margin-left: 8px">+{{ (item as Achievement).point }} poin</span>
                    </template>
                    <template #description>
                      <div>{{ (item as Achievement).title }}</div>
                      <Text type="secondary" class="date-text">{{ formatDate((item as Achievement).createdAt) }}</Text>
                    </template>
                  </ListItemMeta>
                </ListItem>
              </template>
            </List>
            <Empty v-else description="Belum ada prestasi tercatat" />
          </Card>
        </Col>
      </Row>
    </Spin>
  </div>
</template>

<style scoped>
.guru-bk-dashboard {
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

.attention-card,
.recent-card {
  height: 100%;
}

.attention-item,
.recent-item {
  cursor: pointer;
  transition: background-color 0.2s;
}

.attention-item:hover,
.recent-item:hover {
  background-color: #fafafa;
}

.student-name {
  font-weight: 500;
}

.student-stats {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
  margin-top: 4px;
}

.date-text {
  font-size: 12px;
  display: block;
  margin-top: 4px;
}

/* Badge Styles */
.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
}

.status-badge.success {
  background-color: #f6ffed;
  color: #52c41a;
  border: 1px solid #b7eb8f;
}

.status-badge.warning {
  background-color: #fffbe6;
  color: #faad14;
  border: 1px solid #ffe58f;
}

.status-badge.error {
  background-color: #fff2f0;
  color: #ff4d4f;
  border: 1px solid #ffccc7;
}

.status-badge.orange {
  background-color: #fff7e6;
  color: #fa8c16;
  border: 1px solid #ffd591;
}

.status-badge.default {
  background-color: #f5f5f5;
  color: #000000d9;
  border: 1px solid #d9d9d9;
}

.class-badge {
  background-color: #e6f7ff;
  color: #1890ff;
  border: 1px solid #91d5ff;
  padding: 0 8px;
  border-radius: 4px;
  font-size: 12px;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
