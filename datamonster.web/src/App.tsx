import { Routes, Route } from "react-router-dom";
import "./App.css";
import Login from "./routes/login";
import Settlements from "./routes/settlements/settlements";
import PublicLayout from "./routes/publiclayout";
import PrivateLayout from "./routes/privatelayout";

function App() {
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
