import { Button } from "@/components/ui/button";
import { useAuth0 } from "@auth0/auth0-react";
import { Link, Outlet } from "react-router-dom";

export default function Layout() {
  const { logout, isAuthenticated } = useAuth0();
  return (
    <>
      <div id="header" className="sticky top-0 w-full flex-none">
        {isAuthenticated ? (
          <Button onClick={async () => await logout()}>Sign Out</Button>
        ) : (
          <Link to="/select">Sign in</Link>
        )}
      </div>
      <Outlet />
    </>
  );
}
