import 'package:flutter/material.dart';

class WorkerDashboardScreen extends StatelessWidget {
  const WorkerDashboardScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Worker dashboard')),
      body: ListView(
        padding: const EdgeInsets.all(20),
        children: [
          Card(
            child: ListTile(
              title: const Text('Pending slip'),
              subtitle: const Text('Ref TD-0001 • THB 100'),
              trailing: Wrap(
                spacing: 8,
                children: [
                  TextButton(onPressed: () {}, child: const Text('Dispute')),
                  FilledButton(onPressed: () {}, child: const Text('Confirm')),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
