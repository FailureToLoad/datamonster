import { Routes, Route } from "react-router-dom";
import "./App.css";
import Login from "./routes/login";
import SettlementSelector from "./routes/settlementSelector/settlementSelector";
import PublicLayout from "./routes/publiclayout";
import PrivateLayout from "./routes/privatelayout";
import Settlement from "./routes/settlement";
import { SettlementProvider } from "./routes/settlement/context";

function App() {
  return (
    <Routes>
      <Route element={<PrivateLayout />}>
        <Route path="/" element={<SettlementSelector />} />
        <Route
          path="settlement/:settlementId/"
          element={
            <SettlementProvider>
              <Settlement />
            </SettlementProvider>
          }
        />
      </Route>
      <Route element={<PublicLayout />}>
        <Route path="login" element={<Login />} />
      </Route>
    </Routes>
  );
}

export default App;
