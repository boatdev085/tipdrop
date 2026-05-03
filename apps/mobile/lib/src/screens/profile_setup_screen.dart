import 'package:flutter/material.dart';

class ProfileSetupScreen extends StatelessWidget {
  const ProfileSetupScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Public profile')),
      body: const _StubForm(
        fields: ['Username', 'Display name', 'Job title', 'Bio'],
        action: 'Save profile',
      ),
    );
  }
}

class _StubForm extends StatelessWidget {
  const _StubForm({required this.fields, required this.action});

  final List<String> fields;
  final String action;

  @override
  Widget build(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.all(20),
      children: [
        for (final field in fields) ...[
          TextField(decoration: InputDecoration(labelText: field)),
          const SizedBox(height: 12),
        ],
        FilledButton(onPressed: () {}, child: Text(action)),
      ],
    );
  }
}
