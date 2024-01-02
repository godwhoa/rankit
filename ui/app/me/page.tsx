import {
  CardTitle,
  CardDescription,
  CardHeader,
  CardContent,
  Card,
} from "@/components/ui/card";
import { headers } from "next/headers";

const API_BASE_URL = "http://localhost:8000/v1";

export default async function Me() {
  const response = await fetch(`${API_BASE_URL}/users/me`, {
    method: "GET",
    headers: {
      cookie: headers().get("cookie") || "",
      "Content-Type": "application/json",
    },
  });
  const data = await response.json();
  return (
    <main className="min-h-screen bg-gray-100 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <Card className="mx-auto max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-3xl font-bold text-center">Me</CardTitle>
          <CardDescription className="text-center">You details</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <code>{JSON.stringify(data)}</code>
        </CardContent>
      </Card>
    </main>
  );
}
