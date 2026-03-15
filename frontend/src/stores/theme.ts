import { onMounted, ref, watch } from 'vue'
import { useTheme } from 'vuetify'

const THEME_STORAGE_KEY = 'easyjet-theme'

export function useThemeStore () {
  const theme = useTheme()

  const storedTheme = localStorage.getItem(THEME_STORAGE_KEY)
  const defaultValue = storedTheme === 'dark' || (storedTheme === null && theme.global.name.value === 'dark')
  const isDark = ref(defaultValue)

  onMounted(() => {
    if (storedTheme !== null) {
      theme.global.name.value = storedTheme
      isDark.value = storedTheme === 'dark'
    }
  })

  watch(isDark, newVal => {
    const newTheme = newVal ? 'dark' : 'light'
    theme.global.name.value = newTheme
    localStorage.setItem(THEME_STORAGE_KEY, newTheme)
  })

  function toggle () {
    isDark.value = !isDark.value
  }

  return {
    isDark,
    toggle,
  }
}
