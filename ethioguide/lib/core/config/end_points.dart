class EndPoints {
  static final baseUrl = 'https://ethio-guide-backend.onrender.com/api/v1/';
  static final refreshTokenEndPoint = '/auth/refresh';
  static final sendQueryEndPoint = '/ai/guide';
  static final getHistoryEndPoint = '/ai/history';
  static final translateContentEndPoint = '/translate';
  static final createDiscussionEndPoint = 'discussions';
  static final getDiscussionsEndPoint = 'discussions';


  static final registerEndPoint = '/auth/register';
  static final loginEndPoint = '/auth/login';
  static const String verifyAccount = '/auth/verify';
  static const String forgotPassword = '/auth/forgot';
  static const String resetPassword = '/auth/reset';
  static const String refreshToken = '/auth/refresh';
  static const String socialLogin = '/auth/social';
  static const String getProfile = '/auth/me';
  static const String updateProfile = '/auth/me';
  static const String updatePassword = '/auth/me/password';

  static final translateContentEndPoint = 'ai/translate';

}
