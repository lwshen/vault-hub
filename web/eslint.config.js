import js from '@eslint/js';
import globals from 'globals';
import reactHooks from 'eslint-plugin-react-hooks';
import reactRefresh from 'eslint-plugin-react-refresh';
import eslintReact from "@eslint-react/eslint-plugin";
import tseslint from 'typescript-eslint';
import stylistic from '@stylistic/eslint-plugin';

export default tseslint.config(
  { ignores: ['dist'] },
  {
    extends: [
      js.configs.recommended,
      ...tseslint.configs.recommended,
      eslintReact.configs["recommended-typescript"],
    ],
    files: ["**/*.{ts,tsx}"],
    ignores: ['src/components/ui/**/*'],
    languageOptions: {
      ecmaVersion: 2020,
      globals: globals.browser,
    },
    plugins: {
      'react-hooks': reactHooks,
      'react-refresh': reactRefresh,
      '@stylistic': stylistic,
    },
    rules: {
      ...reactHooks.configs.recommended.rules,
      'react-refresh/only-export-components': [
        'warn',
        { allowConstantExport: true },
      ],
      '@stylistic/indent': ['error', 2],
      "@stylistic/indent-binary-ops": ["error", 2],
      '@stylistic/quotes': ['error', 'single', { avoidEscape: true }],
      '@stylistic/brace-style': ['error', '1tbs', { allowSingleLine: false }],
      "@stylistic/semi": "error",
      "@stylistic/eol-last": "error",
      '@stylistic/member-delimiter-style': [
        'error',
        {
          multiline: {
            delimiter: 'semi',
            requireLast: true,
          },
          singleline: {
            delimiter: 'semi',
            requireLast: true,
          },
        },
      ],
      '@stylistic/no-trailing-spaces': 'error',
      '@stylistic/no-multiple-empty-lines': 'error',
      '@stylistic/object-curly-spacing': ['error', 'always'],
      '@stylistic/comma-spacing': ['error', { before: false, after: true }],
      '@stylistic/keyword-spacing': ['error', { before: true, after: true }],
      '@stylistic/space-before-blocks': 'error',
      '@stylistic/space-infix-ops': 'error',
      '@stylistic/arrow-spacing': ['error', { before: true, after: true }],
      '@stylistic/comma-dangle': ['error', 'always-multiline'],
      '@stylistic/array-bracket-spacing': ['error', 'never'],
      '@stylistic/space-in-parens': ['error', 'never'],
      '@stylistic/spaced-comment': ['error', 'always', { exceptions: ['-', '=', '*'] }],
    },
  },
)
