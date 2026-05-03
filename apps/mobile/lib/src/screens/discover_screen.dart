import 'package:flutter/material.dart';

class DiscoverScreen extends StatelessWidget {
  const DiscoverScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Discover')),
      body: ListView(
        padding: const EdgeInsets.all(20),
        children: const [
          TextField(decoration: InputDecoration(labelText: 'Search workers')),
          SizedBox(height: 16),
          ListTile(title: Text('Featured worker'), subtitle: Text('Public profile')),
        ],
      ),
    );
  }
}
