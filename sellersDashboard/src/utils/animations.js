/**
 * Animation Utilities
 * Reusable animation helpers and Vue composable
 */

/**
 * Stagger animation delays for list items
 * @param {number} index - Item index
 * @param {number} delay - Base delay in ms
 * @returns {Object} Style object with animation delay
 */
export function staggerDelay(index, delay = 50) {
  return {
    animationDelay: `${index * delay}ms`
  }
}

/**
 * Vue composable for scroll animations
 * @returns {Object} Intersection observer functions
 */
export function useScrollAnimation() {
  const observeElement = (element, callback) => {
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            callback(entry.target)
            observer.unobserve(entry.target)
          }
        })
      },
      {
        threshold: 0.1,
        rootMargin: '0px 0px -100px 0px'
      }
    )

    observer.observe(element)
    return observer
  }

  return {
    observeElement
  }
}

/**
 * Page transition animations
 */
export const pageTransitions = {
  fade: {
    enterActiveClass: 'animate-fade-in',
    leaveActiveClass: 'animate-fade-in'
  },
  slide: {
    enterActiveClass: 'animate-slide-up',
    leaveActiveClass: 'animate-fade-in'
  },
  scale: {
    enterActiveClass: 'animate-scale-in',
    leaveActiveClass: 'animate-fade-in'
  }
}
