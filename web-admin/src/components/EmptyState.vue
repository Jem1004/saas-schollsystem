<script setup lang="ts">
import { computed } from 'vue'
import { Button, Typography } from 'ant-design-vue'
import {
  PlusOutlined,
  FileSearchOutlined,
  InboxOutlined,
  TeamOutlined,
  BookOutlined,
  TrophyOutlined,
  FileTextOutlined,
  ScheduleOutlined,
} from '@ant-design/icons-vue'

const { Text, Paragraph } = Typography

interface Props {
  type?: 'default' | 'search' | 'no-data' | 'students' | 'classes' | 'achievements' | 'violations' | 'permits' | 'notes' | 'attendance'
  title?: string
  description?: string
  actionText?: string
  showAction?: boolean
  imageSize?: number
}

const props = withDefaults(defineProps<Props>(), {
  type: 'default',
  title: '',
  description: '',
  actionText: '',
  showAction: true,
  imageSize: 100,
})

const emit = defineEmits<{
  (e: 'action'): void
}>()

const typeConfig = {
  default: {
    icon: InboxOutlined,
    title: 'Tidak Ada Data',
    description: 'Belum ada data yang tersedia.',
    actionText: 'Tambah Data',
    color: '#8c8c8c',
  },
  search: {
    icon: FileSearchOutlined,
    title: 'Tidak Ditemukan',
    description: 'Tidak ada hasil yang cocok dengan pencarian Anda.',
    actionText: 'Reset Pencarian',
    color: '#1890ff',
  },
  'no-data': {
    icon: InboxOutlined,
    title: 'Belum Ada Data',
    description: 'Data akan muncul di sini setelah ditambahkan.',
    actionText: 'Tambah Sekarang',
    color: '#8c8c8c',
  },
  students: {
    icon: TeamOutlined,
    title: 'Belum Ada Siswa',
    description: 'Belum ada siswa yang terdaftar. Tambahkan siswa untuk memulai.',
    actionText: 'Tambah Siswa',
    color: '#f97316',
  },
  classes: {
    icon: BookOutlined,
    title: 'Belum Ada Kelas',
    description: 'Belum ada kelas yang dibuat. Buat kelas untuk mengorganisir siswa.',
    actionText: 'Buat Kelas',
    color: '#1890ff',
  },
  achievements: {
    icon: TrophyOutlined,
    title: 'Belum Ada Prestasi',
    description: 'Belum ada prestasi yang dicatat. Catat prestasi siswa untuk memberikan apresiasi.',
    actionText: 'Catat Prestasi',
    color: '#52c41a',
  },
  violations: {
    icon: FileTextOutlined,
    title: 'Tidak Ada Pelanggaran',
    description: 'Tidak ada pelanggaran yang tercatat. Ini adalah hal yang baik!',
    actionText: 'Catat Pelanggaran',
    color: '#ff4d4f',
  },
  permits: {
    icon: FileTextOutlined,
    title: 'Belum Ada Izin Keluar',
    description: 'Belum ada izin keluar yang dicatat hari ini.',
    actionText: 'Buat Izin',
    color: '#722ed1',
  },
  notes: {
    icon: FileTextOutlined,
    title: 'Belum Ada Catatan',
    description: 'Belum ada catatan yang dibuat. Buat catatan untuk mendokumentasikan perkembangan siswa.',
    actionText: 'Buat Catatan',
    color: '#13c2c2',
  },
  attendance: {
    icon: ScheduleOutlined,
    title: 'Belum Ada Data Absensi',
    description: 'Data absensi akan muncul setelah siswa melakukan check-in.',
    actionText: 'Input Manual',
    color: '#fa8c16',
  },
}

const config = computed(() => {
  const baseConfig = typeConfig[props.type] || typeConfig.default
  return {
    ...baseConfig,
    title: props.title || baseConfig.title,
    description: props.description || baseConfig.description,
    actionText: props.actionText || baseConfig.actionText,
  }
})

const handleAction = () => {
  emit('action')
}
</script>

<template>
  <div class="empty-state">
    <div class="empty-icon" :style="{ backgroundColor: `${config.color}10` }">
      <component :is="config.icon" :style="{ fontSize: `${imageSize * 0.5}px`, color: config.color }" />
    </div>
    <Text strong class="empty-title">{{ config.title }}</Text>
    <Paragraph type="secondary" class="empty-description">
      {{ config.description }}
    </Paragraph>
    <Button 
      v-if="showAction" 
      type="primary" 
      @click="handleAction"
    >
      <template #icon><PlusOutlined /></template>
      {{ config.actionText }}
    </Button>
  </div>
</template>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
  text-align: center;
}

.empty-icon {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
}

.empty-title {
  font-size: 18px;
  margin-bottom: 8px;
  display: block;
}

.empty-description {
  max-width: 300px;
  margin-bottom: 24px;
}
</style>
