import Link from "next/link";
import { AlertCircle } from "lucide-react";
import { Button } from "@/components/ui/button";

export default function NotFound() {
  return (
    <div className="container flex flex-col items-center justify-center gap-6 py-24">
      <div className="flex items-center justify-center rounded-full bg-muted p-6">
        <AlertCircle className="h-12 w-12 text-muted-foreground" />
      </div>
      <h1 className="text-4xl font-bold">页面不存在</h1>
      <p className="text-lg text-muted-foreground">
        您访问的页面不存在或已被删除
      </p>
      <Link href="/">
        <Button>返回首页</Button>
      </Link>
    </div>
  );
} 