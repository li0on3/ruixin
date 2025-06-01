// 瑞幸分销系统 - 动画插件
// ============================================

import { nextTick } from 'vue'
import LoadingAnimation from '@/components/LoadingAnimation.vue'
import ResultAnimation from '@/components/ResultAnimation.vue'
import AnimatedNumber from '@/components/AnimatedNumber.vue'
import rippleDirective from '@/directives/ripple'

/**
 * 动画管理器类
 */
class AnimationManager {
  constructor() {
    this.observers = new Map()
    this.config = {
      enableAnimations: true,
      enableIntersectionObserver: true,
      animationDuration: 300,
      animationEasing: 'cubic-bezier(0.4, 0, 0.2, 1)',
      respectPreferredMotion: true
    }
    this.init()
  }

  init() {
    // 检查用户偏好设置
    if (this.config.respectPreferredMotion) {
      this.checkMotionPreferences()
    }

    // 设置全局动画CSS变量
    this.setupGlobalAnimations()

    // 初始化Intersection Observer
    if (this.config.enableIntersectionObserver) {
      this.setupIntersectionObserver()
    }
  }

  checkMotionPreferences() {
    const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)')
    
    if (prefersReducedMotion.matches) {
      this.config.enableAnimations = false
      this.config.animationDuration = 0
      document.documentElement.style.setProperty('--animation-duration', '0ms')
    }

    // 监听偏好设置变化
    prefersReducedMotion.addEventListener('change', (e) => {
      this.config.enableAnimations = !e.matches
      this.config.animationDuration = e.matches ? 0 : 300
      document.documentElement.style.setProperty(
        '--animation-duration', 
        e.matches ? '0ms' : '300ms'
      )
    })
  }

  setupGlobalAnimations() {
    const root = document.documentElement
    root.style.setProperty('--animation-duration', `${this.config.animationDuration}ms`)
    root.style.setProperty('--animation-easing', this.config.animationEasing)
  }

  setupIntersectionObserver() {
    // 创建Intersection Observer用于触发入场动画
    this.intersectionObserver = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            this.triggerElementAnimation(entry.target)
          }
        })
      },
      {
        threshold: 0.1,
        rootMargin: '0px 0px -50px 0px'
      }
    )
  }

  triggerElementAnimation(element) {
    if (!this.config.enableAnimations) return

    const animationType = element.dataset.animation || 'fade-in'
    const delay = element.dataset.animationDelay || '0'
    
    setTimeout(() => {
      element.classList.add('animate-in')
      element.classList.add(`animate-${animationType}`)
    }, parseInt(delay))

    // 动画完成后移除监听
    this.intersectionObserver?.unobserve(element)
  }

  // 注册元素进行动画监听
  observeElement(element) {
    if (this.intersectionObserver) {
      this.intersectionObserver.observe(element)
    }
  }

  // 取消元素动画监听
  unobserveElement(element) {
    if (this.intersectionObserver) {
      this.intersectionObserver.unobserve(element)
    }
  }

  // 手动触发动画
  animate(element, animation, options = {}) {
    if (!this.config.enableAnimations) return Promise.resolve()

    const {
      duration = this.config.animationDuration,
      delay = 0,
      easing = this.config.animationEasing,
      fillMode = 'both'
    } = options

    return new Promise((resolve) => {
      const animationName = `animate-${animation}`
      
      element.style.animation = `${animationName} ${duration}ms ${easing} ${delay}ms ${fillMode}`
      
      const onAnimationEnd = () => {
        element.removeEventListener('animationend', onAnimationEnd)
        element.style.animation = ''
        resolve()
      }

      element.addEventListener('animationend', onAnimationEnd)
    })
  }

  // 数字滚动动画
  animateNumber(element, startValue, endValue, options = {}) {
    if (!this.config.enableAnimations) {
      element.textContent = endValue.toString()
      return Promise.resolve()
    }

    const {
      duration = 1000,
      formatter = (value) => Math.round(value).toString(),
      easing = 'easeOutQuart'
    } = options

    return new Promise((resolve) => {
      const startTime = performance.now()
      const range = endValue - startValue

      const easingFunctions = {
        linear: t => t,
        easeOutQuart: t => 1 - Math.pow(1 - t, 4),
        easeInOutQuart: t => t < 0.5 ? 8 * t * t * t * t : 1 - Math.pow(-2 * t + 2, 4) / 2
      }

      const easingFunc = easingFunctions[easing] || easingFunctions.easeOutQuart

      const updateNumber = (currentTime) => {
        const elapsed = currentTime - startTime
        const progress = Math.min(elapsed / duration, 1)
        const easedProgress = easingFunc(progress)
        const currentValue = startValue + (range * easedProgress)

        element.textContent = formatter(currentValue)

        if (progress < 1) {
          requestAnimationFrame(updateNumber)
        } else {
          element.textContent = formatter(endValue)
          resolve()
        }
      }

      requestAnimationFrame(updateNumber)
    })
  }

  // 销毁动画管理器
  destroy() {
    if (this.intersectionObserver) {
      this.intersectionObserver.disconnect()
    }
    this.observers.clear()
  }
}

