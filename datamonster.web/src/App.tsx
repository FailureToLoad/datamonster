import { useContext, useEffect } from "react";
import { Routes, Route, useNavigate } from "react-router-dom";
import { AuthContext } from "./auth/auth-context";
import "./App.css";
import Login from "./pages/login";
import Home from "./pages/home";
import RequireAuth from "./components/requireAuth";

function App() {
  const { currentUser } = useContext(AuthContext);
  const navigate = useNavigate();

  // NOTE: console log for testing purposes
  console.log("User:", !!currentUser);

  // Check if the current user exists on the initial render.
  useEffect(() => {
    if (currentUser) {
      navigate("/home");
    }
  }, [currentUser]);
  return (
    <Routes>
      <Route index element={<Login />} />
      <Route
        path="home"
        element={
          <RequireAuth>
            <Home />
          </RequireAuth>
        }
      />
    </Routes>
  );
}

export default App;
