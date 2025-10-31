/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          50: '#eef7ff',
          100: '#d9eeff',
          200: '#bce3ff',
          300: '#8ed4ff',
          400: '#5ac0ff',
          500: '#31a5ff',
          600: '#007aff',
          700: '#0062cc',
          800: '#0052a8',
          900: '#004082',
          950: '#00264d',
        },
        neutral: {
          50: '#f8f9fa',
          100: '#e9ecef',
          200: '#dee2e6',
          300: '#ced4da',
          400: '#adb5bd',
          500: '#6c757d',
          600: '#495057',
          700: '#343a40',
          800: '#212529',
          900: '#1c1c1e',
          950: '#0d0d0d',
        },
        success: '#34c759',
        warning: '#ff9500',
        danger: '#ff3b30',
      },
    },
  },
  plugins: [],
}

