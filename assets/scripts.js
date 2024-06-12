import 'htmx.org'

import Alpine from 'alpinejs'
import * as htmx from "htmx.org";

// Add Alpine instance to window object.
window.Alpine = Alpine
window.htmx = htmx

// Start Alpine.
Alpine.start()

window.getSystemThemeName = function() {
    let t = '(prefers-color-scheme: dark)',
        m = window.matchMedia(t);
    if (m.media !== t || m.matches) {
        return 'dark';
    } else {
        return 'light';
    }
}

window.changeTheme = function() {
    try {
        let d = document.documentElement,
            c = d.classList;
        c.remove('light', 'dark');
        let e = localStorage.getItem('theme');
        if ('system' === e || (!e)) {
            if (getSystemThemeName() === 'dark') {
                d.style.colorScheme = 'dark';
                c.add('dark')
            } else {
                d.style.colorScheme = 'light';
                c.add('light')
            }
        } else if (e) {
            c.add(e || '')
        }
        if (e === 'light' || e === 'dark') d.style.colorScheme = e
    } catch (e) {}
}

window.metaData = {}

!function() {
    changeTheme();
}()