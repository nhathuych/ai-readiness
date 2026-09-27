"use client";

import { useState } from "react";
import { Button, buttonVariants } from "../ui/button";
import { clientApiFetch } from "@/lib/api/client";
import { useRouter } from "next/navigation";
import { cn } from "@/lib/utils";

const logoutButtonClass = cn(
  buttonVariants({ variant: "default", size: "lg" }),
  "h-10 rounded-none px-6 text-base font-medium shadow-none",
);

function LogoutButton() {
  const [isPending, setIsPending] = useState(false);
  const router = useRouter();

  async function handleLogout() {
    try {
      setIsPending(true);

      await clientApiFetch<{ message: string }>("/auth/logout", {
        method: "POST",
      });

      router.push("/");
      router.refresh();
    } finally {
      setIsPending(false);
    }
  }

  return (
    <Button
      className={logoutButtonClass}
      variant="default"
      onClick={handleLogout}
      disabled={isPending}
    >
      Logout
    </Button>
  );
}

export default LogoutButton;
