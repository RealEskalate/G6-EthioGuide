import 'package:equatable/equatable.dart';

/// Domain entity representing a government procedure shown in the UI.
class Procedure extends Equatable {

  final String id;
  final String title;
  final String category; // e.g., Transportation, Travel
  final String duration; // e.g., 2-3 weeks
  final String cost; // e.g., 1,200 ETB
  final String icon; // a semantic icon name or asset key
  final bool isQuickAccess; // whether to show in quick access grid
  final String organization; // a short organization
  final String description; // a short description
  final Status status; // a link to the official procedure page
  



  // Detail Page fields
  final List<String> requiredDocuments; // e.g., ["Passport Photo", "Birth Certificate"]
  final List<ProcedureStep> steps; // Step-by-step guide
  final List<Resource> resources; // forms, downloadable items
  final List<FeedbackItem> feedback; // user reviews

  const Procedure({
    required this.id,
    required this.title,
    required this.category,
    required this.duration,
    required this.cost,
    required this.icon,
    required this.isQuickAccess,
    this.requiredDocuments = const [],
    this.steps = const [],
    this.resources = const [],
    this.feedback = const [],
  });

  @override
  List<Object?> get props => [
        id,
        title,
        category,
        duration,
        cost,
        icon,
        isQuickAccess,
        requiredDocuments,
        steps,
        resources,
        feedback,
      ];
}


class ProcedureStep extends Equatable {
  final int number;
  final String title;
  final String description;

  const ProcedureStep({required this.number, required this.title, required this.description});

  @override
  List<Object?> get props => [number, title, description];
}

class Resource extends Equatable {
  final String name;
  final String url;

  const Resource({required this.name, required this.url});

  @override
  List<Object?> get props => [name, url];
}

class FeedbackItem extends Equatable {
  final String user;
  final String comment;
  final String date;
  final bool verified;

  const FeedbackItem({required this.user, required this.comment, required this.date, required this.verified});

  @override
  List<Object?> get props => [user, comment, date, verified];
}


