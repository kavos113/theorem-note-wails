import { defineConfig } from 'eslint/config';

export default defineConfig({
  ignores: ['node_modules', 'dist', 'build'],
  extends: ['plugin:vue/vue3-recommended', 'eslint:recommended', '@vue/typescript/recommended'],
  parserOptions: {
    ecmaVersion: 2020,
    sourceType: 'module'
  },
  rules: {
    'vue/multi-word-component-names': 'off',
    'vue/block-lang': [
      'error',
      {
        script: {
          lang: 'ts'
        }
      }
    ],
    'vue/no-v-html': 'off',
    'vue/block-order': [
      'error',
      {
        order: ['script', 'template', 'style']
      }
    ]
  }
});
