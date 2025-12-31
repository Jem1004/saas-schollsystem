<script setup lang="ts">
import { Alert, Typography, Tooltip } from 'ant-design-vue'
import { EyeOutlined, InfoCircleOutlined, LockOutlined } from '@ant-design/icons-vue'

const { Text } = Typography

interface Props {
  title?: string
  description?: string
  showIcon?: boolean
  type?: 'info' | 'warning'
  compact?: boolean
  tooltipText?: string
}

withDefaults(defineProps<Props>(), {
  title: 'Mode Baca Saja',
  description: 'Anda hanya dapat melihat data ini tanpa melakukan perubahan.',
  showIcon: true,
  type: 'info',
  compact: false,
  tooltipText: '',
})
</script>

<template>
  <div class="read-only-banner" :class="{ compact }">
    <Alert
      v-if="!compact"
      :type="type"
      show-icon
      class="banner-alert"
    >
      <template #icon>
        <EyeOutlined />
      </template>
      <template #message>
        <div class="banner-header">
          <Text strong>{{ title }}</Text>
          <Tooltip v-if="tooltipText" :title="tooltipText">
            <InfoCircleOutlined class="info-icon" />
          </Tooltip>
        </div>
      </template>
      <template #description>
        <slot name="description">
          {{ description }}
        </slot>
      </template>
    </Alert>
    
    <div v-else class="compact-banner" :class="[`type-${type}`]">
      <LockOutlined class="compact-icon" />
      <Text type="secondary">{{ title }}</Text>
      <Tooltip v-if="tooltipText" :title="tooltipText">
        <InfoCircleOutlined class="info-icon" />
      </Tooltip>
    </div>
  </div>
</template>

<style scoped>
.read-only-banner {
  margin-bottom: 16px;
}

.read-only-banner.compact {
  margin-bottom: 8px;
}

.banner-alert {
  border-radius: 8px;
}

.banner-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.info-icon {
  color: #8c8c8c;
  cursor: help;
}

.compact-banner {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 12px;
}

.compact-banner.type-info {
  background-color: #e6f7ff;
  border: 1px solid #91d5ff;
}

.compact-banner.type-warning {
  background-color: #fffbe6;
  border: 1px solid #ffe58f;
}

.compact-icon {
  font-size: 12px;
}
</style>
