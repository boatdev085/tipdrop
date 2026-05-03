import 'package:flutter/material.dart';

class PaymentProfileScreen extends StatelessWidget {
  const PaymentProfileScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Payment verification')),
      body: ListView(
        padding: const EdgeInsets.all(20),
        children: [
          const TextField(decoration: InputDecoration(labelText: 'PromptPay ID')),
          const SizedBox(height: 12),
          const TextField(decoration: InputDecoration(labelText: 'Phone')),
          const SizedBox(height: 12),
          FilledButton(onPressed: () {}, child: const Text('Save payment profile')),
          const SizedBox(height: 12),
          OutlinedButton(onPressed: () {}, child: const Text('Send OTP')),
          const SizedBox(height: 12),
          const TextField(decoration: InputDecoration(labelText: 'OTP')),
          const SizedBox(height: 12),
          FilledButton(onPressed: () {}, child: const Text('Verify')),
        ],
      ),
    );
  }
}
