"use client";

import { Submission } from "@/types/submission";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card";
import { ReadinessBadge } from "../shared/readiness-badge";
import { formatDate } from "@/lib/utils";
import AdminStatusBadge from "../shared/admin-status-badge";
import { Button } from "../ui/button";

const cardClass = "rounded-none border-border bg-card shadow-none";

const cardHeaderClass = "gap-4 border-b border-border pb-5";

const headerContentClass = "flex items-start justify-between gap-4";

const clientInfoClass = "grid gap-1";

const clientNameClass = "text-xl font-semibold tracking-tight";

const clientEmailClass = "text-sm text-muted-foreground";

const cardContentClass = "grid gap-5 pt-5";

const projectGoalClass =
  "line-clamp-2 text-base leading-7 text-muted-foreground";

const metaSectionClass = "grid gap-3 border-t border-border pt-4 text-sm";

const metaRowClass = "flex items-center justify-between gap-4";

const metaLabelClass = "text-muted-foreground";

const metaValueClass = "text-right text-foreground";

const viewButtonClass =
  "h-11 rounded-none px-5 text-sm font-medium shadow-none";

function SubmissionCard({
  submission,
  onView,
}: {
  submission: Submission;
  onView: (submission: Submission) => void;
}) {
  return (
    <Card className={cardClass}>
      <CardHeader className={cardHeaderClass}>
        <section className={headerContentClass}>
          <div className={clientInfoClass}>
            <CardTitle className={clientNameClass}>
              {submission.client_name}
            </CardTitle>
            <p className={clientEmailClass}>{submission.client_email}</p>
          </div>

          <ReadinessBadge status={submission.readiness_status} />
        </section>
      </CardHeader>

      <CardContent className={cardContentClass}>
        <p className={projectGoalClass}>{submission.project_goal}</p>

        <section className={metaSectionClass}>
          <div className={metaRowClass}>
            <span className={metaLabelClass}>Service</span>
            <span className={metaValueClass}>{submission.service_package}</span>
          </div>
          <div className={metaRowClass}>
            <span className={metaLabelClass}>Created</span>
            <span className={metaValueClass}>
              {formatDate(submission.created_at)}
            </span>
          </div>
          <div className={metaRowClass}>
            <span className={metaLabelClass}>Review</span>
            <AdminStatusBadge status={submission.admin_status} />
          </div>
        </section>

        <Button
          onClick={() => onView(submission)}
          className={viewButtonClass}
          type="button"
          variant="default"
        >
          View Details
        </Button>
      </CardContent>
    </Card>
  );
}

export default SubmissionCard;
