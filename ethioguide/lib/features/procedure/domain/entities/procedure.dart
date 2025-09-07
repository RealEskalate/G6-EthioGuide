import 'package:equatable/equatable.dart';

/// Domain entity representing a government procedure shown in the UI.
class Procedure extends Equatable {

  final String id;
  final String title;
  final ProcessTime duration; // e.g., 2-3 weeks
  final String cost; // e.g., 1,200 ETB
  final List<String> requiredDocuments; // e.g., ["Passport Photo", "Birth Certificate"]
  final List<ProcedureStep> steps; // Step-by-step guide


   final List<Resource> resources; // forms, downloadable items
   final List<FeedbackItem> feedback; // user reviews
  // final String icon; // a semantic icon name or asset key
  // final bool isQuickAccess; // whether to show in quick access grid // a link to the official procedure /* page */

  const Procedure({
    required this.id,
    required this.title,
    required this.duration,
    required this.cost,
    this.requiredDocuments = const [],
    this.steps = const [],
    this.resources = const [],
     this.feedback = const [],
    // required this.icon,
    // required this.isQuickAccess,
  });

  @override
  List<Object?> get props => [
        id,
        title,
        duration,
        cost,
        requiredDocuments,
        steps,
        // resources,
        // feedback,
        //    icon,
        // isQuickAccess,
      ];
}


class ProcedureStep extends Equatable {
  final int number;
  final String title;


  const ProcedureStep({required this.number, required this.title});

  @override
  List<Object?> get props => [number, title];
}

class Resource extends Equatable {
  final String name;
  final String url;

  const Resource({required this.name, required this.url});

  @override
  List<Object?> get props => [name, url];
}

class ProcessTime extends Equatable {
  final int minday;
  final int maxday;

  const ProcessTime({required this.minday, required this.maxday});


  @override
  List<Object?> get props => [minday, maxday];
}

class Fee extends Equatable {
    final int amount;
  final String currency;

  const Fee({required this.amount, required this.currency});

  @override
  List<Object?> get props => [amount, currency];
}

class FeedbackItem extends Equatable {
  final String user;
  final String comment;
  final String date;


  const FeedbackItem({required this.user, required this.comment, required this.date});

  @override
  List<Object?> get props => [user, comment, date];
}


