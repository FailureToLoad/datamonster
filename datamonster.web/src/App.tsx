import { useContext, useEffect } from "react";
import { Routes, Route, useNavigate } from "react-router-dom";
import { AuthContext } from "./auth/auth-context";
import "./App.css";
import Login from "./pages/login";
import Settlements from "./pages/settlements";
import RequireAuth from "./components/requireAuth";

function App() {
  const { currentUser } = useContext(AuthContext);
  const navigate = useNavigate();

  // NOTE: console log for testing purposes
  console.log("User:", !!currentUser);

  // Check if the current user exists on the initial render.
  useEffect(() => {
    if (currentUser) {
      navigate("/");
    }
  }, [currentUser]);
  return (
    <Routes>
      <Route
        index
        element={
          <RequireAuth>
            <Settlements />
          </RequireAuth>
        }
      />
      <Route path="login" element={<Login />} />
    </Routes>
  );
}

export default App;
