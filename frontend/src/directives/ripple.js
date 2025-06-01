// 波纹效果指令
export default {
  mounted(el, binding) {
    // 添加必要的类名和样式
    el.classList.add('ripple-container')
    el.style.position = 'relative'
    el.style.overflow = 'hidden'
    
    // 点击事件处理
    const handleClick = (e) => {
      // 获取点击位置
      const rect = el.getBoundingClientRect()
      const x = e.clientX - rect.left
      const y = e.clientY - rect.top
      
      // 创建波纹元素
      const ripple = document.createElement('span')
      ripple.className = 'ripple-effect'
      ripple.style.left = x + 'px'
      ripple.style.top = y + 'px'
      
      // 设置波纹大小
      const size = Math.max(rect.width, rect.height) * 2
      ripple.style.width = size + 'px'
      ripple.style.height = size + 'px'
      
      // 添加到元素中
      el.appendChild(ripple)
      
      // 动画结束后移除
      setTimeout(() => {
        ripple.remove()
      }, 600)
    }
    
    // 绑定事件
    el._rippleHandler = handleClick
    el.addEventListener('click', handleClick)
  },
  
  unmounted(el) {
    // 清理事件监听
    if (el._rippleHandler) {
      el.removeEventListener('click', el._rippleHandler)
      delete el._rippleHandler
    }
  }
}

// 波纹效果样式（需要添加到全局样式中）
const rippleStyles = `
.ripple-container {
  position: relative;
  overflow: hidden;
}

.ripple-effect {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.5);
  transform: translate(-50%, -50%) scale(0);
  animation: ripple-animation 0.6s ease-out;
  pointer-events: none;
}

@keyframes ripple-animation {
  to {
    transform: translate(-50%, -50%) scale(1);
    opacity: 0;
  }
}
`

// 自动注入样式
if (typeof document !== 'undefined') {
  const style = document.createElement('style')
  style.textContent = rippleStyles
  document.head.appendChild(style)
}