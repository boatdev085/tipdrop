import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'screens/discover_screen.dart';
import 'screens/donor_tip_flow_screen.dart';
import 'screens/home_screen.dart';
import 'screens/leaderboard_screen.dart';
import 'screens/login_screen.dart';
import 'screens/payment_profile_screen.dart';
import 'screens/profile_setup_screen.dart';
import 'screens/qr_scanner_screen.dart';
import 'screens/splash_screen.dart';
import 'screens/worker_dashboard_screen.dart';

final _router = GoRouter(
  initialLocation: '/splash',
  routes: [
    GoRoute(path: '/splash', builder: (context, state) => const SplashScreen()),
    GoRoute(path: '/login', builder: (context, state) => const LoginScreen()),
    GoRoute(path: '/', builder: (context, state) => const HomeScreen()),
    GoRoute(
      path: '/profile/setup',
      builder: (context, state) => const ProfileSetupScreen(),
    ),
    GoRoute(
      path: '/payment-profile',
      builder: (context, state) => const PaymentProfileScreen(),
    ),
    GoRoute(
      path: '/tips/new',
      builder: (context, state) => const DonorTipFlowScreen(),
    ),
    GoRoute(
      path: '/dashboard',
      builder: (context, state) => const WorkerDashboardScreen(),
    ),
    GoRoute(
      path: '/leaderboard',
      builder: (context, state) => const LeaderboardScreen(),
    ),
    GoRoute(path: '/discover', builder: (context, state) => const DiscoverScreen()),
    GoRoute(path: '/scan', builder: (context, state) => const QrScannerScreen()),
  ],
);

class TipDropApp extends StatelessWidget {
  const TipDropApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      title: 'TipDrop',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: const Color(0xFF0F766E)),
        useMaterial3: true,
      ),
      routerConfig: _router,
    );
  }
}
