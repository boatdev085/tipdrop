import 'package:flutter/material.dart';

class DonorTipFlowScreen extends StatelessWidget {
  const DonorTipFlowScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Send tip')),
      body: ListView(
        padding: const EdgeInsets.all(20),
        children: [
          const TextField(decoration: InputDecoration(labelText: 'Worker username')),
          const SizedBox(height: 12),
          const TextField(decoration: InputDecoration(labelText: 'Amount')),
          const SizedBox(height: 12),
          FilledButton(onPressed: () {}, child: const Text('Initiate tip')),
          const SizedBox(height: 12),
          OutlinedButton(onPressed: () {}, child: const Text('Upload slip')),
          const SizedBox(height: 12),
          const Text('Status: not started'),
        ],
      ),
    );
  }
}
