export default {
    prefix: 'tw-',
    content: [
      "./components/**/*.{js,vue,ts}",
      "./layouts/**/*.vue",
      "./pages/**/*.vue",
      "./plugins/**/*.{js,ts}",
      "./nuxt.config.{js,ts}",
      "./app.vue",
    ],
    theme: {
      extend: {
        colors: {
          'primary': '#222222',
          'secondary': '#444444',
          'light-gray': '#888888',
          'white': '#ffffff',
        },
        fontFamily: {
          'sans': ['Roboto', 'sans-serif'],
        },
        container: {
          center: true,
          padding: {
            DEFAULT: '1rem',
            sm: '1.5rem',
            lg: '2rem',
          },
          screens: {
            sm: '640px',
            md: '768px',
            lg: '1024px',
            xl: '1280px',
          },
        },
        transitionDuration: {
          DEFAULT: '300ms',
        },
        transitionTimingFunction: {
          DEFAULT: 'cubic-bezier(0.4, 0, 0.2, 1)',
        },
        animation: {
          'fade-in': 'fadeIn 0.7s ease-in-out forwards',
          'slide-up': 'slideUp 0.5s ease-out forwards',
        },
        keyframes: {
          fadeIn: {
            '0%': { opacity: '0' },
            '100%': { opacity: '1' },
          },
          slideUp: {
            '0%': { transform: 'translateY(20px)', opacity: '0' },
            '100%': { transform: 'translateY(0)', opacity: '1' },
          },
        },
      },
    },
    plugins: [],
  }