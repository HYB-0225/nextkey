/**
 * NextKey 动画工具函数库
 * 基于GSAP实现高性能动画效果
 */

import gsap from 'gsap'

/**
 * 默认动画配置
 */
export const ANIMATION_CONFIG = {
  // 时长
  duration: {
    fast: 0.15,
    normal: 0.3,
    slow: 0.5,
  },
  
  // 缓动函数
  ease: {
    out: 'power2.out',
    inOut: 'power2.inOut',
    elastic: 'elastic.out(1, 0.5)',
    back: 'back.out(1.7)',
    bounce: 'bounce.out',
  },
  
  // 交错延迟
  stagger: 0.05,
}

/**
 * 列表交错动画 - 滑入效果
 * @param {Array|NodeList} elements - 要动画的元素列表
 * @param {Object} options - 配置选项
 */
export function staggerFadeInUp(elements, options = {}) {
  const {
    duration = ANIMATION_CONFIG.duration.normal,
    stagger = ANIMATION_CONFIG.stagger,
    ease = ANIMATION_CONFIG.ease.out,
    y = 20,
    clearProps = 'all',
  } = options

  return gsap.fromTo(
    elements,
    {
      opacity: 0,
      y: y,
    },
    {
      opacity: 1,
      y: 0,
      duration,
      stagger,
      ease,
      clearProps,
    }
  )
}

/**
 * 列表交错动画 - 缩放弹入效果
 * @param {Array|NodeList} elements - 要动画的元素列表
 * @param {Object} options - 配置选项
 */
export function staggerScaleIn(elements, options = {}) {
  const {
    duration = ANIMATION_CONFIG.duration.normal,
    stagger = ANIMATION_CONFIG.stagger,
    ease = ANIMATION_CONFIG.ease.back,
    scale = 0.8,
    clearProps = 'all',
  } = options

  return gsap.fromTo(
    elements,
    {
      opacity: 0,
      scale: scale,
    },
    {
      opacity: 1,
      scale: 1,
      duration,
      stagger,
      ease,
      clearProps,
    }
  )
}

/**
 * 弹出动画 - 适用于对话框
 * @param {Element} element - 要动画的元素
 * @param {Object} options - 配置选项
 */
export function popIn(element, options = {}) {
  const {
    duration = ANIMATION_CONFIG.duration.normal,
    ease = ANIMATION_CONFIG.ease.back,
    scale = 0.9,
  } = options

  return gsap.fromTo(
    element,
    {
      opacity: 0,
      scale: scale,
    },
    {
      opacity: 1,
      scale: 1,
      duration,
      ease,
      clearProps: 'all',
    }
  )
}

/**
 * 弹出动画 - 退出效果
 * @param {Element} element - 要动画的元素
 * @param {Object} options - 配置选项
 */
export function popOut(element, options = {}) {
  const {
    duration = ANIMATION_CONFIG.duration.fast,
    ease = ANIMATION_CONFIG.ease.inOut,
    scale = 0.9,
  } = options

  return gsap.to(element, {
    opacity: 0,
    scale: scale,
    duration,
    ease,
  })
}

/**
 * 底部滑入动画 - 移动端对话框
 * @param {Element} element - 要动画的元素
 * @param {Object} options - 配置选项
 */
export function slideUpIn(element, options = {}) {
  const {
    duration = ANIMATION_CONFIG.duration.normal,
    ease = ANIMATION_CONFIG.ease.out,
    y = '100%',
  } = options

  return gsap.fromTo(
    element,
    {
      y: y,
    },
    {
      y: 0,
      duration,
      ease,
      clearProps: 'all',
    }
  )
}

/**
 * 底部滑出动画
 * @param {Element} element - 要动画的元素
 * @param {Object} options - 配置选项
 */
export function slideDownOut(element, options = {}) {
  const {
    duration = ANIMATION_CONFIG.duration.normal,
    ease = ANIMATION_CONFIG.ease.inOut,
    y = '100%',
  } = options

  return gsap.to(element, {
    y: y,
    duration,
    ease,
  })
}

/**
 * 淡入动画
 * @param {Element|Array} elements - 要动画的元素
 * @param {Object} options - 配置选项
 */
export function fadeIn(elements, options = {}) {
  const {
    duration = ANIMATION_CONFIG.duration.normal,
    ease = ANIMATION_CONFIG.ease.out,
    stagger = 0,
  } = options

  return gsap.fromTo(
    elements,
    {
      opacity: 0,
    },
    {
      opacity: 1,
      duration,
      ease,
      stagger,
      clearProps: 'all',
    }
  )
}

/**
 * 淡出动画
 * @param {Element|Array} elements - 要动画的元素
 * @param {Object} options - 配置选项
 */
export function fadeOut(elements, options = {}) {
  const {
    duration = ANIMATION_CONFIG.duration.fast,
    ease = ANIMATION_CONFIG.ease.inOut,
  } = options

  return gsap.to(elements, {
    opacity: 0,
    duration,
    ease,
  })
}

