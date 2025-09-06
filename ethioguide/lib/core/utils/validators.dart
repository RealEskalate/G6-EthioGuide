// This class will hold all our reusable validation logic.
class Validators {
  // Regular expression for a basic email format check.
  static final RegExp _emailRegExp = RegExp(
    r'^[a-zA-Z0-9.]+@[a-zA-Z0-9]+\.[a-zA-Z]+',
  );

  // Method to validate if a string is a valid email format.
  static String? validateEmail(String? value) {
    if (value == null || value.isEmpty) {
      return 'Email is required.';
    }
    if (!_emailRegExp.hasMatch(value)) {
      return 'Please enter a valid email address.';
    }
    return null; // Return null if the value is valid.
  }

  // Method to validate if a string is not empty.
  static String? validateRequired(String? value, String fieldName) {
    if (value == null || value.isEmpty) {
      return '$fieldName is required.';
    }
    return null;
  }
  
  // Method to validate password length.
  static String? validatePassword(String? value) {
    if (value == null || value.isEmpty) {
      return 'Password is required.';
    }
    if (value.length < 8) {
      return 'Password must be at least 8 characters long.';
    }
    // You can add more complex password rules here (uppercase, numbers, etc.)
    return null;
  }
}