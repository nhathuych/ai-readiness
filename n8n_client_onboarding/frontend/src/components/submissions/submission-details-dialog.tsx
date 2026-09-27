"use client";

import { Submission } from "@/types/submission";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card";
import { ReadinessBadge } from "../shared/readiness-badge";
import AdminStatusBadge from "../shared/admin-status-badge";
import { formatDate } from "@/lib/utils";
import { Button } from "../ui/button";

const detailRowClass =
  "grid gap-1 border-b border-border pb-3 last:border-b-0 last:pb-0";

const detailLabelClass =
  "text-xs font-medium uppercase tracking-[0.16em] text-muted-foreground";

const detailValueClass = "text-sm leading-6 text-foreground";

const textBlockClass = "grid gap-2";

const textBlockTitleClass = "text-sm font-medium text-foreground";

const textBlockValueClass =
  "min-h-20 whitespace-pre-wrap border border-border bg-muted/40 p-4 text-sm leading-6 text-muted-foreground";

const overlayClass =
  "fixed inset-0 z-50 overflow-y-auto bg-foreground/10 px-4 py-8 backdrop-blur-sm";

const dialogWrapperClass = "mx-auto w-full max-w-5xl";

const cardClass = "rounded-none border-border shadow-none";

const cardHeaderClass = "gap-5 border-b border-border pb-5";

const headerContentClass =
  "flex flex-col gap-4 md:flex-row md:items-start md:justify-between";

const clientInfoClass = "grid gap-2";

const clientNameClass = "text-2xl font-semibold tracking-tight";

const clientEmailClass = "text-sm text-muted-foreground";

const badgeWrapperClass = "flex flex-wrap gap-2 md:justify-end";

const cardContentClass = "grid gap-6 pt-6";

const detailsGridClass = "grid gap-4 md:grid-cols-2 xl:grid-cols-3";

const aiSectionClass = "grid gap-4 border-t border-border pt-6 md:grid-cols-2";

const missingItemsSectionClass = "grid gap-3";

const missingItemsTitleClass = "text-sm font-medium text-foreground";

const missingItemsWrapperClass = "flex flex-wrap gap-2";

const missingItemClass =
  "border border-amber-200 bg-amber-50 px-3 py-1.5 text-sm text-amber-800";

const emptyMissingItemsClass =
  "border border-border bg-muted/40 p-4 text-sm text-muted-foreground";

const footerClass = "flex justify-end border-t border-border pt-5";

const closeButtonClass =
  "h-11 rounded-none px-5 text-sm font-medium shadow-none";

function DetailRow({ value, label }: { label: string; value: string }) {
  return (
    <div className={detailRowClass}>
      <span className={detailLabelClass}>{label}</span>
      <span className={detailValueClass}>{value}</span>
    </div>
  );
}

function TextBlock({ title, value }: { title: string; value: string }) {
  return (
    <section className={textBlockClass}>
      <p className={textBlockTitleClass}>{title}</p>
      <div className={textBlockValueClass}>{value}</div>
    </section>
  );
}

function SubmissionDetails({
  submission,
  onClose,
}: {
  submission: Submission | null;
  onClose: () => void;
}) {
  if (!submission) return null;

  return (
    <div className={overlayClass}>
      <div className={dialogWrapperClass}>
        <Card className={cardClass}>
          <CardHeader className={cardHeaderClass}>
            <div className={headerContentClass}>
              <div className={clientInfoClass}>
                <CardTitle className={clientNameClass}>
                  {submission?.client_name}
                </CardTitle>
                <p className={clientEmailClass}>{submission.client_email}</p>
              </div>
              <div className={badgeWrapperClass}>
                <ReadinessBadge status={submission.readiness_status} />
                <AdminStatusBadge status={submission.admin_status} />
              </div>
            </div>
          </CardHeader>

          <CardContent className={cardContentClass}>
            <section className={detailsGridClass}>
              <DetailRow label="Client Email" value={submission.client_email} />
              <DetailRow label="Service" value={submission.service_package} />
              <DetailRow label="Timeline" value={submission.desired_timeline} />
              <DetailRow label="Assets" value={submission.assets_provided} />
              <DetailRow
                label="Created"
                value={formatDate(submission.created_at)}
              />
              <DetailRow
                label="Updated"
                value={formatDate(submission.updated_at)}
              />
            </section>
            <TextBlock title="Project Goal" value={submission.project_goal} />
            <section className={aiSectionClass}>
              <TextBlock title="AI Summary" value={submission.ai_summary} />
              <TextBlock
                title="Recommended Next Action"
                value={submission.recommended_next_action}
              />
            </section>
            <section className={missingItemsSectionClass}>
              <p className={missingItemsTitleClass}>Mission Items</p>
              {submission.missing_items.length ? (
                <div className={missingItemsWrapperClass}>
                  {submission.missing_items.map((missingItem) => (
                    <span className={missingItemClass} key={missingItem}>
                      {missingItem}
                    </span>
                  ))}
                </div>
              ) : (
                <p className={emptyMissingItemsClass}>No missing items found</p>
              )}
            </section>

            <section className={footerClass}>
              <Button
                className={closeButtonClass}
                onClick={onClose}
                type="button"
              >
                Close
              </Button>
            </section>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

export default SubmissionDetails;
