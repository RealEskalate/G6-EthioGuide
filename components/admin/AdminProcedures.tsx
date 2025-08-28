import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import {Trash2 } from "lucide-react"
import { BiSolidEdit } from "react-icons/bi";
import { FaEye } from "react-icons/fa";
import {
  TableCell,
} from "@/components/ui/table";
export default function AdminProcedures() {

  const procedures = [
    {
      name: "Driver's License Application",
      description: "Complete procedure for new driver's license registration",
      updated: "Dec 15, 2024",
    },
    {
      name: "Vehicle Registration Renewal",
      description: "Annual vehicle registration renewal process",
      updated: "Dec 12, 2024",
    },
    {
      name: "Commercial License Application",
      description: "Process for obtaining commercial driving permits",
      updated: "Dec 10, 2024",
    },
    {
      name: "Traffic Violation Appeal",
      description: "Process for appealing traffic violation penalties",
      updated: "Nov 28, 2024",
    },
    {
      name: "International Driving Permit",
      description: "Application process for international driving permits",
      updated: "Dec 8, 2024",
    },
  ]

  return (
    <div className="p-6 space-y-6 w-full">

      {/* Search + Filter Row */}
      <div className="flex items-center justify-between gap-4">

      <h1 className="text-xl font-semibold text-primary-dark">Procedure Management</h1>
        <Input placeholder="Search procedures..." className="max-w-sm" />
      </div>

      {/* Table */}
      <div className="border rounded-md shadow-sm">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b">
            <tr>
              <th className="text-left py-3 px-4 font-medium text-gray-600">Procedure Name</th>
              <th className="text-left py-3 px-4 font-medium text-gray-600">Last Updated</th>
              <th className="text-left py-3 px-4 font-medium text-gray-600 ">Actions</th>
            </tr>
          </thead>
          <tbody>
            {procedures.map((p, idx) => (
              <tr key={idx} className="border-b hover:bg-gray-50">
                <td className="py-3 px-4">
                  <div className="font-medium text-gray-900">{p.name}</div>
                  <div className="text-gray-500">{p.description}</div>
                </td>
                <td className="py-3 px-4 text-gray-600">{p.updated}</td>
                <TableCell className="flex space-x-2 mt-3">
                      <FaEye className="w-4 h-4 text-primary cursor-pointer" />
                      <BiSolidEdit className="w-4 h-4 text-primary cursor-pointer" />
                      <Trash2 className="w-4 h-4 text-red-600 cursor-pointer" />
                    </TableCell>
              </tr>
            ))}
          </tbody>
        </table>

        {/* Pagination */}
        <div className="flex items-center justify-between p-4 text-sm text-gray-500">
          <span>Showing 1 to 5 of 23 procedures</span>
          <div className="flex gap-2">
            <Button variant="outline" size="sm">Previous</Button>
            <Button variant="default" className="text-white" size="sm">1</Button>
            <Button variant="outline" size="sm">Next</Button>
          </div>
        </div>
      </div>
    </div>
  )
}
