"use client";

import { cn } from "@/lib/utils";
import {
  CreateSubmissionInput,
  CreateSubmissionResponse,
} from "@/types/submission";
import { Button, buttonVariants } from "../ui/button";
import { FormEvent, ReactNode, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card";
import { Label } from "../ui/label";
import { Input } from "../ui/input";
import { Textarea } from "../ui/textarea";
import { clientApiFetch } from "@/lib/api/client";

const initialForm: CreateSubmissionInput = {
  client_name: "",
  client_email: "",
  service_package: "Website / Landing Page",
  project_goal: "",
  desired_timeline: "",
  assets_provided: "Ready",
};

const serviceOptions = [
  "Website / Landing Page",
  "Automation / Operations",
  "Software / App Development",
  "Other",
];

const assetOptions = ["Ready", "Partially Ready", "Not Ready"];

const fieldClass = "grid gap-2";

const openButtonClass =
  "h-11 rounded-none px-5 text-sm font-medium shadow-none";

const overlayClass =
  "fixed inset-0 z-50 overflow-y-auto bg-foreground/10 px-4 py-8 backdrop-blur-sm";

const dialogCardClass =
  "mx-auto w-full max-w-4xl rounded-none border-border shadow-none";

const dialogHeaderClass = "border-b border-border pb-5";

const dialogTitleClass = "text-2xl font-semibold tracking-tight";

const dialogDescriptionClass = "text-sm leading-6 text-muted-foreground";

const dialogContentClass = "pt-6";

const formClass = "grid gap-5";

const formGridClass = "grid gap-5 md:grid-cols-2";

const selectClass =
  "h-11 rounded-none border border-input bg-background px-3 text-sm shadow-none focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring";

const errorMessageClass =
  "border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700";

const formActionsClass = "flex flex-wrap gap-3 border-t border-border pt-5";

const submitButtonClass =
  "h-11 rounded-none px-5 text-sm font-medium shadow-none";

const cancelButtonClass = cn(
  buttonVariants({ variant: "outline" }),
  "h-11 rounded-none border-border px-5 text-sm font-medium shadow-none",
);

function Field({ children, label }: { children: ReactNode; label: string }) {
  return (
    <section className={fieldClass}>
      <Label>{label}</Label>
      {children}
    </section>
  );
}

function CreateSubmissionClient({
  onCreated,
}: {
  onCreated: (payload: CreateSubmissionResponse) => void;
}) {
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState<CreateSubmissionInput>(initialForm);
  const [isSubmitting, setIsSubmitting] = useState(false);

  function closeDialog() {
    if (isSubmitting) {
      return;
    }

    setOpen(false);
    setForm(initialForm);
  }

  function setField(name: keyof CreateSubmissionInput, value: string) {
    setForm((current) => ({
      ...current,
      [name]: value,
    }));
  }

  async function handleFormSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    try {
      setIsSubmitting(true);

      const response = await clientApiFetch<CreateSubmissionResponse>(
        "/submissions/",
        {
          method: "POST",
          body: JSON.stringify(form),
        },
      );

      onCreated(response);
      closeDialog();
    } catch (err) {
      console.log(err);
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <>
      <Button className={openButtonClass} onClick={() => setOpen(true)}>
        New Onboarding Request
      </Button>
      {open ? (
        <div className={overlayClass}>
          <Card className={dialogCardClass}>
            <CardHeader className={dialogHeaderClass}>
              <CardTitle className={dialogTitleClass}>
                Submit Onboarding Request
              </CardTitle>
            </CardHeader>
            <CardContent className={dialogContentClass}>
              <form onSubmit={handleFormSubmit} className={formClass}>
                <section className={formGridClass}>
                  <Field label="Client Name">
                    <Input
                      value={form.client_name}
                      onChange={(event) =>
                        setField("client_name", event.target.value)
                      }
                      placeholder="Sangam Mukherjee"
                      required
                    />
                  </Field>
                  <Field label="Client Email">
                    <Input
                      value={form.client_email}
                      onChange={(event) =>
                        setField("client_email", event.target.value)
                      }
                      placeholder="sangam@gmail.com"
                      required
                    />
                  </Field>
                  <Field label="Service / Package">
                    <select
                      className={selectClass}
                      value={form.service_package}
                      onChange={(event) =>
                        setField("service_package", event.target.value)
                      }
                    >
                      {serviceOptions.map((option) => (
                        <option key={option} value={option}>
                          {option}
                        </option>
                      ))}
                    </select>
                  </Field>
                  <Field label="Desired Timeline">
                    <Input
                      value={form.desired_timeline}
                      onChange={(event) =>
                        setField("desired_timeline", event.target.value)
                      }
                      placeholder="4 weeks"
                      required
                    />
                  </Field>
                  <Field label="Assests Provided">
                    <select
                      className={selectClass}
                      value={form.assets_provided}
                      onChange={(event) =>
                        setField("assets_provided", event.target.value)
                      }
                    >
                      {assetOptions.map((option) => (
                        <option key={option} value={option}>
                          {option}
                        </option>
                      ))}
                    </select>
                  </Field>
                </section>
                <Field label="Project Goal">
                  <Textarea
                    rows={5}
                    value={form.project_goal}
                    onChange={(event) =>
                      setField("project_goal", event.target.value)
                    }
                    placeholder="We need a landing page redesign..."
                    required
                  />
                </Field>
                <section className={formActionsClass}>
                  <Button
                    className={submitButtonClass}
                    disabled={isSubmitting}
                    type="submit"
                  >
                    {isSubmitting ? "Checking with n8n..." : "Submit Request"}
                  </Button>
                  <Button
                    className={submitButtonClass}
                    type="button"
                    variant="destructive"
                    onClick={closeDialog}
                  >
                    Cancel
                  </Button>
                </section>
              </form>
            </CardContent>
          </Card>
        </div>
      ) : null}
    </>
  );
}

export default CreateSubmissionClient;
