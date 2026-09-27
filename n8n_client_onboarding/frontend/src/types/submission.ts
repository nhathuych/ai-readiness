export type ReadinessStatus = "ready" | "missing_info";
export type AdminStatus = "pending" | "approved" | "rejected";

export type Submission = {
  id: string;
  user_id: string;
  client_name: string;
  client_email: string;
  service_package: string;
  project_goal: string;
  desired_timeline: string;
  assets_provided: string;
  readiness_status: ReadinessStatus;
  missing_items: string[];
  ai_summary: string;
  recommended_next_action: string;
  admin_status: AdminStatus;
  created_at: string;
  updated_at: string;
};

export type CreateSubmissionInput = {
  client_name: string;
  client_email: string;
  service_package: string;
  project_goal: string;
  desired_timeline: string;
  assets_provided: string;
};

export type CreateSubmissionResponse = {
  submission: Submission;
};

export type SubmissionsResponse = {
  submissions: Submission[];
};

export type UpdateSubmissionStatusResponse = {
  submission: Submission;
};
