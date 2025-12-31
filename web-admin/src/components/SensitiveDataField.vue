<script setup lang="ts">
import { computed, watch } from 'vue'
import { Button, Tooltip, Tag } from 'ant-design-vue'
import {
  EyeOutlined,
  EyeInvisibleOutlined,
  LockOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons-vue'
import { useSensitiveData } from '@/composables/useSensitiveData'

interface Props {
  value: string
  label?: string
  blurByDefault?: boolean
  requireConfirmation?: boolean
  confirmTitle?: string
  confirmDescription?: string
  showIndicator?: boolean
  maskChar?: string
  visibleChars?: number
}

const props = withDefaults(defineProps<Props>(), {
  label: '',
  blurByDefault: true,
  requireConfirmation: true,
  confirmTitle: 'Data Sensitif',
  confirmDescription: 'Data ini bersifat rahasia. Apakah Anda yakin ingin melihatnya?',
  showIndicator: true,
  maskChar: '•',
  visibleChars: 0,
})

const emit = defineEmits<{
  (e: 'reveal'): void
  (e: 'hide'): void
}>()

const { isRevealed, toggle } = useSensitiveData({
  title: props.confirmTitle,
  description: props.confirmDescription,
  requireConfirmation: props.requireConfirmation,
  blurByDefault: props.blurByDefault,
})

const displayText = computed(() => {
  if (isRevealed.value) return props.value
  
  if (!props.value) return ''
  
  if (props.visibleChars > 0) {
    const visible = props.value.slice(0, props.visibleChars)
    const masked = props.maskChar.repeat(Math.min(props.value.length - props.visibleChars, 15))
    return visible + masked
  }
  
  return props.maskChar.repeat(Math.min(props.value.length, 20))
})

const handleToggle = async () => {
  const result = await toggle()
  if (result) {
    if (isRevealed.value) {
      emit('reveal')
    } else {
      emit('hide')
    }
  }
}

watch(isRevealed, (newVal) => {
  if (newVal) {
    emit('reveal')
  } else {
    emit('hide')
  }
})
</script>

<template>
  <div class="sensitive-data-field">
    <div v-if="label" class="field-label">
      <span>{{ label }}</span>
      <Tag v-if="showIndicator && !isRevealed" color="warning" class="confidential-tag">
        <LockOutlined /> Rahasia
      </Tag>
    </div>
    
    <div class="field-content">
      <div 
        class="field-value" 
        :class="{ 'is-blurred': !isRevealed }"
      >
        {{ displayText }}
      </div>
      
      <Tooltip :title="isRevealed ? 'Sembunyikan' : 'Tampilkan'">
        <Button 
          type="text" 
          size="small" 
          class="toggle-btn"
          @click="handleToggle"
        >
          <template #icon>
            <EyeInvisibleOutlined v-if="isRevealed" />
            <EyeOutlined v-else />
          </template>
        </Button>
      </Tooltip>
    </div>
    
    <div v-if="!isRevealed && showIndicator" class="field-hint">
      <ExclamationCircleOutlined />
      <span>Klik ikon mata untuk melihat data</span>
    </div>
  </div>
</template>

<style scoped>
.sensitive-data-field {
  width: 100%;
}

.field-label {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
  font-weight: 500;
  color: #262626;
}

.confidential-tag {
  font-size: 11px;
  padding: 0 6px;
  line-height: 18px;
}

.field-content {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 8px 12px;
  background-color: #fafafa;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
}

.field-value {
  flex: 1;
  word-break: break-word;
  line-height: 1.6;
  transition: filter 0.3s ease;
}

.field-value.is-blurred {
  filter: blur(4px);
  user-select: none;
  color: #8c8c8c;
}

.toggle-btn {
  flex-shrink: 0;
  color: #8c8c8c;
}

.toggle-btn:hover {
  color: #f97316;
}

.field-hint {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 4px;
  font-size: 12px;
  color: #8c8c8c;
}
</style>
