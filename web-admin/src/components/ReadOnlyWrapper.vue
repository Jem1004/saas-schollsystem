<script setup lang="ts">
import { computed } from 'vue'
import { Tooltip } from 'ant-design-vue'

interface Props {
  readOnly?: boolean
  tooltipText?: string
  disableInputs?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  readOnly: false,
  tooltipText: 'Anda tidak memiliki izin untuk mengubah data ini',
  disableInputs: true,
})

const wrapperClass = computed(() => ({
  'read-only-wrapper': true,
  'is-read-only': props.readOnly,
  'disable-inputs': props.readOnly && props.disableInputs,
}))
</script>

<template>
  <Tooltip v-if="readOnly" :title="tooltipText" placement="top">
    <div :class="wrapperClass">
      <slot />
      <div v-if="readOnly" class="read-only-overlay" />
    </div>
  </Tooltip>
  <div v-else :class="wrapperClass">
    <slot />
  </div>
</template>

<style scoped>
.read-only-wrapper {
  position: relative;
}

.read-only-wrapper.is-read-only {
  cursor: not-allowed;
}

.read-only-wrapper.disable-inputs :deep(input),
.read-only-wrapper.disable-inputs :deep(textarea),
.read-only-wrapper.disable-inputs :deep(select),
.read-only-wrapper.disable-inputs :deep(button:not(.ant-btn-link)),
.read-only-wrapper.disable-inputs :deep(.ant-select),
.read-only-wrapper.disable-inputs :deep(.ant-picker),
.read-only-wrapper.disable-inputs :deep(.ant-input),
.read-only-wrapper.disable-inputs :deep(.ant-input-number) {
  pointer-events: none;
  opacity: 0.7;
}

.read-only-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: transparent;
  z-index: 1;
}
</style>
