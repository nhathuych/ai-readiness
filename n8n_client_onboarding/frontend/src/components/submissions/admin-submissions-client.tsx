"use client";

import {
  AdminStatus,
  Submission,
  UpdateSubmissionStatusResponse,
} from "@/types/submission";
import { useState } from "react";
import { Card, CardContent } from "../ui/card";
import { ReadinessBadge } from "../shared/readiness-badge";
import AdminStatusBadge from "../shared/admin-status-badge";
import { formatDate } from "@/lib/utils";
import { Button } from "../ui/button";
import SubmissionDetails from "./submission-details-dialog";
import { clientApiFetch } from "@/lib/api/client";

const summaryCardClass = "rounded-none border-border shadow-none";

const summaryCardContentClass = "grid gap-1 p-5";

const summaryLabelClass = "text-sm text-muted-foreground";

const summaryValueClass = "text-3xl font-semibold tracking-tight";

const mainClass = "mx-auto grid max-w-7xl gap-6 px-6 py-8 lg:py-10";

const pageHeaderClass = "grid gap-2";

const pageTitleClass = "text-3xl font-semibold tracking-tight";

const pageDescriptionClass =
  "max-w-3xl text-base leading-7 text-muted-foreground";

const statsGridClass = "grid gap-4 md:grid-cols-3";

const errorMessageClass =
  "border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700";

const tableCardClass = "rounded-none border-border shadow-none";

const tableCardContentClass = "p-0";

const tableScrollClass = "overflow-x-auto";

const tableClass = "w-full min-w-[900px] border-collapse text-sm";

const tableHeadClass =
  "border-b border-border bg-muted/50 text-left text-xs uppercase tracking-[0.16em] text-muted-foreground";

const tableHeaderCellClass = "px-5 py-4 font-medium";

const tableRowClass = "border-b border-border last:border-b-0";

const tableCellClass = "px-5 py-4 align-top";

const mutedTableCellClass = "px-5 py-4 align-top text-muted-foreground";

const clientInfoClass = "grid gap-1";

const clientNameClass = "font-medium text-foreground";

const clientEmailClass = "text-muted-foreground";

const tableActionsClass = "flex flex-wrap gap-2";

const tableButtonClass = "h-9 rounded-none px-3 text-xs shadow-none";

const emptyStateClass = "px-6 py-10";

const emptyTitleClass = "text-lg font-medium tracking-tight";

const emptyDescriptionClass = "mt-2 text-base leading-7 text-muted-foreground";

function replaceSubmission(items: Submission[], nextSubmission: Submission) {
  return items.map((item) =>
    item.id === nextSubmission.id ? nextSubmission : item,
  );
}

function AdminSubmissionsClient({
  initialSubmissions,
}: {
  initialSubmissions: Submission[];
}) {
  const [submissions, setSubmissions] = useState(initialSubmissions);
  const [selectedSubmission, setSelectedSubmission] =
    useState<Submission | null>(null);
  const [pendingSubmissionId, setPendingSubmissionId] = useState<string | null>(
    null,
  );

  async function updateStatus(
    currentSubmissionID: string,
    adminStatus: AdminStatus,
  ) {
    try {
      setPendingSubmissionId(currentSubmissionID);

      const response = await clientApiFetch<UpdateSubmissionStatusResponse>(
        `/admin/submissions/${currentSubmissionID}/status`,
        {
          method: "PATCH",
          body: JSON.stringify({ admin_status: adminStatus }),
        },
      );

      setSubmissions((current) =>
        replaceSubmission(current, response.submission),
      );

      setSelectedSubmission((current) =>
        current?.id === response.submission.id ? response.submission : current,
      );
    } catch (error) {
      console.log(error);
    } finally {
      setPendingSubmissionId(null);
    }
  }

  return (
    <main className={mainClass}>
      <section className={pageHeaderClass}>
        <h1 className={pageTitleClass}>Onboaring Review Queue</h1>
      </section>

      <Card className={tableCardClass}>
        <CardContent className={tableCardContentClass}>
          {submissions.length ? (
            <div className={tableScrollClass}>
              <table className={tableClass}>
                <thead className={tableHeadClass}>
                  <tr>
                    <th className={tableHeaderCellClass}>Client</th>
                    <th className={tableHeaderCellClass}>Service</th>
                    <th className={tableHeaderCellClass}>Readiness</th>
                    <th className={tableHeaderCellClass}>Review</th>
                    <th className={tableHeaderCellClass}>Created</th>
                    <th className={tableHeaderCellClass}>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {submissions.map((submission) => {
                    const isPending = pendingSubmissionId === submission.id;
                    const isApproved = submission.admin_status === "approved";
                    const isRejected = submission.admin_status === "rejected";

                    return (
                      <tr className={tableRowClass} key={submission.id}>
                        <td className={tableCellClass}>
                          <div className={clientInfoClass}>
                            <span className={clientNameClass}>
                              {submission.client_name}
                            </span>
                            <span className={clientEmailClass}>
                              {submission.client_email}
                            </span>
                          </div>
                        </td>
                        <td className={mutedTableCellClass}>
                          {submission.service_package}
                        </td>
                        <td className={tableCellClass}>
                          <ReadinessBadge
                            status={submission.readiness_status}
                          />
                        </td>
                        <td className={tableCellClass}>
                          <AdminStatusBadge status={submission.admin_status} />
                        </td>
                        <td className={mutedTableCellClass}>
                          {formatDate(submission.created_at)}
                        </td>
                        <td className={tableCellClass}>
                          <div className={tableActionsClass}>
                            <Button
                              className={tableButtonClass}
                              type="button"
                              onClick={() => setSelectedSubmission(submission)}
                            >
                              View
                            </Button>
                            <Button
                              disabled={isPending || isApproved || isRejected}
                              className={tableButtonClass}
                              type="button"
                              onClick={() =>
                                updateStatus(submission.id, "approved")
                              }
                            >
                              {isApproved ? "Approved" : "Approve"}
                            </Button>
                            <Button
                              disabled={isPending || isApproved || isRejected}
                              className={tableButtonClass}
                              onClick={() =>
                                updateStatus(submission.id, "rejected")
                              }
                              type="button"
                            >
                              {isRejected ? "Rejected" : "Reject"}
                            </Button>
                          </div>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          ) : (
            <div className={emptyStateClass}>
              <p className={emptyTitleClass}>No Submissions yet</p>
            </div>
          )}
        </CardContent>
      </Card>

      <SubmissionDetails
        submission={selectedSubmission}
        onClose={() => setSelectedSubmission(null)}
      />
    </main>
  );
}

export default AdminSubmissionsClient;
