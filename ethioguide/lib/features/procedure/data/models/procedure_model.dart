import 'package:ethioguide/features/procedure/domain/entities/procedure.dart';

class ProcedureModel extends Procedure {
  const ProcedureModel({
    required super.id,
    required super.title,
    required super.category,
    required super.duration,
    required super.cost,
    required super.icon,
    required super.isQuickAccess,
    super.requiredDocuments = const [],
    super.steps = const [],
    super.resources = const [],
    super.feedback = const [],
  });

  factory ProcedureModel.fromJson(Map<String, dynamic> json) {
    return ProcedureModel(
      id: json['id'].toString(),
      title: json['title'] ?? '',
      category: json['category'] ?? '',
      duration: json['duration'] ?? '',
      cost: json['cost'] ?? '',
      icon: json['icon'] ?? 'badge',
      isQuickAccess: json['isQuickAccess'] == true,
      requiredDocuments: (json['requiredDocuments'] as List?)?.map((e) => e.toString()).toList() ?? const [],
      steps: ((json['steps'] as List?) ?? const [])
          .map((e) => ProcedureStep(
                number: (e['number'] as num?)?.toInt() ?? 0,
                title: e['title']?.toString() ?? '',
                description: e['description']?.toString() ?? '',
              ))
          .toList(),
      resources: ((json['resources'] as List?) ?? const [])
          .map((e) => Resource(
                name: e['name']?.toString() ?? '',
                url: e['url']?.toString() ?? '',
              ))
          .toList(),
      feedback: ((json['feedback'] as List?) ?? const [])
          .map((e) => FeedbackItem(
                user: e['user']?.toString() ?? '',
                comment: e['comment']?.toString() ?? '',
                date: e['date']?.toString() ?? '',
                verified: e['verified'] == true,
              ))
          .toList(),
    );
  }

  Map<String, dynamic> toJson() => {
        'id': id,
        'title': title,
        'category': category,
        'duration': duration,
        'cost': cost,
        'icon': icon,
        'isQuickAccess': isQuickAccess,
        'requiredDocuments': requiredDocuments,
        'steps': steps
            .map((s) => {
                  'number': s.number,
                  'title': s.title,
                  'description': s.description,
                })
            .toList(),
        'resources': resources.map((r) => {'name': r.name, 'url': r.url}).toList(),
        'feedback': feedback
            .map((f) => {
                  'user': f.user,
                  'comment': f.comment,
                  'date': f.date,
                  'verified': f.verified,
                })
            .toList(),
      };
}


