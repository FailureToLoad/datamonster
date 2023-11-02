import { useContext, useEffect } from "react";
import { Routes, Route, useNavigate } from "react-router-dom";
import { AuthContext } from "./auth/auth-context";
import "./App.css";
import Login from "./routes/login";
import Settlements from "./routes/settlements";
import PublicLayout from "./routes/publiclayout";
import PrivateLayout from "./routes/privatelayout";

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
      <Route element={<PrivateLayout />}>
        <Route path="/" element={<Settlements />} />
      </Route>
      <Route element={<PublicLayout />}>
        <Route path="login" element={<Login />} />
      </Route>
    </Routes>
  );
}

export default App;
