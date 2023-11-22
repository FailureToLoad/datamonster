import { Routes, Route } from "react-router-dom";
import "./App.css";
import Login from "./routes/login";
import SettlementSelector from "./routes/settlementSelector";
import PublicLayout from "./routes/publiclayout";
import PrivateLayout from "./routes/privatelayout";
import Settlement from "./routes/settlement";
import Population from "./routes/settlement/population";
import Timeline from "./routes/settlement/timeline";
import SettlementStorage from "./routes/settlement/settlementStorage";

function App() {
  return (
    <Routes>
      <Route element={<PrivateLayout />}>
        <Route path="/select" element={<SettlementSelector />} />
        <Route path="/" element={<Settlement />}>
          <Route path="timeline" element={<Timeline />} />
          <Route path="population" element={<Population />} />
          <Route path="storage" element={<SettlementStorage />} />
        </Route>
      </Route>
      <Route element={<PublicLayout />}>
        <Route path="login" element={<Login />} />
      </Route>
    </Routes>
  );
}

export default App;
