/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./components/*.templ"],
  theme: {
    screens: {
      scTwo: {'max':'700px'},
      scTextOne: {'max':'620px'},
      scTextTwo: {'max':'555px'},
      scOne: {'max':'550px'},
    },
    extend: {
      spacing: {
        '99%': '99%',
        'chart-highlight': 'calc(100vh - 37rem)',
      }
    },
  },
  plugins: [],
}

