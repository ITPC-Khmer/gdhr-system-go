/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // DarkOrange palette (primary). 600 ≈ #FF8C00 (darkorange).
        primary: {
          50: '#fff7ed',
          100: '#ffedd5',
          200: '#fed7aa',
          300: '#fdba74',
          400: '#fb923c',
          500: '#ff8c00',
          600: '#f97316',
          700: '#c2560c',
          800: '#9a3d0c',
          900: '#7c3410',
          950: '#431807',
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'Avenir', 'Helvetica', 'Arial', 'sans-serif'],
      },
      boxShadow: {
        soft: '0 2px 8px -1px rgba(0,0,0,0.06), 0 4px 16px -2px rgba(0,0,0,0.05)',
      },
    },
  },
  plugins: [],
}
