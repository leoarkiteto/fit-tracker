// Extends app.json. The Google Sign-In config plugin is added only when the
// iOS client ID is available, because it needs the reversed client ID as URL scheme.
module.exports = ({ config }) => {
  const iosClientId = process.env.EXPO_PUBLIC_GOOGLE_IOS_CLIENT_ID;
  if (!iosClientId) {
    return config;
  }

  const suffix = ".apps.googleusercontent.com";
  const prefix = iosClientId.endsWith(suffix)
    ? iosClientId.slice(0, -suffix.length)
    : iosClientId;
  const iosUrlScheme = `com.googleusercontent.apps.${prefix}`;

  return {
    ...config,
    plugins: [
      ...(config.plugins ?? []),
      ["@react-native-google-signin/google-signin", { iosUrlScheme }],
    ],
  };
};
