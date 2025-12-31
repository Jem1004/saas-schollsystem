<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { Modal, Input, Typography, Alert } from 'ant-design-vue'
import { DeleteOutlined, WarningOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue'

const { Text, Paragraph } = Typography

interface Props {
  open: boolean
  title?: string
  message?: string
  warningMessage?: string
  confirmText?: string
  cancelText?: string
  requireConfirmation?: boolean
  confirmationWord?: string
  loading?: boolean
  type?: 'delete' | 'warning' | 'danger'
}

const props = withDefaults(defineProps<Props>(), {
  title: 'Konfirmasi Hapus',
  message: 'Apakah Anda yakin ingin menghapus data ini?',
  warningMessage: 'Tindakan ini tidak dapat dibatalkan.',
  confirmText: 'Ya, Hapus',
  cancelText: 'Batal',
  requireConfirmation: false,
  confirmationWord: 'HAPUS',
  loading: false,
  type: 'delete',
})

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'confirm'): void
  (e: 'cancel'): void
}>()

const inputValue = ref('')

const isConfirmEnabled = computed(() => {
  if (!props.requireConfirmation) return true
  return inputValue.value === props.confirmationWord
})

const typeConfig = {
  delete: {
    icon: DeleteOutlined,
    iconColor: '#ff4d4f',
    alertType: 'error' as const,
  },
  warning: {
    icon: WarningOutlined,
    iconColor: '#faad14',
    alertType: 'warning' as const,
  },
  danger: {
    icon: ExclamationCircleOutlined,
    iconColor: '#ff4d4f',
    alertType: 'error' as const,
  },
}

const config = typeConfig[props.type]

watch(() => props.open, (val) => {
  if (!val) {
    inputValue.value = ''
  }
})

const handleOk = () => {
  if (!isConfirmEnabled.value) return
  emit('confirm')
}

const handleCancel = () => {
  emit('update:open', false)
  emit('cancel')
}
</script>

<template>
  <Modal
    :open="open"
    :title="null"
    :confirm-loading="loading"
    :ok-text="confirmText"
    :cancel-text="cancelText"
    ok-type="primary"
    :ok-button-props="{ danger: true, disabled: !isConfirmEnabled }"
    centered
    @ok="handleOk"
    @cancel="handleCancel"
  >
    <div class="destructive-dialog">
      <div class="dialog-header">
        <div class="icon-wrapper" :style="{ backgroundColor: `${config.iconColor}15` }">
          <component :is="config.icon" :style="{ color: config.iconColor, fontSize: '24px' }" />
        </div>
        <div class="header-text">
          <Text strong class="dialog-title">{{ title }}</Text>
          <Paragraph class="dialog-message">{{ message }}</Paragraph>
        </div>
      </div>

      <Alert
        :type="config.alertType"
        :message="warningMessage"
        show-icon
        class="warning-alert"
      />

      <div v-if="requireConfirmation" class="confirmation-input">
        <Text type="secondary">
          Ketik <Text code>{{ confirmationWord }}</Text> untuk konfirmasi:
        </Text>
        <Input
          v-model:value="inputValue"
          :placeholder="`Ketik ${confirmationWord}`"
          :status="inputValue && !isConfirmEnabled ? 'error' : ''"
          class="input-field"
        />
      </div>
    </div>
  </Modal>
</template>

<style scoped>
.destructive-dialog {
  padding: 8px 0;
}

.dialog-header {
  display: flex;
  gap: 16px;
  align-items: flex-start;
  margin-bottom: 16px;
}

.icon-wrapper {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.header-text {
  flex: 1;
}

.dialog-title {
  font-size: 16px;
  display: block;
  margin-bottom: 4px;
}

.dialog-message {
  color: #595959;
  margin-bottom: 0;
}

.warning-alert {
  margin-bottom: 16px;
}

.confirmation-input {
  margin-top: 16px;
}

.confirmation-input > span:first-child {
  display: block;
  margin-bottom: 8px;
}

.input-field {
  width: 100%;
}
</style>
