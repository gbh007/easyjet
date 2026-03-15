/**
 * plugins/vuetify.ts
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Composables
import { createVuetify } from 'vuetify'
// Styles
import '@mdi/font/css/materialdesignicons.css'

import 'vuetify/styles'

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  theme: {
    defaultTheme: 'dark',
    themes: {
      light: {
        dark: false,
        colors: {
          'background': '#ffffff',
          'surface': '#f5f5f5',
          'primary': '#000000',
          'secondary': '#424242',
          'accent': '#82b1ff',
          'error': '#ff5252',
          'info': '#2196f3',
          'success': '#4caf50',
          'warning': '#ff9800',
          'on-primary': '#ffffff',
          'on-surface': '#000000',
        },
      },
      dark: {
        dark: true,
        colors: {
          'background': '#121212',
          'surface': '#1e1e1e',
          'primary': '#ffffff',
          'secondary': '#b0b0b0',
          'accent': '#82b1ff',
          'error': '#ff5252',
          'info': '#2196f3',
          'success': '#4caf50',
          'warning': '#ff9800',
          'on-primary': '#000000',
          'on-surface': '#ffffff',
        },
      },
    },
  },
})
