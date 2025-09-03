import 'package:ethioguide/features/authentication/domain/entities/tokens.dart';


class TokensModel extends Tokens {
  const TokensModel({
    required super.accessToken,
    required super.refreshToken,
  });

  factory TokensModel.fromJson(Map<String, dynamic> json) {
    return TokensModel(
      accessToken: json['token'],
      refreshToken: json['refreshToken'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'token': accessToken,
      'refreshToken': refreshToken,
    };
  }
}