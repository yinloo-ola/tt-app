// uno.config.ts
import { defineConfig } from 'unocss'
import presetIcons from '@unocss/preset-icons'
import presetUno from '@unocss/preset-uno'

export default defineConfig({
    presets: [
        presetIcons({}),
        presetUno(),
    ],
    shortcuts: {
        'tab-pill': 'px-4 py-2 flex items-center rounded-md cursor-pointer',
        'tab-selected': 'bg-yellow-400'
    },
})