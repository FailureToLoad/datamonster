export default function Home() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-background">
      <h1 className="mb-4 text-5xl font-extrabold leading-none tracking-tight">
        Datamonster
      </h1>
      <a href="/auth/login">Sign In</a>
    </div>
  );
}
