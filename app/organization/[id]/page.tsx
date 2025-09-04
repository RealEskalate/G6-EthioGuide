import EditProcedurePage from "../(orgRelatedPages)/editProcedure/page";

// Next.js server component
export default async function EditProcedure({
  params,
}: {
  params: { id: string };
}) {
  const{ id } = await params
  // const { id } = params;

  // Fetch the procedure (dummy for now, replace with your API call)
  // const res = await fetch(`http://localhost:3000/api/procedures/${id}`);
  // const procedure = await res.json();
  const procedure = {
    id: "1",
    orgId: "org-001",
    title: "Driver's License Application",
    requirements: [
      { text: "Proof of identity" },
      { text: "Medical certificate" },
    ],
    steps: [
      { order: 1, text: "Fill out the application form" },
      { order: 2, text: "Submit required documents" },
      { order: 3, text: "Take vision and written tests" },
    ],
    fees: [{ label: "Application Fee", amount: 500, currency: "Birr" }],
    processingTime: { minDays: 5, maxDays: 10 },
    updatedAt: "2024-12-15T12:00:00Z",
    createdAt: "2024-10-01T09:00:00Z",
  };

  return <EditProcedurePage procedure={procedure} />;
}
