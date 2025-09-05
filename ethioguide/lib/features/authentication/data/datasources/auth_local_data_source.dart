import 'package:google_sign_in/google_sign_in.dart';
import 'package:ethioguide/core/error/exception.dart';
import 'package:flutter/foundation.dart';

abstract class AuthLocalDataSource {
  Future<String> getGoogleServerAuthCode();
}

class AuthLocalDataSourceImpl implements AuthLocalDataSource {
  AuthLocalDataSourceImpl();

  static const String _webClientId =
      '821529620719-ciu7h4g2k0vrjbc9gfs8p7r08e9mf4ht.apps.googleusercontent.com';

  @override
  Future<String> getGoogleServerAuthCode() async {
    try {
      final GoogleSignIn googleSignIn = GoogleSignIn(
        scopes: const ['email', 'profile'],
        clientId: kIsWeb ? _webClientId : null,
      );

      // Ensure a clean session
      await googleSignIn.disconnect();

      final GoogleSignInAccount? account = await googleSignIn.signIn();
      if (account == null) {
        throw CacheException(message: 'Sign-in cancelled by user.');
      }

      // Prefer the account.serverAuthCode when available
      String? authCode = account.serverAuthCode;

      // Fallback to authentication.serverAuthCode for older platforms
      if (authCode == null || authCode.isEmpty) {
        final GoogleSignInAuthentication authentication =
            await account.authentication;
        authCode = authentication.serverAuthCode;
      }

      if (authCode == null || authCode.isEmpty) {
        throw CacheException(
          message: 'Failed to retrieve Google server auth code.',
        );
      }

      return authCode;
    } catch (e) {
      if (e is CacheException) rethrow;
      throw CacheException(message: 'An error occurred during Google Sign-In.');
    }
  }
}