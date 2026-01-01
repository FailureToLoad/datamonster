import styles from "./Home.module.css";

export default function Home() {
  return (
    <div className={styles.page}>
      <h1 className={styles.title}>
        Datamonster
      </h1>
      <a href="/auth/login">Sign In</a>
    </div>
  );
}
