import { defineConfig } from 'orval';

export default defineConfig({
  'easyjet-api': {
    input: '../openapi.yaml',
    output: {
      target: 'src/api/generated.ts',
      client: 'axios',
      mode: 'split',
      clean: true,
      override: {
        axios: {
          axiosInstance: 'src/api/axios',
        },
      },
      fileHeader: '/* eslint-disable */',
    },
  },
});
