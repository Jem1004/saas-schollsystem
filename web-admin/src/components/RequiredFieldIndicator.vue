<script setup lang="ts">
import { Tooltip } from 'ant-design-vue'

interface Props {
  showTooltip?: boolean
  tooltipText?: string
}

withDefaults(defineProps<Props>(), {
  showTooltip: true,
  tooltipText: 'Field ini wajib diisi',
})
</script>

<template>
  <Tooltip v-if="showTooltip" :title="tooltipText">
    <span class="required-indicator">*</span>
  </Tooltip>
  <span v-else class="required-indicator">*</span>
</template>

<style scoped>
.required-indicator {
  color: #ff4d4f;
  font-weight: bold;
  margin-left: 2px;
}
</style>
