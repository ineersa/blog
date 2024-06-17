import 'htmx.org'

import Alpine from 'alpinejs'
import * as htmx from 'htmx.org'
import * as hljs from 'highlight.js/lib/common'

// Add Alpine instance to window object.
window.Alpine = Alpine
window.htmx = htmx
window.hljs = hljs

hljs.highlightAll()

window.getSystemThemeName = function () {
  let t = '(prefers-color-scheme: dark)',
    m = window.matchMedia(t)
  if (m.media !== t || m.matches) {
    return 'dark'
  } else {
    return 'light'
  }
}

window.changeTheme = function () {
  try {
    let d = document.documentElement,
      c = d.classList
    c.remove('light', 'dark')
    let e = localStorage.getItem('theme')
    if ('system' === e || !e) {
      if (getSystemThemeName() === 'dark') {
        d.style.colorScheme = 'dark'
        c.add('dark')
        localStorage.setItem('theme', 'dark')
      } else {
        d.style.colorScheme = 'light'
        c.add('light')
        localStorage.setItem('theme', 'dark')
      }
    } else if (e) {
      c.add(e || 'dark')
      localStorage.setItem('theme', e || 'dark')
    }
    if (e === 'light' || e === 'dark') d.style.colorScheme = e
  } catch (e) {
    console.error(e)
  }
  return localStorage.getItem('theme')
}

window.metaData = {}

const domReady = (callback) => {
  document.addEventListener('DOMContentLoaded', callback)
}

domReady(() => {
  // Display body when DOM is loaded
  document.body.style.visibility = 'visible'
  changeTheme()
  // Start Alpine.
  Alpine.start()
})
