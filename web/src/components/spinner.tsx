function Spinner() {
  return (
    <div className="flex h-screen w-full flex-col items-center justify-center">
      <div
        className="inline-block h-6 w-6 animate-spin rounded-full border-[3px] border-current border-t-transparent text-gray-800 dark:text-white"
        role="status"
        aria-label="loading"
      >
        <span className="sr-only">Loading...</span>
      </div>
    </div>
  );
}

export default Spinner;
