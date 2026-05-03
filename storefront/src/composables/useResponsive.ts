import { ref, onMounted, onUnmounted } from 'vue'

export function useResponsive() {
    const isMobile = ref(false)

    const checkResponsive = () => {
        isMobile.value = window.innerWidth < 768
    }

    onMounted(() => {
        checkResponsive()
        window.addEventListener('resize', checkResponsive)
    })

    onUnmounted(() => {
        window.removeEventListener('resize', checkResponsive)
    })

    const getProp = (val: any, defaultValue: any = '') => {
        if (val === undefined || val === null) return defaultValue

        if (typeof val === 'object' && val !== null) {
            if (isMobile.value && 'mobile' in val) {
                return val.mobile ?? val.desktop ?? defaultValue
            }
            return val.desktop ?? defaultValue
        }

        return val ?? defaultValue
    }

    return { isMobile, getProp }
}
