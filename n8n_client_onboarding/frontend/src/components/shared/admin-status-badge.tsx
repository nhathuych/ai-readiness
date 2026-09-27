import { AdminStatus } from "@/types/submission";
import { Badge } from "../ui/badge";

const labels: Record<AdminStatus, string> = {
  pending: "Pending review",
  approved: "Approved",
  rejected: "Rejected",
};

function AdminStatusBadge({ status }: { status: AdminStatus }) {
  if (status === "approved") {
    return <Badge variant="default">{labels[status]}</Badge>;
  }

  if (status === "rejected") {
    return <Badge variant="destructive">{labels[status]}</Badge>;
  }

  return <Badge variant="ghost">{labels[status]}</Badge>;
}

export default AdminStatusBadge;
