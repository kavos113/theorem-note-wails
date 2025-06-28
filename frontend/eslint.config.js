import js from '@eslint/js';
import globals from 'globals';
import tseslint from 'typescript-eslint';
import pluginVue from 'eslint-plugin-vue';
import { defineConfig } from 'eslint/config';

export default defineConfig([
  { ignores: ['**/node_modules/**', '**/dist/**', '**/wailsjs/**'] },
  { files: ['**/*.{js,mjs,cjs,ts,mts,cts,vue}'], plugins: { js }, extends: ['js/recommended'] },
  { files: ['**/*.{js,mjs,cjs,ts,mts,cts,vue}'], languageOptions: { globals: globals.browser } },
  tseslint.configs.recommended,
  pluginVue.configs['flat/essential'],
  {
    files: ['**/*.vue'],
    languageOptions: { parserOptions: { parser: tseslint.parser } },
    rules: {
      'vue/no-v-html': 'off',
      'vue/block-lang': [
        'error',
        {
          script: { lang: 'ts' }
        }
      ],
      'vue/block-order': [
        'error',
        {
          order: ['script', 'template', 'style']
        }
      ]
    }
  }
]);
