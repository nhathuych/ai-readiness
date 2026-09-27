import { cn } from "@/lib/utils";
import { buttonVariants } from "../ui/button";
import { User } from "@/types/auth";
import Link from "next/link";
import { publicApiBaseUrl } from "@/lib/config";
import LogoutButton from "./logout-button";

const brandLinkClass = "flex items-center gap-3";

const brandLogoClass =
  "flex h-11 w-11 items-center justify-center border border-border bg-card text-sm font-semibold text-foreground";

const brandTextWrapperClass = "grid gap-0.5";

const brandTitleClass = "text-base font-semibold tracking-tight";

const headerClass = "sticky top-0 z-40 border-b bg-background/95 backdrop-blur";

const navbarContainerClass =
  "mx-auto flex max-w-7xl items-center justify-between gap-4 px-6 py-4";

const navActionsClass = "flex items-center gap-3";

const workspaceButtonClass = cn(
  buttonVariants({ variant: "ghost", size: "sm" }),
  "h-10 rounded-none px-4 text-sm font-medium shadow-none",
);

const logoutButtonClass = cn(
  buttonVariants({ variant: "outline", size: "sm" }),
  "h-10 rounded-none border-border px-5 text-sm font-medium shadow-none",
);

const loginButtonClass = cn(
  buttonVariants({ variant: "outline", size: "sm" }),
  "h-10 rounded-none border-border bg-primary px-5 text-sm font-medium text-white shadow-none",
);

function SiteNavbar({ user }: { user: User | null }) {
  const workspaceHref = user?.role === "admin" ? "/admin" : "/submissions";

  return (
    <header className={headerClass}>
      <div className={navbarContainerClass}>
        <Link className={brandLinkClass} href="/">
          <span className={brandLogoClass}>AC</span>
          <div className={brandTextWrapperClass}>
            <span className={brandTitleClass}>Agency Client Portal</span>
          </div>
        </Link>

        {user ? (
          <div className={navActionsClass}>
            <Link href={workspaceHref} className={workspaceButtonClass}>
              Workspace
            </Link>
            <LogoutButton />
          </div>
        ) : (
          <Link
            href={`${publicApiBaseUrl}/auth/google`}
            className={loginButtonClass}
          >
            Login
          </Link>
        )}
      </div>
    </header>
  );
}

export default SiteNavbar;
