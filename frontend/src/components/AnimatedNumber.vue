<template>
  <span class="animated-number">
    {{ displayValue }}
  </span>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'

const props = defineProps({
  value: {
    type: Number,
    required: true
  },
  duration: {
    type: Number,
    default: 1000
  },
  format: {
    type: Function,
    default: (val) => val.toFixed(0)
  },
  delay: {
    type: Number,
    default: 0
  },
  easing: {
    type: String,
    default: 'easeOutCubic'
  },
  prefix: {
    type: String,
    default: ''
  },
  suffix: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['complete'])

const currentValue = ref(0)
const startTimestamp = ref(null)
const startValue = ref(0)

// 缓动函数
const easingFunctions = {
  linear: (t) => t,
  easeInQuad: (t) => t * t,
  easeOutQuad: (t) => t * (2 - t),
  easeInOutQuad: (t) => t < 0.5 ? 2 * t * t : -1 + (4 - 2 * t) * t,
  easeInCubic: (t) => t * t * t,
  easeOutCubic: (t) => (--t) * t * t + 1,
  easeInOutCubic: (t) => t < 0.5 ? 4 * t * t * t : (t - 1) * (2 * t - 2) * (2 * t - 2) + 1,
  easeInQuart: (t) => t * t * t * t,
  easeOutQuart: (t) => 1 - (--t) * t * t * t,
  easeInOutQuart: (t) => t < 0.5 ? 8 * t * t * t * t : 1 - 8 * (--t) * t * t * t,
  easeInQuint: (t) => t * t * t * t * t,
  easeOutQuint: (t) => 1 + (--t) * t * t * t * t,
  easeInOutQuint: (t) => t < 0.5 ? 16 * t * t * t * t * t : 1 + 16 * (--t) * t * t * t * t
}

const displayValue = computed(() => {
  try {
    let formattedValue
    // 确保 format 是一个函数
    if (typeof props.format === 'function') {
      formattedValue = props.format(currentValue.value)
    } else {
      // 如果不是函数，使用默认格式化
      formattedValue = currentValue.value.toFixed(0)
    }
    
    // 添加前缀和后缀
    return `${props.prefix}${formattedValue}${props.suffix}`
  } catch (error) {
    console.warn('AnimatedNumber format error:', error)
    // 出错时返回原始值
    return `${props.prefix}${currentValue.value || 0}${props.suffix}`
  }
})

const animate = (timestamp) => {
  if (!startTimestamp.value) {
    startTimestamp.value = timestamp
  }
  
  const progress = Math.min((timestamp - startTimestamp.value) / props.duration, 1)
  const easedProgress = easingFunctions[props.easing](progress)
  
  currentValue.value = startValue.value + (props.value - startValue.value) * easedProgress
  
  if (progress < 1) {
    requestAnimationFrame(animate)
  } else {
    currentValue.value = props.value
    emit('complete')
  }
}

const startAnimation = () => {
  startValue.value = currentValue.value
  startTimestamp.value = null
  
  if (props.delay > 0) {
    setTimeout(() => {
      requestAnimationFrame(animate)
    }, props.delay)
  } else {
    requestAnimationFrame(animate)
  }
}

watch(() => props.value, () => {
  startAnimation()
})

onMounted(() => {
  startAnimation()
})
</script>

<style lang="scss" scoped>
.animated-number {
  display: inline-block;
  font-variant-numeric: tabular-nums;
}
</style>