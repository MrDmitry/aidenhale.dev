/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ['./web/templates/*.html'],
    theme: {
        extend: {
            colors: {
                'tokyo': {
                    'dark-one': '#1B1E2C',
                    'dark-two': '#24283B',
                    'accent-one': '#1ABC9C',
                    'accent-two': '#DB4B4B',
                    'accent-three': '#E0AF68',
                    'accent-four': '#51A0CF',
                    'accent-five': '#9CCB69',
                    'accent-six': '#9375CA',
                },
            },
            typography: ({ theme }) => ({
                tokyo: {
                    css: {
                        '--tw-prose-body': 'rgb(255 255 255 / 90%)',
                        '--tw-prose-headings': theme("colors")['tokyo']['accent-five'],
                        '--tw-prose-lead': 'rgb(255 255 255 / 90%)',
                        '--tw-prose-links': 'rgb(255 255 255 / 90%)',
                        '--tw-prose-bold': 'rgb(255 255 255 / 90%)',
                        '--tw-prose-counters': 'rgb(255 255 255 / 90%)',
                        '--tw-prose-bullets': 'rgb(255 255 255 / 90%)',
                        '--tw-prose-hr': theme("colors")['tokyo']['accent-five'],
                        '--tw-prose-quotes': 'rgb(255 255 255 / 90%)',
                        '--tw-prose-quote-borders': 'rgb(255 255 255 / 90%)',
                        '--tw-prose-captions': 'rgb(255 255 255 / 90%)',
                        '--tw-prose-kbd': 'rgb(255 255 255 / 90%)',
                        '--tw-prose-kbd-shadows': 'rgb(0 0 0 / 50%)',
                        '--tw-prose-code': 'rgb(255 255 255 / 90%)', 
                        '--tw-prose-pre-code': theme("colors")['tokyo']['accent-one'],
                        '--tw-prose-pre-bg': theme("colors")['tokyo']['dark-two'],
                        '--tw-prose-th-borders': 'rgb(255 255 255 / 90%)',
                        '--tw-prose-td-borders': 'rgb(255 255 255 / 90%)',
                    },
                },
            }),
        },
    },
    plugins: [
        require('@tailwindcss/typography'),
    ],
    darkMode: 'class',
}
