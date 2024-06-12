/** @type {import('tailwindcss').Config} */
const colors = require('tailwindcss/colors')
module.exports = {
  content: ['**/*.{html,templ}'],
  darkMode: 'selector',
  theme: {
    extend: {},
    colors: {
      transparent: 'transparent',
      current: 'currentColor',
      black: colors.black,
      white: colors.white,
      gray: colors.slate,
      green: colors.emerald,
      purple: colors.violet,
      yellow: colors.amber,
      pink: colors.fuchsia,
      primary: colors.sky,
      secondary: colors.amber,
      neutral: colors.gray,
    },
  },
  plugins: [require('@tailwindcss/forms'), require('@tailwindcss/typography')],
}
