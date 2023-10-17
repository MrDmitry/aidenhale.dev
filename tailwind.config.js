/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ['./web/templates/*.html'],
    theme: {
        extend: {
            colors : {
                'material' : {
                    'primary' : '#1f2335',
                    'secondary' : '#c0caf5',
                    'highlight' : '#24283b',
                },
            },
        },
    },
    plugins: [
        require('@tailwindcss/typography'),
    ],
    darkMode: 'class',
}
