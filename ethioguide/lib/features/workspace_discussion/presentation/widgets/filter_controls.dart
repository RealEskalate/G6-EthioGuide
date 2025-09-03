import 'package:flutter/material.dart';

/// Widget that displays filter controls for discussions
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
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 20),
          decoration: BoxDecoration(
            color: Colors.grey[50],
            borderRadius: BorderRadius.circular(8),
          ),
          child: Row(
            children: [
              // Categories dropdown
  /*             Expanded(
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
                    ...categories.map(
                      (category) => DropdownMenuItem<String>(
                        value: category,
                        child: Text(category),
                      ),
                    ),
                  ],
                  onChanged: onCategoryChanged,
                ),
              ), */

              Expanded(
  child: DropdownButtonFormField<String>(
    value: selectedCategory,
    decoration: const InputDecoration(
      border: InputBorder.none,
      contentPadding: EdgeInsets.zero,
    ),
    dropdownColor: Colors.white, // background of dropdown
    icon: const Icon(Icons.arrow_drop_down, color: Colors.black),
    style: const TextStyle(color: Colors.black, fontSize: 14),
    items: [
      DropdownMenuItem<String>(
        value: null,
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
          decoration: BoxDecoration(
            color: Colors.teal.shade100, // ✅ background like screenshot
            borderRadius: BorderRadius.circular(6),
          ),
          child: const Text(
            'All Categories',
            style: TextStyle(
              color: Colors.teal, // ✅ green text
              fontWeight: FontWeight.bold,
            ),
          ),
        ),
      ),
      ...categories.map(
        (category) => DropdownMenuItem<String>(
          value: category,
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
            child: Text(category),
          ),
        ),
      ),
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
              
              // Filter icon
              
            ],
          ),
        ),
        const SizedBox(height: 8),

        // Horizontal scroll indicator
      ],
    );
  }
}