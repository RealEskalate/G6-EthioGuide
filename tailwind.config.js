/** @type {import('tailwindcss').Config} */
const config = {
  content: [
    './app/**/*.{js,ts,jsx,tsx}',
    './components/**/*.{js,ts,jsx,tsx}',
    './pages/**/*.{js,ts,jsx,tsx}',
    './src/**/*.{js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#3A6A8D',
          light: '#5B8DB1',
          dark: '#2A4B66',
        },
        secondary: {
          DEFAULT: '#5E9C8D',
          light: '#80B7A9',
          dark: '#44776A',
        },
        accent: {
          DEFAULT: 'rgba(58,106,141,0.1)',
          light: 'rgba(58,106,141,0.1)',
        },
        neutral: {
          DEFAULT: '#A7B3B9',
          light: '#C8D1D6',
          dark: '#85939A',
        },
        background: {
          light: '#F5F7F8',
          dark: '#2E4D57',
        },
        error: '#DC2626',
      },
      fontFamily: {
        sans: ['Inter', 'sans-serif'],
        amharic: ['Noto Sans Ethiopic', 'sans-serif'],
      },
    },
  },
  plugins: [
    
  ],
  darkMode: 'class',
};

module.exports = config;
