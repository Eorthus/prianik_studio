import { defineNuxtConfig } from "nuxt/config";

export default defineNuxtConfig({
  ssr: true,

  i18n:{
    strategy: "prefix_except_default",
    defaultLocale: "es",
    // lazy: false,
    locales: [
      {
        code: "es",
        name: "Español",
        file: "es.json",
      },
      {
        code: "en",
        name: "English", 
        file: "en.json",
      },
      {
        code: "ru",
        name: "Русский",
        file: "ru.json",
      },
    ],
    // detectBrowserLanguage: {
    //   useCookie: true,
    //   cookieKey: "i18n_redirected",
    //   redirectOn: "root"
    // },
  },

  app: {
    head: {
      title: "Мастерская лазерной резки и 3D печати",
      htmlAttrs: { lang: "ru" },
      meta: [
        { charset: "utf-8" },
        { name: "viewport", content: "width=device-width, initial-scale=1" },
        {
          name: "description",
          content:
            "Мастерская лазерной резки и 3D печати. Изготовление индивидуальных сувениров, декора и аксессуаров по вашему дизайну.",
        },
        { name: "robots", content: "index, follow" },
        {
          property: "og:title",
          content: "Мастерская лазерной резки и 3D печати",
        },
        {
          property: "og:description",
          content:
            "Изготовление индивидуальных сувениров, декора и аксессуаров по вашему дизайну.",
        },
        { property: "og:type", content: "website" },
      ],
      link: [
        { rel: "icon", type: "image/x-icon", href: "/favicon.ico" },
        {
          rel: "stylesheet",
          href: "https://fonts.googleapis.com/css2?family=Roboto:wght@300;400;500;700&display=swap",
        },
      ],
    },
  },

  css: ["@/assets/scss/main.scss"],
  modules: ["@nuxtjs/tailwindcss", "@nuxtjs/i18n",     "@nuxtjs/sitemap"],
  plugins: [
    "~/plugins/gsap.ts"
  ],
  compatibilityDate: "2025-04-13",
  imports: {
    autoImport: true // ← это значение по умолчанию
  },

  nitro: {
    prerender: {
      //@ts-expect-error
      enabled: false
    }
  },

  runtimeConfig: {
    public: {
      apiBaseUrl: process.env.API_BASE_URL || 'http://localhost:8080/api',
    }
  },
});