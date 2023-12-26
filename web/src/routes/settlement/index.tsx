import { Settlement } from "@/api/settlement";
import { Navigate } from "react-router-dom";
export { default as Settlement } from "./settlement";

export function SettlementLoader() {
  const settlementJson = localStorage.getItem("settlement");
  const settlement = settlementJson
    ? (JSON.parse(settlementJson) as Settlement)
    : null;
  if (!settlement) {
    return <Navigate to="/select" />;
  }
  return settlement;
}
