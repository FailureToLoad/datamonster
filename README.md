# datamonster

An app for managing the complexity of Kingdom Death: Monster campaigns.

## Set up
Everything is fairly manual while I work on getting the container structure straightened out.

Run the docker-compose in `datamonster.records` to stand up the postgres DB container. I know its not really needed but raw docker commands are a pain.
Afterwards, run the following command to get the container's IP address.
`docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' dm_db`

Set up an environment variable for the connection string in the form of 
EXPORT CONN_STRING="postgres://appuser:Password1@IPADDRESSHERE:5432/records"

Those creds aren't sticking around forever, its just convenient to have them here for the time being.

Follow the [Firebase SDK setup instructions](https://firebase.google.com/docs/admin/setup) for go. This project relies on the `GOOGLE_APPLICATION_CREDENTIALS` environment variable

In `datamonster.web` create an `env.local` file with the following values. They can be gotten from your firebase dash.
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
