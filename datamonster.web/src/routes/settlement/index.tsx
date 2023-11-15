import { useParams } from "react-router-dom";
import { SettlementProvider } from "./context";

function Settlement() {
  const { settlementId } = useParams();
  return (
    <SettlementProvider>
      <h1>Settlement Id: {settlementId}</h1>
    </SettlementProvider>
  );
}

export default Settlement;
