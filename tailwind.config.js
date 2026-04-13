/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['"DM Sans"', 'sans-serif'],
        mono: ['"DM Mono"', 'monospace'],
      },
      colors: {
        bg: '#f2f2f7',
        surface: '#ffffff',
        border: '#e0e0e5',
        'border-strong': '#c8c8d0',
        primary: '#1c1c1e',
        secondary: '#6c6c70',
        muted: '#aeaeb2',
        'green-main': '#34c759',
        'green-bg': '#e8f8ed',
        'yellow-main': '#ff9f0a',
        'yellow-bg': '#fff5e6',
        'red-main': '#ff3b30',
        'red-bg': '#fff0ef',
      },
    },
  },
  plugins: [],
}
