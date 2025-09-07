import 'package:ethioguide/features/procedure/domain/entities/procedure_step.dart';
import 'package:flutter/material.dart';
import 'package:ethioguide/features/procedure/domain/entities/workspace_procedure.dart';




/// Widget that displays filters for procedures
class WorkspaceFilters extends StatelessWidget {
  final ProcedureStatus? selectedStatus;
  final String? selectedOrganization;
  final Function(ProcedureStatus?) onStatusChanged;
  final Function(String?) onOrganizationChanged;

  const WorkspaceFilters({
    super.key,
    this.selectedStatus,
    this.selectedOrganization,
    required this.onStatusChanged,
    required this.onOrganizationChanged,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Expanded(
          child: _FilterDropdown<ProcedureStatus>(
            label: 'Status:',
            value: selectedStatus,
            items: [
              const DropdownMenuItem(
                value: null,
                child: Text('All'),
              ),
              ...ProcedureStatus.values.map((status) => DropdownMenuItem(
                value: status,
                child: Text(status.displayName),
              )),
            ],
            onChanged: onStatusChanged,
          ),
        ),
        const SizedBox(width: 16),
        Expanded(
          child: _FilterDropdown<String>(
            label: 'Organization:',
            value: selectedOrganization,
            items: [
                          const DropdownMenuItem(
              value: null,
              child: Text('All Organizations'),
            ),
            ...WorkspaceSampleData.getSampleOrganizations().map((org) => DropdownMenuItem(
              value: org,
              child: Text(org),
            )),
            ],
            onChanged: onOrganizationChanged,
          ),
        ),
      ],
    );
  }
}

class _FilterDropdown<T> extends StatelessWidget {
  final String label;
  final T? value;
  final List<DropdownMenuItem<T>> items;
  final Function(T?) onChanged;

  const _FilterDropdown({
    required this.label,
    required this.value,
    required this.items,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
            fontWeight: FontWeight.w500,
          ),
        ),
        const SizedBox(height: 8),
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 12),
          decoration: BoxDecoration(
            border: Border.all(color: Colors.grey.shade300),
            borderRadius: BorderRadius.circular(8),
          ),
          child: DropdownButtonHideUnderline(
            child: DropdownButton<T>(
              value: value,
              items: items,
              onChanged: onChanged,
              isExpanded: true,
              hint: const Text('Select'),
              icon: const Icon(Icons.keyboard_arrow_down),
            ),
          ),
        ),
      ],
    );
  }
}


class WorkspaceSampleData {
  static List<String> getSampleOrganizations() {
    return [
      "Ministry of Education",
      "Ministry of Health",
      "Immigration Office",
      "Commercial Bank of Ethiopia",
      "Addis Ababa City Administration",
    ];
  }
}
