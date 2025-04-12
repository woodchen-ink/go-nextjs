import type { Metadata } from "next";
import { Inter as FontSans } from "next/font/google";
import { cn } from "@/lib/utils";
import { Toaster } from "@/components/ui/sonner";

import "./globals.css";

const fontSans = FontSans({
  subsets: ["latin"],
  variable: "--font-sans",
});


// 动态生成元数据
export async function generateMetadata(): Promise<Metadata> {
  // 在构建时使用默认值，避免动态fetches
  if (process.env.NODE_ENV === 'production') {
    return {
      title: "VPS促销监控",
      description: "监控并展示各家VPS促销信息",
    };
  }
  return {
    title: "VPS促销监控",
    description: "监控并展示各家VPS促销信息",
  };
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN" suppressHydrationWarning className="overflow-y-scroll">
      <head>
        <link rel="icon" href="/logo.svg" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0, viewport-fit=cover" />
      </head>
      <body
        className={cn(
          "min-h-screen bg-background font-sans antialiased overflow-y-scroll overflow-x-hidden",
          fontSans.variable
        )}
      >
        <div className="relative flex min-h-screen flex-col">
          <div className="w-full max-w-screen-2xl mx-auto px-4">
            <div className="flex-1">{children}</div>
          </div>
        </div>
        <Toaster />
      </body>
    </html>
  );
}
