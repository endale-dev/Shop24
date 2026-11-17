/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,ts}'],
  theme: {
    extend: {
      colors: {
        bg: '#e6e3de',
        surface: '#f4f2ee',
        border: '#d0cbc3',
        text: '#1b1714',
        muted: '#5c564e',
        accent: '#7c3f1d',
        'accent-hover': '#5e2e14',
        danger: '#8e1e1e',
        'danger-hover': '#721818',
        'good-bg': '#dcecde',
        'good-fg': '#145c32',
        'warn-bg': '#f6e6c8',
        'warn-fg': '#7a3e10',
        'bad-bg': '#f8e2dc',
        'bad-fg': '#8e1e1e',
      },
      fontFamily: {
        sans: ['Figtree', 'ui-sans-serif', 'system-ui', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
