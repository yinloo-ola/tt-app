// uno.config.ts
import { defineConfig } from 'unocss'
import presetIcons from '@unocss/preset-icons'
import presetUno from '@unocss/preset-uno'
import presetTypography from '@unocss/preset-typography'
import presetWebFonts from '@unocss/preset-web-fonts'

export default defineConfig({
    presets: [
        presetIcons({}),
        presetUno(),
        presetTypography(),
        presetWebFonts({
            provider: 'google',
            fonts: {
                sans: 'Roboto',
            },
        }),
    ],
    shortcuts: {
        'tab-pill': 'px-4 py-2 flex items-center rounded-md cursor-pointer',
    },
})