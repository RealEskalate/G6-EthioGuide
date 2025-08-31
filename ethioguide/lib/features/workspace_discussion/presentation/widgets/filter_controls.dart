import 'package:flutter/material.dart';

/// Widget that displays filter controls for discussions
class FilterControls extends StatelessWidget {
  final String? selectedCategory;
  final String? selectedFilter;
  final List<String> categories;
  final Function(String?)? onCategoryChanged;
  final Function(String?)? onFilterChanged;

  const FilterControls({
    super.key,
    this.selectedCategory,
    this.selectedFilter,
    required this.categories,
    this.onCategoryChanged,
    this.onFilterChanged,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        // Filter bar
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
          decoration: BoxDecoration(
            color: Colors.grey[50],
            borderRadius: BorderRadius.circular(8),
          ),
          child: Row(
            children: [
              // Categories dropdown
              Expanded(
                child: DropdownButtonFormField<String>(
                  value: selectedCategory,
                  decoration: const InputDecoration(
                    labelText: 'Categories',
                    border: InputBorder.none,
                    contentPadding: EdgeInsets.zero,
                  ),
                  items: [
                    const DropdownMenuItem<String>(
                      value: null,
                      child: Text('All Categories'),
                    ),
                    ...categories.map((category) => DropdownMenuItem<String>(
                      value: category,
                      child: Text(category),
                    )),
                  ],
                  onChanged: onCategoryChanged,
                ),
              ),
              const SizedBox(width: 16),
              // Filter dropdown
              Expanded(
                child: DropdownButtonFormField<String>(
                  value: selectedFilter,
                  decoration: const InputDecoration(
                    labelText: 'Most Recent',
                    border: InputBorder.none,
                    contentPadding: EdgeInsets.zero,
                  ),
                  items: const [
                    DropdownMenuItem<String>(
                      value: 'recent',
                      child: Text('Most Recent'),
                    ),
                    DropdownMenuItem<String>(
                      value: 'trending',
                      child: Text('Trending'),
                    ),
                    DropdownMenuItem<String>(
                      value: 'popular',
                      child: Text('Most Popular'),
                    ),
                  ],
                  onChanged: onFilterChanged,
                ),
              ),
              const SizedBox(width: 16),
              // Filter icon
              Icon(
                Icons.filter_list,
                color: Colors.grey[600],
                size: 20,
              ),
            ],
          ),
        ),
        const SizedBox(height: 8),
        // Horizontal scroll indicator
        Container(
          height: 4,
          decoration: BoxDecoration(
            color: Colors.grey[300],
            borderRadius: BorderRadius.circular(2),
          ),
          child: Row(
            children: [
              Container(
                width: 100,
                height: 4,
                decoration: BoxDecoration(
                  color: Colors.grey[600],
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const Spacer(),
              IconButton(
                onPressed: () {},
                icon: const Icon(Icons.arrow_back_ios, size: 16),
                padding: EdgeInsets.zero,
                constraints: const BoxConstraints(),
              ),
              IconButton(
                onPressed: () {},
                icon: const Icon(Icons.arrow_forward_ios, size: 16),
                padding: EdgeInsets.zero,
                constraints: const BoxConstraints(),
              ),
            ],
          ),
        ),
      ],
    );
  }
}
