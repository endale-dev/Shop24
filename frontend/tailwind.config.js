/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,ts}'],
  theme: {
    extend: {
      colors: {
        bg: '#f4f5f7',
        surface: '#ffffff',
        border: '#d0d5dd',
        text: '#111418',
        muted: '#5c6570',
        accent: '#1e40af',
        'accent-hover': '#1e3a8a',
        danger: '#9f1c1c',
        'danger-hover': '#7f1616',
        'good-bg': '#e5f4ea',
        'good-fg': '#146c36',
        'warn-bg': '#fbf3d8',
        'warn-fg': '#8a5a00',
        'bad-bg': '#fde8e8',
        'bad-fg': '#9f1c1c',
      },
      fontFamily: {
        sans: ['Figtree', 'ui-sans-serif', 'system-ui', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
