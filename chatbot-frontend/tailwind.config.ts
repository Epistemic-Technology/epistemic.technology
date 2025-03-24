import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/**/*.{js,ts,jsx,tsx,mdx}",
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        // You can add custom colors here
      },
      fontFamily: {
        // You can add custom fonts here
      },
      spacing: {
        // Custom spacing can be added here
      },
      screens: {
        // Custom breakpoints can be defined here
      },
    },
  },
  plugins: [],
};

export default config;
