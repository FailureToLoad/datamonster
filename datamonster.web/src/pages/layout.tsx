import "./App.css";
import Header from "@/components/header";

function Content() {
  return (
    <div
      id="content-container"
      className="flex flex-1 items-center justify-center"
    >
      Content
    </div>
  );
}

function Layout() {
  return (
    <>
      <div className="flex h-screen w-full flex-col">
        <Header />
        <Content />
      </div>
    </>
  );
}

export default Layout;
