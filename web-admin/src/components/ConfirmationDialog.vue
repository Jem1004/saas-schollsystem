<script setup lang="ts">
import { ref, watch } from 'vue'
import { Modal, Typography, Input, Alert } from 'ant-design-vue'
import {
  ExclamationCircleOutlined,
  DeleteOutlined,
  LockOutlined,
} from '@ant-design/icons-vue'

const { Text, Paragraph } = Typography

interface Props {
  open: boolean
  title?: string
  message?: string
  description?: string
  type?: 'warning' | 'danger' | 'sensitive' | 'info'
  confirmText?: string
  cancelText?: string
  requireInput?: boolean
  inputPlaceholder?: string
  inputValidation?: string
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  title: 'Konfirmasi',
  message: 'Apakah Anda yakin ingin melanjutkan?',
  description: '',
  type: 'warning',
  confirmText: 'Ya, Lanjutkan',
  cancelText: 'Batal',
  requireInput: false,
  inputPlaceholder: 'Ketik untuk konfirmasi',
  inputValidation: '',
  loading: false,
})

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'confirm'): void
  (e: 'cancel'): void
}>()

const inputValue = ref('')
const isInputValid = ref(!props.requireInput)

const typeConfig = {
  warning: {
    icon: ExclamationCircleOutlined,
    iconColor: '#faad14',
    okType: 'primary' as const,
  },
  danger: {
    icon: DeleteOutlined,
    iconColor: '#ff4d4f',
    okType: 'primary' as const,
    okDanger: true,
  },
  sensitive: {
    icon: LockOutlined,
    iconColor: '#fa8c16',
    okType: 'primary' as const,
  },
  info: {
    icon: ExclamationCircleOutlined,
    iconColor: '#1890ff',
    okType: 'primary' as const,
  },
}

const config = typeConfig[props.type]

watch(inputValue, (val) => {
  if (props.requireInput && props.inputValidation) {
    isInputValid.value = val === props.inputValidation
  } else if (props.requireInput) {
    isInputValid.value = val.length > 0
  }
})

watch(() => props.open, (val) => {
  if (!val) {
    inputValue.value = ''
    isInputValid.value = !props.requireInput
  }
})

const handleOk = () => {
  if (props.requireInput && !isInputValid.value) return
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
    :ok-type="config.okType"
    :ok-button-props="{ 
      danger: type === 'danger',
      disabled: requireInput && !isInputValid 
    }"
    centered
    @ok="handleOk"
    @cancel="handleCancel"
  >
    <div class="confirmation-content">
      <div class="confirmation-header">
        <div class="icon-wrapper" :style="{ backgroundColor: `${config.iconColor}15` }">
          <component :is="config.icon" :style="{ color: config.iconColor, fontSize: '24px' }" />
        </div>
        <div class="header-text">
          <Text strong class="confirmation-title">{{ title }}</Text>
          <Paragraph class="confirmation-message">{{ message }}</Paragraph>
        </div>
      </div>

      <div v-if="description" class="confirmation-description">
        <Alert :type="type === 'danger' ? 'error' : type === 'sensitive' ? 'warning' : 'info'" show-icon>
          <template #message>{{ description }}</template>
        </Alert>
      </div>

      <div v-if="requireInput" class="confirmation-input">
        <Text type="secondary" class="input-label">
          <template v-if="inputValidation">
            Ketik <Text code>{{ inputValidation }}</Text> untuk konfirmasi:
          </template>
          <template v-else>
            {{ inputPlaceholder }}
          </template>
        </Text>
        <Input
          v-model:value="inputValue"
          :placeholder="inputPlaceholder"
          :status="inputValue && !isInputValid ? 'error' : ''"
          class="confirmation-input-field"
        />
      </div>
    </div>
  </Modal>
</template>

<style scoped>
.confirmation-content {
  padding: 8px 0;
}

.confirmation-header {
  display: flex;
  gap: 16px;
  align-items: flex-start;
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

.confirmation-title {
  font-size: 16px;
  display: block;
  margin-bottom: 4px;
}

.confirmation-message {
  color: #595959;
  margin-bottom: 0;
}

.confirmation-description {
  margin-top: 16px;
}

.confirmation-input {
  margin-top: 16px;
}

.input-label {
  display: block;
  margin-bottom: 8px;
}

.confirmation-input-field {
  width: 100%;
}
</style>
