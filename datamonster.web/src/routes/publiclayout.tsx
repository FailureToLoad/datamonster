import { AuthContext } from "@/auth/auth-context";
import { useContext } from "react";
import { Navigate, Outlet } from "react-router-dom";

function PublicLayout() {
  const { isLoggedIn } = useContext(AuthContext);
  if (isLoggedIn) return <Navigate to="/select" />;
  return <Outlet />;
}

export default PublicLayout;
