import type { NextConfig } from "next";

  

const nextConfig: NextConfig = {

  /* config options here */

  rewrites: async () => [
    {
      source: "/api/:path*",
      destination: "http://localhost:8080/api/:path*",
    },
  ],

  // 使用 standalone 模式替代静态导出

  output: "standalone",
  typescript: {
    ignoreBuildErrors: true
  },
  eslint: {
    ignoreDuringBuilds: true
  }
};

  

export default nextConfig;