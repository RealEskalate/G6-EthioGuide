"use client";

import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";

type Requirement = { text: string; optional?: boolean };
type Step = { order: number; text: string };
type Fee = { label: string; amount: number; currency: string };
type ProcessingTime = { minDays?: number; maxDays?: number };

export default function AddProcedurePage() {
  const [title, setTitle] = useState("");
  const [requirements, setRequirements] = useState<Requirement[]>([]);
  const [steps, setSteps] = useState<Step[]>([]);
  const [fees, setFees] = useState<Fee[]>([]);
  const [processingTime, setProcessingTime] = useState<ProcessingTime>({});

  // Handlers for arrays
  const addRequirement = () => {
    if (
      requirements.length === 0 ||
      requirements[requirements.length - 1].text.trim() !== ""
    ) {
      setRequirements((prev) => [...prev, { text: "", optional: false }]);
    }
  };

  const updateRequirement = <K extends keyof Requirement>(
    i: number,
    field: K,
    value: Requirement[K]
  ) => {
    const copy = [...requirements];
    copy[i][field] = value;
    setRequirements(copy);
  };

  const addStep = () => {
    if (steps.length === 0 || steps[steps.length - 1].text.trim() !== "") {
      setSteps((prev) => [...prev, { order: prev.length + 1, text: "" }]);
    }
  };

  const updateStep = (i: number, value: string) => {
    const copy = [...steps];
    copy[i].text = value;
    setSteps(copy);
  };

  const addFee = () => {
    const last = fees[fees.length - 1];
    if (fees.length === 0 || (last.label.trim() !== "" && last.amount > 0)) {
      setFees((prev) => [...prev, { label: "", amount: 0, currency: "Birr" }]);
    }
  };

  const updateFee = <K extends keyof Fee>(
    i: number,
    field: K,
    value: Fee[K]
  ) => {
    const copy = [...fees];
    copy[i][field] = value;
    setFees(copy);
  };

  const handleSubmit = () => {
    const payload = {
      title,
      requirements,
      steps,
      fees,
      processingTime,
      updatedAt: new Date().toISOString(),
      createdAt: new Date().toISOString(),
    };
    console.log("Submitting:", payload);
    // TODO: send to API
  };

  return (
    <div className="max-w-2xl mx-auto space-y-6 text-primary-dark">
      <h2 className="text-xl font-semibold">Add New Procedure</h2>

      {/* Title */}
      <div>
        <Label htmlFor="title">Title</Label>
        <Input
          id="title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />
      </div>

      {/* Requirements */}
      <div>
        <h3 className="font-medium">Requirements</h3>
        {requirements.map((req, i) => (
          <div key={i} className="flex gap-2 items-center mb-2">
            <Input
              placeholder="Requirement text"
              value={req.text}
              onChange={(e) => updateRequirement(i, "text", e.target.value)}
            />
            <label className="flex items-center gap-1 text-sm">
              <input
                type="checkbox"
                checked={req.optional}
                onChange={(e) =>
                  updateRequirement(i, "optional", e.target.checked)
                }
              />
              Optional
            </label>
          </div>
        ))}
        <Button variant="outline" onClick={addRequirement}>
          + Add Requirement
        </Button>
      </div>

      {/* Steps */}
      <div>
        <h3 className="font-medium">Steps</h3>
        {steps.map((s, i) => (
          <div key={i} className="mb-2">
            <Input
              placeholder={`Step ${i + 1}`}
              value={s.text}
              onChange={(e) => updateStep(i, e.target.value)}
            />
          </div>
        ))}
        <Button variant="outline" onClick={addStep}>
          + Add Step
        </Button>
      </div>

      {/* Fees */}
      <div>
        <h3 className="font-medium">Fees</h3>
        {fees.map((f, i) => (
          <div key={i} className="flex gap-2 mb-2">
            <Input
              placeholder="Label"
              value={f.label}
              onChange={(e) => updateFee(i, "label", e.target.value)}
            />
            <Input
              type="number"
              min={0}
              placeholder="Amount"
              value={f.amount}
              onChange={(e) => updateFee(i, "amount", Number(e.target.value))}
            />
            <Input
              placeholder="Currency"
              value={f.currency}
              onChange={(e) => updateFee(i, "currency", e.target.value)}
            />
          </div>
        ))}
        <Button variant="outline" onClick={addFee}>
          + Add Fee
        </Button>
      </div>

      {/* Processing Time */}
      <div>
        <h3 className="font-medium">Processing Time</h3>
        <div className="flex gap-2">
          <Input
            type="number"
            min={0}
            placeholder="Min Days"
            value={processingTime.minDays ?? ""}
            onChange={(e) =>
              setProcessingTime({
                ...processingTime,
                minDays: Number(e.target.value),
              })
            }
          />
          <Input
            type="number"
            min={0}
            placeholder="Max Days"
            value={processingTime.maxDays ?? ""}
            onChange={(e) =>
              setProcessingTime({
                ...processingTime,
                maxDays: Number(e.target.value),
              })
            }
          />
        </div>
      </div>

      {/* Documents */}
      {/* <div>
        <h3 className="font-medium">Documents Required</h3>
        {documents.map((d, i) => (
          <div key={i} className="flex gap-2 mb-2">
            <Input
              placeholder="Document Name"
              value={d.name}
              onChange={(e) => updateDocument(i, "name", e.target.value)}
            />
            <Input
              placeholder="Template URL"
              value={d.templateUrl ?? ""}
              onChange={(e) => updateDocument(i, "templateUrl", e.target.value)}
            />
          </div>
        ))}
        <Button variant="outline" onClick={addDocument}>+ Add Document</Button>
      </div> */}

      {/* Tags */}
      {/* <div>
        <h3 className="font-medium">Tags</h3>
        <Input
          placeholder="Comma separated tags"
          value={tags.join(",")}
          onChange={(e) => setTags(e.target.value.split(","))}
        />
      </div> */}

      <Button className="bg-primary text-white" onClick={handleSubmit}>
        Save Procedure
      </Button>
    </div>
  );
}
