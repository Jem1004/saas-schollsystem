<script setup lang="ts">
import { computed } from 'vue'
import { DatePicker, RangePicker } from 'ant-design-vue'
import type { Dayjs } from 'dayjs'
import dayjs from 'dayjs'
import 'dayjs/locale/id'
import localeData from 'dayjs/plugin/localeData'
import weekday from 'dayjs/plugin/weekday'

// Setup dayjs plugins and locale
dayjs.extend(localeData)
dayjs.extend(weekday)
dayjs.locale('id')

type DateValue = Dayjs | string | null
type RangeValue = [Dayjs, Dayjs] | [string, string] | null

interface Props {
  modelValue?: DateValue | RangeValue
  type?: 'date' | 'week' | 'month' | 'quarter' | 'year' | 'range'
  placeholder?: string | [string, string]
  format?: string
  showTime?: boolean
  disabled?: boolean
  allowClear?: boolean
  disabledDate?: (current: Dayjs) => boolean
  valueFormat?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'date',
  format: 'DD MMMM YYYY',
  showTime: false,
  disabled: false,
  allowClear: true,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: DateValue | RangeValue): void
  (e: 'change', value: DateValue | RangeValue, dateString: string | [string, string]): void
}>()

// Indonesian locale configuration
const idLocale = {
  lang: {
    locale: 'id_ID',
    placeholder: 'Pilih tanggal',
    rangePlaceholder: ['Tanggal mulai', 'Tanggal akhir'] as [string, string],
    today: 'Hari ini',
    now: 'Sekarang',
    backToToday: 'Kembali ke hari ini',
    ok: 'OK',
    clear: 'Hapus',
    month: 'Bulan',
    year: 'Tahun',
    timeSelect: 'Pilih waktu',
    dateSelect: 'Pilih tanggal',
    monthSelect: 'Pilih bulan',
    yearSelect: 'Pilih tahun',
    decadeSelect: 'Pilih dekade',
    yearFormat: 'YYYY',
    dateFormat: 'D/M/YYYY',
    dayFormat: 'D',
    dateTimeFormat: 'D/M/YYYY HH:mm:ss',
    monthFormat: 'MMMM',
    monthBeforeYear: true,
    previousMonth: 'Bulan sebelumnya',
    nextMonth: 'Bulan berikutnya',
    previousYear: 'Tahun sebelumnya',
    nextYear: 'Tahun berikutnya',
    previousDecade: 'Dekade sebelumnya',
    nextDecade: 'Dekade berikutnya',
    previousCentury: 'Abad sebelumnya',
    nextCentury: 'Abad berikutnya',
    shortWeekDays: ['Min', 'Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab'],
    shortMonths: ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des'],
  },
  timePickerLocale: {
    placeholder: 'Pilih waktu',
  },
}

// Computed format based on type
const computedFormat = computed(() => {
  if (props.format !== 'DD MMMM YYYY') return props.format
  
  switch (props.type) {
    case 'week':
      return 'YYYY-wo'
    case 'month':
      return 'MMMM YYYY'
    case 'quarter':
      return 'YYYY-[Q]Q'
    case 'year':
      return 'YYYY'
    default:
      return props.showTime ? 'DD MMMM YYYY HH:mm' : 'DD MMMM YYYY'
  }
})

// Computed placeholder
const computedPlaceholder = computed(() => {
  if (props.placeholder) return props.placeholder
  
  if (props.type === 'range') {
    return ['Tanggal mulai', 'Tanggal akhir'] as [string, string]
  }
  
  switch (props.type) {
    case 'week':
      return 'Pilih minggu'
    case 'month':
      return 'Pilih bulan'
    case 'quarter':
      return 'Pilih kuartal'
    case 'year':
      return 'Pilih tahun'
    default:
      return 'Pilih tanggal'
  }
})

// Convert value to Dayjs
const internalValue = computed(() => {
  if (!props.modelValue) return undefined
  
  if (props.type === 'range') {
    const rangeVal = props.modelValue as RangeValue
    if (!rangeVal) return undefined
    return [
      typeof rangeVal[0] === 'string' ? dayjs(rangeVal[0]) : rangeVal[0],
      typeof rangeVal[1] === 'string' ? dayjs(rangeVal[1]) : rangeVal[1],
    ] as [Dayjs, Dayjs]
  }
  
  const val = props.modelValue as DateValue
  return typeof val === 'string' ? dayjs(val) : (val ?? undefined)
})

// Handle change for range picker
const handleRangeChange = (value: unknown, dateString: [string, string]) => {
  const typedValue = value as [Dayjs, Dayjs] | null
  let emitValue: RangeValue = typedValue
  
  if (props.valueFormat && typedValue) {
    emitValue = [
      typedValue[0].format(props.valueFormat),
      typedValue[1].format(props.valueFormat),
    ] as [string, string]
  }
  
  emit('update:modelValue', emitValue)
  emit('change', emitValue, dateString)
}

// Handle change for date picker
const handleDateChange = (value: unknown, dateString: string) => {
  const typedValue = value as Dayjs | null
  let emitValue: DateValue = typedValue
  
  if (props.valueFormat && typedValue) {
    emitValue = typedValue.format(props.valueFormat)
  }
  
  emit('update:modelValue', emitValue)
  emit('change', emitValue, dateString)
}

// Get picker type
const pickerType = computed(() => {
  if (props.type === 'range') return 'date'
  return props.type
})
</script>

<template>
  <RangePicker
    v-if="type === 'range'"
    :value="internalValue as [Dayjs, Dayjs] | undefined"
    :format="computedFormat"
    :placeholder="computedPlaceholder as [string, string]"
    :show-time="showTime"
    :disabled="disabled"
    :allow-clear="allowClear"
    :disabled-date="disabledDate"
    :locale="idLocale"
    style="width: 100%"
    @change="handleRangeChange"
  />
  <DatePicker
    v-else
    :value="internalValue as Dayjs | undefined"
    :picker="pickerType"
    :format="computedFormat"
    :placeholder="computedPlaceholder as string"
    :show-time="showTime"
    :disabled="disabled"
    :allow-clear="allowClear"
    :disabled-date="disabledDate"
    :locale="idLocale"
    style="width: 100%"
    @change="handleDateChange"
  />
</template>
