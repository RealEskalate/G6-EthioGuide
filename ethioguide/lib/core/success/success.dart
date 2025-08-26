import 'package:equatable/equatable.dart';

class Success extends Equatable{
  final String message;

  const Success({this.message = 'success'}); // default success message
  
  @override
  List<Object?> get props => [message]; 
}
