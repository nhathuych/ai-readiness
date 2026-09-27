"use client";

import { useState } from "react";
import CreateSubmissionClient from "./create-submission-dialog";
import { CreateSubmissionResponse, Submission } from "@/types/submission";
import { Card, CardContent } from "../ui/card";
import SubmissionCard from "./submission-card";
import SubmissionDetails from "./submission-details-dialog";

const mainClass = "mx-auto grid max-w-7xl gap-6 px-6 py-8 lg:py-10";

const headerSectionClass =
  "flex flex-col gap-4 md:flex-row md:items-end md:justify-between";

const headerTextWrapperClass = "grid gap-2";

const pageTitleClass = "text-3xl font-semibold tracking-tight";

const pageDescriptionClass =
  "max-w-2xl text-base leading-7 text-muted-foreground";

const submissionsGridClass = "grid gap-4 md:grid-cols-2 xl:grid-cols-3";

const emptyCardClass = "rounded-none border-border shadow-none";

const emptyCardContentClass = "px-6 py-10";

const emptyTitleClass = "text-lg font-medium tracking-tight";

const emptyDescriptionClass = "mt-2 text-base leading-7 text-muted-foreground";

function SubmissionPageClient({
  initialSubmissions,
}: {
  initialSubmissions: Submission[];
}) {
  const [submissions, setSubmissions] = useState(initialSubmissions);
  const [selectedSubmission, setSelectedSubmission] =
    useState<Submission | null>(null);

  function handleCreated(payload: CreateSubmissionResponse) {
    setSubmissions((current) => [payload.submission, ...current]);
  }

  return (
    <main className={mainClass}>
      <section className={headerSectionClass}>
        <div className={headerTextWrapperClass}>
          <h1 className={pageTitleClass}>My onboarding submissions</h1>
          <p className={pageDescriptionClass}>
            Submit your onboarding details, review the AI readiness result, and
            track the final admin decision from one page.
          </p>
        </div>
        <CreateSubmissionClient onCreated={handleCreated} />
      </section>

      {submissions.length ? (
        <section className={submissionsGridClass}>
          {submissions.map((item) => (
            <SubmissionCard
              onView={setSelectedSubmission}
              key={item.id}
              submission={item}
            />
          ))}
        </section>
      ) : (
        <Card className={emptyCardClass}>
          <CardContent className={emptyCardContentClass}>
            <p className={emptyTitleClass}>No onboarding requests</p>
          </CardContent>
        </Card>
      )}

      <SubmissionDetails
        submission={selectedSubmission}
        onClose={() => setSelectedSubmission(null)}
      />
    </main>
  );
}

export default SubmissionPageClient;
