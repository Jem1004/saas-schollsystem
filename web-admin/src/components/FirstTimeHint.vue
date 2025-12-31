<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Alert, Typography } from 'ant-design-vue'
import { BulbOutlined } from '@ant-design/icons-vue'

const { Text } = Typography

interface Props {
  hintKey: string
  title?: string
  message: string
  type?: 'info' | 'success' | 'warning'
  closable?: boolean
  showOnce?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  title: 'Tips',
  type: 'info',
  closable: true,
  showOnce: true,
})

const emit = defineEmits<{
  (e: 'close'): void
}>()

const isVisible = ref(true)

const storageKey = `hint_dismissed_${props.hintKey}`

onMounted(() => {
  if (props.showOnce) {
    const dismissed = localStorage.getItem(storageKey)
    if (dismissed === 'true') {
      isVisible.value = false
    }
  }
})

const handleClose = () => {
  isVisible.value = false
  if (props.showOnce) {
    localStorage.setItem(storageKey, 'true')
  }
  emit('close')
}
</script>

<template>
  <Alert
    v-if="isVisible"
    :type="type"
    :closable="closable"
    class="first-time-hint"
    @close="handleClose"
  >
    <template #icon>
      <BulbOutlined />
    </template>
    <template #message>
      <Text strong>{{ title }}</Text>
    </template>
    <template #description>
      <slot>{{ message }}</slot>
    </template>
  </Alert>
</template>

<style scoped>
.first-time-hint {
  margin-bottom: 16px;
  border-radius: 8px;
}

.first-time-hint :deep(.ant-alert-icon) {
  font-size: 18px;
}
</style>
