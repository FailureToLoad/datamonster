# datamonster

An app for managing the complexity of Kingdom Death: Monster campaigns.

## Set up

Everything is fairly manual while I work on getting the container structure straightened out.

The database was moved to its own private repo since I plan to load it with static data at some point.

Follow the [Firebase SDK setup instructions](https://firebase.google.com/docs/admin/setup) for go. This project relies on the `GOOGLE_APPLICATION_CREDENTIALS` environment variable

In `web` create an `env.local` file with the following values. They can be gotten from your firebase dash.

```
VITE_FIREBASE_API_KEY=
VITE_FIREBASE_AUTH_DOMAIN=
VITE_FIREBASE_DATABASE_URL=
VITE_FIREBASE_PROJECT_ID=
VITE_FIREBASE_STORAGE_BUCKET=
VITE_FIREBASE_MESSAGING_SENDER_ID=
VITE_FIREBASE_APP_ID=
VITE_FIREBASE_MEASUREMENT_ID=
```

The `VITE_` prefixing is required for vite to correctly process the files at startup.