/**
 * 删除动画 - 缩小淡出
 * @param {Element} element - 要动画的元素
 * @param {Function} onComplete - 完成回调
 */
export function removeAnimation(element, onComplete) {
  return gsap.to(element, {
    scale: 0.8,
    opacity: 0,
    duration: ANIMATION_CONFIG.duration.fast,
    ease: ANIMATION_CONFIG.ease.inOut,
    onComplete,
  })
}

/**
 * 抖动动画 - 用于错误提示
 * @param {Element} element - 要动画的元素
 * @param {Object} options - 配置选项
 */
export function shake(element, options = {}) {
  const {
    duration = 0.5,
    strength = 10,
  } = options

  return gsap.fromTo(
    element,
    { x: -strength },
    {
      x: strength,
      duration: 0.1,
      repeat: 5,
      yoyo: true,
      ease: 'power1.inOut',
      onComplete: () => {
        gsap.set(element, { x: 0 })
      },
    }
  )
}

/**
 * 脉冲动画 - 用于加载状态
 * @param {Element} element - 要动画的元素
 */
export function pulse(element) {
  return gsap.to(element, {
    scale: 1.1,
    duration: 0.6,
    repeat: -1,
    yoyo: true,
    ease: 'power1.inOut',
  })
}

/**
 * 涟漪效果 - 按钮点击反馈
 * @param {Event} event - 点击事件
 * @param {Element} container - 容器元素
 */
export function ripple(event, container) {
  const ripple = document.createElement('span')
  const rect = container.getBoundingClientRect()
  const size = Math.max(rect.width, rect.height)
  const x = event.clientX - rect.left - size / 2
  const y = event.clientY - rect.top - size / 2

  ripple.style.cssText = `
    position: absolute;
    width: ${size}px;
    height: ${size}px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.6);
    left: ${x}px;
    top: ${y}px;
    pointer-events: none;
    transform: scale(0);
  `

  container.appendChild(ripple)

  gsap.to(ripple, {
    scale: 2,
    opacity: 0,
    duration: 0.6,
    ease: 'power2.out',
    onComplete: () => {
      ripple.remove()
    },
  })
}

/**
 * 页面切换动画 - 滑入
 * @param {Element} element - 要动画的元素
 * @param {String} direction - 方向 ('left' | 'right')
 */
export function pageSlideIn(element, direction = 'right') {
  const x = direction === 'right' ? '100%' : '-100%'

  return gsap.fromTo(
    element,
    {
      x: x,
      opacity: 0,
    },
    {
      x: 0,
      opacity: 1,
      duration: ANIMATION_CONFIG.duration.normal,
      ease: ANIMATION_CONFIG.ease.out,
      clearProps: 'all',
    }
  )
}

/**
 * 页面切换动画 - 滑出
 * @param {Element} element - 要动画的元素
 * @param {String} direction - 方向 ('left' | 'right')
 */
export function pageSlideOut(element, direction = 'left') {
  const x = direction === 'left' ? '-100%' : '100%'

  return gsap.to(element, {
    x: x,
    opacity: 0,
    duration: ANIMATION_CONFIG.duration.normal,
    ease: ANIMATION_CONFIG.ease.inOut,
  })
}

/**
 * 高亮闪烁动画
 * @param {Element} element - 要动画的元素
 * @param {String} color - 高亮颜色
 */
export function highlight(element, color = '#FF8C42') {
  const originalBg = window.getComputedStyle(element).backgroundColor

  return gsap.fromTo(
    element,
    {
      backgroundColor: color,
    },
    {
      backgroundColor: originalBg,
      duration: 1,
      ease: 'power2.out',
    }
  )
}

/**
 * 创建时间轴动画
 * @returns {gsap.core.Timeline}
 */
export function createTimeline(options = {}) {
  return gsap.timeline(options)
}

/**
 * 表单项交错动画
 * @param {Array|NodeList} formItems - 表单项元素
 */
export function staggerFormItems(formItems) {
  return gsap.to(
    formItems,
    {
      opacity: 1,
      x: 0,
      duration: ANIMATION_CONFIG.duration.normal,
      stagger: 0.05,
      delay: 0.1,
      ease: ANIMATION_CONFIG.ease.out,
      clearProps: 'transform',
    }
  )
}

/**
 * 骨架屏脉冲动画
 * @param {Element} element - 骨架屏元素
 */
export function skeletonPulse(element) {
  return gsap.to(element, {
    opacity: 0.5,
    duration: 1,
    repeat: -1,
    yoyo: true,
    ease: 'power1.inOut',
  })
}

export default {
  staggerFadeInUp,
  staggerScaleIn,
  popIn,
  popOut,
  slideUpIn,
  slideDownOut,
  fadeIn,
  fadeOut,
  removeAnimation,
  shake,
  pulse,
  ripple,
  pageSlideIn,
  pageSlideOut,
  highlight,
  createTimeline,
  staggerFormItems,
  skeletonPulse,
  ANIMATION_CONFIG,
}

