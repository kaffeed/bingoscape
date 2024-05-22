/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    'views/**/*.templ',
  ],
  darkMode: 'class',
  theme: {
    extend: {
      fontFamily: {
        mono: ['runescape'],
      }
    },
  },
  daisyui: {
    themes: [
      "dark", "light"
    ],
  },
  plugins: [
    require('daisyui'),
    require('@tailwindcss/forms'),
  ],
  corePlugins: {
    preflight: true,
  }
}
