import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  //docker container
  server: {
    host: '0.0.0.0',
    port: 5173,
  },
  // server: {
  //   proxy: {
  //     '/api': {
  //       target: 'http://localhost:5000', // Docker service name
  //       changeOrigin: true,
  //       rewrite: (path) => path.replace(/^\/api/, '')
  //     }
  //   }
  // },  

  //local server
  // server: {
  //   origin: 'http://localhost:5173',
  //   port: 5173,
  // },  
})
