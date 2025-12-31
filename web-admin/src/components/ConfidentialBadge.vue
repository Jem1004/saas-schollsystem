<script setup lang="ts">
import { Tag, Tooltip } from 'ant-design-vue'
import { LockOutlined, SafetyOutlined, EyeInvisibleOutlined } from '@ant-design/icons-vue'

interface Props {
  type?: 'internal' | 'confidential' | 'restricted'
  size?: 'small' | 'default'
  showTooltip?: boolean
  tooltipText?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'confidential',
  size: 'default',
  showTooltip: true,
  tooltipText: '',
})

const badgeConfig = {
  internal: {
    color: 'orange',
    icon: LockOutlined,
    label: 'Internal',
    tooltip: 'Data ini hanya dapat diakses oleh pihak internal yang berwenang',
  },
  confidential: {
    color: 'red',
    icon: SafetyOutlined,
    label: 'Rahasia',
    tooltip: 'Data ini bersifat rahasia dan dilindungi',
  },
  restricted: {
    color: 'purple',
    icon: EyeInvisibleOutlined,
    label: 'Terbatas',
    tooltip: 'Akses ke data ini dibatasi berdasarkan peran',
  },
}

const config = badgeConfig[props.type]
const tooltipContent = props.tooltipText || config.tooltip
</script>

<template>
  <Tooltip v-if="showTooltip" :title="tooltipContent">
    <Tag 
      :color="config.color" 
      :class="['confidential-badge', `size-${size}`]"
    >
      <component :is="config.icon" />
      <span class="badge-label">{{ config.label }}</span>
    </Tag>
  </Tooltip>
  <Tag 
    v-else
    :color="config.color" 
    :class="['confidential-badge', `size-${size}`]"
  >
    <component :is="config.icon" />
    <span class="badge-label">{{ config.label }}</span>
  </Tag>
</template>

<style scoped>
.confidential-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-weight: 500;
}

.confidential-badge.size-small {
  font-size: 11px;
  padding: 0 6px;
  line-height: 18px;
}

.confidential-badge.size-default {
  font-size: 12px;
  padding: 2px 8px;
  line-height: 20px;
}

.badge-label {
  margin-left: 2px;
}
</style>
