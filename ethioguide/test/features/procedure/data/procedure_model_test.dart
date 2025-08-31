import 'package:ethioguide/features/procedure/data/models/procedure_model.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  test('ProcedureModel serialization', () {
    final json = {
            'id': '1',
            'title': 'Passport Renewal',
            'category': 'Travel',
            'duration': '2-3 weeks',
            'cost': '1,200 ETB',
            'icon': 'passport',
            'isQuickAccess': true,
            'requiredDocuments': [],
            'steps': [],
            'resources': [],
            'feedback': []
          };

    final model = ProcedureModel.fromJson(json);
    expect(model.toJson(), json);
  });
}


