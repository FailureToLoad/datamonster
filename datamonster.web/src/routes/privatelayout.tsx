import { AuthContext } from "@/auth/auth-context";
import { useContext } from "react";
import { Navigate, Outlet } from "react-router-dom";
import Spinner from "@/components/spinner";
function PrivateLayout() {
  const { isLoggedIn, isLoading } = useContext(AuthContext);
  if (isLoading) return <Spinner />;
  if (!isLoggedIn) return <Navigate to="/login" />;
  return <Outlet />;
}

export default PrivateLayout;
