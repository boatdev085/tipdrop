import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final destinations = <({String label, String path})>[
      (label: 'Public profile setup', path: '/profile/setup'),
      (label: 'Payment verification', path: '/payment-profile'),
      (label: 'Donor tip flow', path: '/tips/new'),
      (label: 'Worker dashboard', path: '/dashboard'),
      (label: 'Leaderboard', path: '/leaderboard'),
      (label: 'Discover', path: '/discover'),
      (label: 'QR scanner', path: '/scan'),
    ];

    return Scaffold(
      appBar: AppBar(title: const Text('TipDrop')),
      body: ListView.separated(
        padding: const EdgeInsets.all(20),
        itemCount: destinations.length,
        separatorBuilder: (_, __) => const SizedBox(height: 8),
        itemBuilder: (context, index) {
          final destination = destinations[index];
          return FilledButton.tonal(
            onPressed: () => context.go(destination.path),
            child: Text(destination.label),
          );
        },
      ),
    );
  }
}
