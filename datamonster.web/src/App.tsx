import { Routes, Route } from "react-router-dom";
import "./App.css";
import Login from "./routes/login";
import SettlementSelector from "./routes/settlementSelector";
import PublicLayout from "./routes/publiclayout";
import PrivateLayout from "./routes/privatelayout";
import Settlement from "./routes/settlement";

function App() {
  return (
    <Routes>
      <Route element={<PrivateLayout />}>
        <Route path="/select" element={<SettlementSelector />} />
        <Route path="/" element={<Settlement />}>
          <Route path="timeline" element={<div>timeline</div>} />
          <Route path="population" element={<div>population</div>} />
          <Route path="storage" element={<div>storage</div>} />
        </Route>
      </Route>
      <Route element={<PublicLayout />}>
        <Route path="login" element={<Login />} />
      </Route>
    </Routes>
  );
}

export default App;
