import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class LoginScreen extends StatelessWidget {
  const LoginScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Sign in')),
      body: ListView(
        padding: const EdgeInsets.all(20),
        children: [
          FilledButton(
            onPressed: () => context.go('/'),
            child: const Text('Continue with Google'),
          ),
          const SizedBox(height: 12),
          OutlinedButton(
            onPressed: () => context.go('/'),
            child: const Text('Continue with Facebook'),
          ),
        ],
      ),
    );
  }
}
