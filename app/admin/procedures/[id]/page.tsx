import EditProcedurePage from "../../(adminRelatedPages)/editProcedure/page";
import ProcedureProp from "@/types/procedure";

interface Props {
  params: { id: string };
}

export default async function EditProcedure({ params }: Props) {
  const { id } = await params;

  const res = await fetch(
    `https://ethio-guide-backend.onrender.com/api/v1/procedures/68bad97d299bfa90117809e`,
    { cache: "no-store" } // optional: ensures fresh data
  );

  if (!res.ok) {
    throw new Error("Failed to fetch procedure");
  }

  const procedure: ProcedureProp = await res.json();

  return <EditProcedurePage procedure={procedure} />;
}
