<script setup lang="ts">
import { ref, computed } from 'vue'
import { Modal, Button, Space, Spin, Descriptions, DescriptionsItem, Divider, Typography } from 'ant-design-vue'
import { PrinterOutlined, DownloadOutlined } from '@ant-design/icons-vue'
import type { Permit } from '@/types/bk'

const { Title, Text } = Typography

interface Props {
  open: boolean
  permit: Permit | null
  schoolName?: string
  schoolAddress?: string
  schoolPhone?: string
}

const props = withDefaults(defineProps<Props>(), {
  open: false,
  schoolName: 'SMP Negeri 1 Contoh',
  schoolAddress: 'Jl. Pendidikan No. 1, Kota Contoh',
  schoolPhone: '(021) 1234567',
})

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'close'): void
}>()

const loading = ref(false)
const printRef = ref<HTMLDivElement | null>(null)

// Format date to Indonesian locale
const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  return date.toLocaleDateString('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
}

// Format time
const formatTime = (dateString: string): string => {
  const date = new Date(dateString)
  return date.toLocaleTimeString('id-ID', {
    hour: '2-digit',
    minute: '2-digit',
  })
}

// Current date for document
const currentDate = computed(() => {
  return new Date().toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
})

// Handle close
const handleClose = () => {
  emit('update:open', false)
  emit('close')
}

// Handle print
const handlePrint = () => {
  const printContent = printRef.value
  if (!printContent) return
  
  const printWindow = window.open('', '_blank')
  if (!printWindow) {
    alert('Popup diblokir. Silakan izinkan popup untuk mencetak.')
    return
  }
  
  printWindow.document.write(`
    <!DOCTYPE html>
    <html>
    <head>
      <title>Surat Izin Keluar - ${props.permit?.studentName || 'Siswa'}</title>
      <style>
        * {
          margin: 0;
          padding: 0;
          box-sizing: border-box;
        }
        body {
          font-family: 'Times New Roman', Times, serif;
          font-size: 12pt;
          line-height: 1.5;
          padding: 20mm;
          color: #000;
        }
        .document {
          max-width: 210mm;
          margin: 0 auto;
        }
        .header {
          text-align: center;
          border-bottom: 3px double #000;
          padding-bottom: 10px;
          margin-bottom: 20px;
        }
        .header h1 {
          font-size: 16pt;
          margin-bottom: 5px;
        }
        .header p {
          font-size: 10pt;
          margin: 2px 0;
        }
        .title {
          text-align: center;
          margin: 20px 0;
        }
        .title h2 {
          font-size: 14pt;
          text-decoration: underline;
          margin-bottom: 5px;
        }
        .title p {
          font-size: 10pt;
        }
        .content {
          margin: 20px 0;
        }
        .content p {
          margin: 10px 0;
          text-align: justify;
        }
        .student-info {
          margin: 15px 0 15px 30px;
        }
        .student-info table {
          border-collapse: collapse;
        }
        .student-info td {
          padding: 3px 10px 3px 0;
          vertical-align: top;
        }
        .student-info td:first-child {
          width: 120px;
        }
        .signature {
          margin-top: 40px;
          display: flex;
          justify-content: space-between;
        }
        .signature-box {
          text-align: center;
          width: 200px;
        }
        .signature-box .name {
          margin-top: 60px;
          border-bottom: 1px solid #000;
          padding-bottom: 2px;
        }
        .footer {
          margin-top: 30px;
          font-size: 9pt;
          color: #666;
          text-align: center;
        }
        @media print {
          body {
            padding: 0;
          }
          .no-print {
            display: none;
          }
        }
      </style>
    </head>
    <body>
      ${printContent.innerHTML}
    </body>
    </html>
  `)
  
  printWindow.document.close()
  printWindow.focus()
  
  setTimeout(() => {
    printWindow.print()
    printWindow.close()
  }, 250)
}

// Handle download as PDF (using print to PDF)
const handleDownload = () => {
  handlePrint()
}
</script>

