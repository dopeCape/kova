import { defineConfig } from 'tsup';

export default defineConfig({
  entry: {
    index: 'src/index.ts',
    server: 'src/server.ts',
    client: 'src/client.ts',
  },
  format: ['cjs', 'esm'],
  dts: true,
  sourcemap: true,
  external: ['react', 'react-dom'],
  clean: true,
  splitting: false,
  minify: false,
});
