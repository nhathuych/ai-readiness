import SiteNavbar from "@/components/shared/site-navbar";
import { buttonVariants } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { getUserInfo } from "@/lib/api/auth";
import { publicApiBaseUrl } from "@/lib/config";
import { cn } from "@/lib/utils";
import Link from "next/link";

const workflowSteps = [
  "Submit onboarding details",
  "Get an AI readiness review",
  "Wait for admin decision",
];

const benefits = [
  {
    title: "Clear intake",
    description:
      "Collect the service, project goal, timeline, assets, and notes in one structured onboarding flow.",
  },
  {
    title: "Faster review",
    description:
      "The team gets a readiness summary, missing items, next action, and suggested follow-up before making a decision.",
  },
  {
    title: "Simple approval flow",
    description:
      "Clients can see whether the agency approved or rejected their onboarding request without any extra checkout step.",
  },
];

const mainClass =
  "mx-auto grid max-w-7xl gap-10 px-6 py-14 lg:grid-cols-[1.15fr_0.85fr] lg:items-start";

const heroSectionClass = "grid gap-8";

const heroContentClass = "grid gap-5";

const eyebrowClass =
  "text-sm font-semibold uppercase tracking-[0.2em] text-muted-foreground";

const headingClass =
  "max-w-4xl text-5xl font-semibold leading-tight tracking-tight sm:text-6xl";

const descriptionClass = "max-w-3xl text-lg leading-8 text-muted-foreground";

const actionsClass = "flex flex-wrap gap-3";

const primaryButtonClass = cn(
  buttonVariants({ variant: "default", size: "lg" }),
  "h-12 rounded-none px-6 text-base font-medium shadow-none",
);

const secondaryButtonClass = cn(
  buttonVariants({ variant: "outline", size: "lg" }),
  "h-12 rounded-none border-border px-6 text-base font-medium shadow-none",
);

const benefitsGridClass = "grid gap-4 sm:grid-cols-3";

const benefitTitleClass = "text-lg";

const benefitDescriptionClass = "text-base leading-7 text-muted-foreground";

const workflowContentClass = "grid gap-5";

const workflowItemClass = "flex gap-4";

const workflowNumberClass =
  "flex h-9 w-9 shrink-0 items-center justify-center border border-border bg-card text-sm font-semibold";

const workflowTextClass = "grid gap-1";

const workflowTitleClass = "font-medium text-foreground";

const workflowDescriptionClass = "text-sm leading-6 text-muted-foreground";

export default async function Home() {
  const user = await getUserInfo();
  const workspaceHref = user?.role === "admin" ? "/admin" : "/submissions";
  const loginHref = `${publicApiBaseUrl}/auth/google`;
  const viewProcessHref = user ? workspaceHref : "/";

  return (
    <>
      {/* navbar */}
      <SiteNavbar user={user} />
      <main className={mainClass}>
        <section className={heroSectionClass}>
          <div className={heroContentClass}>
            <p className={eyebrowClass}>
              Client Onboarding portal for agencies
            </p>
            <h1 className={headingClass}>
              Start every client project with the right details in place.
            </h1>
            <p className={descriptionClass}></p>
          </div>

          <div className={actionsClass}>
            {user ? (
              <Link href={workspaceHref} className={primaryButtonClass}>
                Open Workspace
              </Link>
            ) : (
              <a href={loginHref} className={primaryButtonClass}>
                Start Onboarding
              </a>
            )}

            <Link className={secondaryButtonClass} href={viewProcessHref}>
              View Process
            </Link>
          </div>

          <div className={benefitsGridClass}>
            {benefits.map((benefit) => (
              <Card key={benefit.title}>
                <CardHeader>
                  <CardTitle className={benefitTitleClass}>
                    {benefit.title}
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <p className={benefitDescriptionClass}>
                    {benefit.description}
                  </p>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>

        <Card>
          <CardHeader>
            <CardTitle>How Onboarding Works</CardTitle>
          </CardHeader>
          <CardContent className={workflowContentClass}>
            {workflowSteps.map((step, index) => (
              <div key={step} className={workflowItemClass}>
                <span className={workflowNumberClass}>{index + 1}</span>

                <div className={workflowTextClass}>
                  <p className={workflowTitleClass}>{step}</p>
                </div>
              </div>
            ))}
          </CardContent>
        </Card>
      </main>
    </>
  );
}
