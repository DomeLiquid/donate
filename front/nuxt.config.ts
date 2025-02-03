import Tailwind from "@tailwindcss/vite";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  ssr: false,
  devtools: {
    enabled: true,
  },
  compatibilityDate: "2024-11-23",
  modules: [
    "@nuxt/eslint",
    "radix-vue/nuxt",
    "@nuxt/icon",
    "@nuxt/image",
    "@nuxt/fonts",
    "nuxt-security",
    "@vueuse/nuxt",
    "@pinia/nuxt",
  ],
  css: ["~/assets/css/main.css"],
  vite: {
    plugins: [Tailwind()],
  },
  fonts: {
    experimental: {
      processCSSVariables: true,
    },
  },
  postcss: {
    plugins: {
      cssnano: {
        preset: "default",
      },
    },
  },
  image: {
    quality: 80,
    format: ["webp"],
  },
  app: {
    head: {
      title: "Welcome",
      titleTemplate: "%s - Nuxt 4 Starter Template",
      link: [{ rel: "icon", type: "image/png", href: "/favicon.png" }],
    },
  },
  runtimeConfig: {
    public: {
      buildAt: new Date().toLocaleString("zh-CN", {
        timeZone: "Asia/Shanghai",
      }),
      environment: "production",
    },
  },
  future: {
    compatibilityVersion: 4,
  },
  typescript: {
    strict: true,
  },
  telemetry: false,
});