// 创建全局动画管理器
const animationManager = new AnimationManager()

// Vue指令：自动触发入场动画
const vAnimateOnScroll = {
  mounted(el, binding) {
    const { value = 'fade-in', modifiers } = binding
    
    el.dataset.animation = value
    if (modifiers.delay) {
      el.dataset.animationDelay = '100'
    }
    
    animationManager.observeElement(el)
  },
  
  unmounted(el) {
    animationManager.unobserveElement(el)
  }
}

// Vue指令：悬停动画
const vHoverAnimation = {
  mounted(el, binding) {
    const animation = binding.value || 'hover-lift'
    let isAnimating = false

    const onMouseEnter = () => {
      if (isAnimating) return
      isAnimating = true
      el.classList.add(animation)
    }

    const onMouseLeave = () => {
      el.classList.remove(animation)
      isAnimating = false
    }

    el.addEventListener('mouseenter', onMouseEnter)
    el.addEventListener('mouseleave', onMouseLeave)
    
    el._hoverAnimation = { onMouseEnter, onMouseLeave }
  },
  
  unmounted(el) {
    if (el._hoverAnimation) {
      el.removeEventListener('mouseenter', el._hoverAnimation.onMouseEnter)
      el.removeEventListener('mouseleave', el._hoverAnimation.onMouseLeave)
      delete el._hoverAnimation
    }
  }
}

// Vue插件
export default {
  install(app, options = {}) {
    // 合并配置
    Object.assign(animationManager.config, options)

    // 注册全局组件
    app.component('LoadingAnimation', LoadingAnimation)
    app.component('ResultAnimation', ResultAnimation)
    app.component('AnimatedNumber', AnimatedNumber)

    // 注册指令
    app.directive('ripple', rippleDirective)
    app.directive('animate-on-scroll', vAnimateOnScroll)
    app.directive('hover-animation', vHoverAnimation)

    // 提供动画管理器
    app.provide('animations', animationManager)
    app.config.globalProperties.$animations = animationManager

    // 添加全局方法
    app.config.globalProperties.$animate = animationManager.animate.bind(animationManager)
    app.config.globalProperties.$animateNumber = animationManager.animateNumber.bind(animationManager)

    // 添加页面切换动画类到根元素
    app.mixin({
      mounted() {
        if (this.$el && this.$el.classList && this.$parent && !this.$parent.$parent) {
          this.$el.classList.add('animate-fade-in')
          setTimeout(() => {
            this.$el.classList.remove('animate-fade-in')
          }, 300)
        }
      }
    })

    // 添加全局过渡动画配置
    app.config.globalProperties.$transitions = {
      page: 'page-transition',
      modal: 'modal-bounce-in',
      drawer: 'slide-right',
      collapse: 'collapse',
      list: 'list'
    }

    // 开发环境下暴露到 window
    if (process.env.NODE_ENV === 'development') {
      window.__ANIMATION_MANAGER__ = animationManager
    }
  }
}

// 导出动画管理器
export { animationManager }

// 组合式API
export function useAnimations() {
  return {
    animationManager,
    animate: animationManager.animate.bind(animationManager),
    animateNumber: animationManager.animateNumber.bind(animationManager)
  }
}