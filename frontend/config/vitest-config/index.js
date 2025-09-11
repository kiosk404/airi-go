module.exports = {
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/setup.ts'],
    css: true,
  },
  coverage: {
    provider: 'v8',
    reporter: ['text', 'json', 'html'],
    exclude: [
      'node_modules/',
      'src/setup.ts',
      '**/*.d.ts',
      '**/*.config.*',
    ],
  },
};