<template>
  <Modal
    :open="open"
    title="Preview Surat Izin Keluar"
    width="800px"
    :footer="null"
    @cancel="handleClose"
  >
    <Spin v-if="loading" />
    
    <template v-else-if="permit">
      <!-- Action buttons -->
      <div class="action-bar">
        <Space>
          <Button type="primary" @click="handlePrint">
            <template #icon><PrinterOutlined /></template>
            Cetak
          </Button>
          <Button @click="handleDownload">
            <template #icon><DownloadOutlined /></template>
            Download PDF
          </Button>
        </Space>
      </div>
      
      <!-- Document preview -->
      <div ref="printRef" class="document-preview">
        <div class="document">
          <!-- Header -->
          <div class="document-header">
            <Title :level="4" style="margin: 0">{{ schoolName }}</Title>
            <Text type="secondary">{{ schoolAddress }}</Text>
            <br />
            <Text type="secondary">Telp: {{ schoolPhone }}</Text>
          </div>
          
          <Divider style="border-color: #000; border-width: 2px" />
          
          <!-- Title -->
          <div class="document-title">
            <Title :level="4" style="margin: 0; text-decoration: underline">
              SURAT IZIN KELUAR
            </Title>
            <Text type="secondary">No: {{ permit.id }}/SIK/{{ new Date(permit.createdAt).getFullYear() }}</Text>
          </div>
          
          <!-- Content -->
          <div class="document-content">
            <p>Yang bertanda tangan di bawah ini, Guru BK {{ schoolName }}, menerangkan bahwa:</p>
            
            <Descriptions :column="1" bordered size="small" class="student-info">
              <DescriptionsItem label="Nama">{{ permit.studentName }}</DescriptionsItem>
              <DescriptionsItem label="NIS">{{ permit.studentNis }}</DescriptionsItem>
              <DescriptionsItem label="NISN">{{ permit.studentNisn }}</DescriptionsItem>
              <DescriptionsItem label="Kelas">{{ permit.studentClass }}</DescriptionsItem>
            </Descriptions>
            
            <p>
              Diberikan izin untuk meninggalkan sekolah pada:
            </p>
            
            <Descriptions :column="1" bordered size="small" class="permit-info">
              <DescriptionsItem label="Hari/Tanggal">{{ formatDate(permit.exitTime) }}</DescriptionsItem>
              <DescriptionsItem label="Jam Keluar">{{ formatTime(permit.exitTime) }}</DescriptionsItem>
              <DescriptionsItem v-if="permit.returnTime" label="Jam Kembali">
                {{ formatTime(permit.returnTime) }}
              </DescriptionsItem>
              <DescriptionsItem label="Alasan">{{ permit.reason }}</DescriptionsItem>
              <DescriptionsItem label="Guru Penanggung Jawab">
                {{ permit.responsibleTeacherName }}
              </DescriptionsItem>
            </Descriptions>
            
            <p>
              Demikian surat izin ini dibuat untuk dapat dipergunakan sebagaimana mestinya.
            </p>
          </div>
          
          <!-- Signature -->
          <div class="document-signature">
            <div class="signature-left">
              <Text>{{ schoolAddress?.split(',')[0] }}, {{ currentDate }}</Text>
              <br />
              <Text>Guru BK,</Text>
              <div class="signature-space"></div>
              <Text strong style="text-decoration: underline">{{ permit.createdByName }}</Text>
            </div>
            
            <div class="signature-right">
              <Text>Mengetahui,</Text>
              <br />
              <Text>Orang Tua/Wali</Text>
              <div class="signature-space"></div>
              <Text>(...............................)</Text>
            </div>
          </div>
          
          <!-- Footer -->
          <div class="document-footer">
            <Divider dashed />
            <Text type="secondary" style="font-size: 10px">
              Dokumen ini dicetak secara otomatis oleh Sistem Manajemen Sekolah pada {{ new Date().toLocaleString('id-ID') }}
            </Text>
          </div>
        </div>
      </div>
    </template>
    
    <template v-else>
      <div class="empty-state">
        <Text type="secondary">Tidak ada data izin untuk ditampilkan</Text>
      </div>
    </template>
  </Modal>
</template>

<style scoped>
.action-bar {
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.document-preview {
  background: #fff;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  padding: 24px;
  max-height: 600px;
  overflow-y: auto;
}

.document {
  font-family: 'Times New Roman', Times, serif;
  line-height: 1.6;
}

.document-header {
  text-align: center;
  margin-bottom: 8px;
}

.document-title {
  text-align: center;
  margin: 24px 0;
}

.document-content {
  margin: 24px 0;
}

.document-content p {
  margin: 16px 0;
  text-align: justify;
}

.student-info,
.permit-info {
  margin: 16px 0;
}

.document-signature {
  display: flex;
  justify-content: space-between;
  margin-top: 40px;
  padding: 0 40px;
}

.signature-left,
.signature-right {
  text-align: center;
}

.signature-space {
  height: 80px;
}

.document-footer {
  margin-top: 40px;
  text-align: center;
}

.empty-state {
  text-align: center;
  padding: 40px;
}

/* Print styles */
@media print {
  .action-bar {
    display: none;
  }
  
  .document-preview {
    border: none;
    padding: 0;
    max-height: none;
    overflow: visible;
  }
}
</style>
