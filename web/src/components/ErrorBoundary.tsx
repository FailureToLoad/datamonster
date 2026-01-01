import { isRouteErrorResponse, useRouteError } from "react-router";
import styles from "./ErrorBoundary.module.css";

export function ErrorBoundary() {
  const error = useRouteError();

  let message = "Oops!";
  let details = "An unexpected error occurred.";

  if (isRouteErrorResponse(error)) {
    message = error.status === 404 ? "404" : `Error ${error.status}`;
    details =
      error.status === 404
        ? "The requested page could not be found."
        : error.statusText || details;
  } else if (error instanceof Error) {
    details = error.message;
  }

  return (
    <main className={styles.page}>
      <h1 className={styles.title}>{message}</h1>
      <p className={styles.details}>{details}</p>
      <a href="/" className={styles.homeLink}>Go home</a>
    </main>
  );
}
