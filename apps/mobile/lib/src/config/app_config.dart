class AppConfig {
  const AppConfig({
    required this.apiBaseUrl,
    required this.appScheme,
  });

  factory AppConfig.fromEnvironment() {
    return const AppConfig(
      apiBaseUrl: String.fromEnvironment(
        'API_BASE_URL',
        defaultValue: 'http://localhost:8080',
      ),
      appScheme: String.fromEnvironment(
        'FLUTTER_APP_SCHEME',
        defaultValue: 'tipdrop',
      ),
    );
  }

  final String apiBaseUrl;
  final String appScheme;
}
